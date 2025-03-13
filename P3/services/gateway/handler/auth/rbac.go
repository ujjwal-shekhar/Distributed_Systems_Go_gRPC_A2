package auth

import (
	"context"
	"errors"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	pb "github.com/ujjwal-shekhar/stripe-clone/services/common/genproto/comms"
)

func RBACUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Convert method name from "/package.service/method" to "package.service.method"
	// and then I find the method descriptor to extract the policy later on
	methodName := strings.Replace(strings.TrimPrefix(info.FullMethod, "/"), "/", ".", -1)
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(methodName))
	if err != nil {
		return nil, err
	}

	method, ok := desc.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, status.Errorf(codes.Internal, "invalid method descriptor: %v", desc)
	}

	log.Printf("RBACUnaryInterceptor: %s", info.FullMethod)
	// Now, I will extract the RBAC policy which is stored in the method options
	// It contains the allowed roles and whether unauthenticated users are allowed
	// I just I iterate over the method options to find the policy
	// I then unmarshal the policy into a proto message
	var policy *pb.RBAC
	method.Options().ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.FullName() != pb.E_AccessControl.TypeDescriptor().FullName() {
			return true // Continue iterating
		}

		b, err := proto.Marshal(v.Message().Interface())
		if err != nil {
			return false // Stop iteration, handle error later
		}

		policy = &pb.RBAC{}
		if err := proto.Unmarshal(b, policy); err != nil {
			return false // Stop iteration, handle error later
		}

		return false // Stop iteration once found
	})

	log.Printf("RBACUnaryInterceptor: policy=%v", policy)

	if policy == nil {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: no policy defined")
	}

	if policy.AllowUnauthenticated {
		return handler(ctx, req)
	}

	// Extract user role from context (assumed to be stored in metadata)
	role, err := extractUserRole(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required")
	}

	log.Printf("RBACUnaryInterceptor: role=%s", role)

	// Check if role is allowed
	if !policy.AllowUnauthenticated {
		allowed := false
		for _, allowedRole := range policy.AllowedRoles {
			r, ok := pb.Role_value[role]
			if !ok {
				continue
			}

			if r == int32(allowedRole) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied for role: %s", role)
		}
	}

	// Proceed with the RPC
	return handler(ctx, req)
}

func extractUserRole(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found in context")
	}

	token, ok := md["authorization"]
	if !ok {
		return "", errors.New("role not found in metadata")
	}

	// We extract the bearer token from the metadata
	// and then we extract the role from the token
	// the token is parsed by jwt then
	tokenTrimmed := strings.TrimPrefix(token[0], "Bearer ")
	claims, err := ParseJWT(tokenTrimmed)
	if err != nil {
		return "", err
	}

	return claims.Role, nil
}
