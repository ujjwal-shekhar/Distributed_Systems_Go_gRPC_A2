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
	log.Printf("RBACUnaryInterceptor: method descriptor found : %s", info.FullMethod)

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
	if policy == nil {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: no policy defined")
	}
	log.Printf("RBACUnaryInterceptor: policy found=%v", policy)

	// For the method, if unauthenticated users are allowed, proceed with the RPC
	if policy.AllowUnauthenticated {
		return handler(ctx, req)
	}

	// Extract user role from context (assumed to be stored in metadata)
	claims, err := extractUserRole(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("RBACUnaryInterceptor: claims role=%s", claims.Role)

	// Check if role is allowed
	if !policy.AllowUnauthenticated {
		allowed := false
		for _, allowedRole := range policy.AllowedRoles {
			r, ok := pb.Role_value[claims.Role]
			if !ok {
				continue
			}
			log.Printf("RBACUnaryInterceptor: allowedRole=%s clientRole=%s", pb.Role_name[int32(allowedRole)], claims.Role)

			if r == int32(allowedRole) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied for role: %s", claims.Role)
		}
	}
	log.Printf("RBACUnaryInterceptor: permission granted for role: %s", claims.Role)

	// If the method is from the stripe service
	// Overwrite the userdetails with the ones obtained from the JWT token
	// This is done to ensure that the user cannot spoof their role
	// by sending a different role in the request
	if strings.Contains(info.FullMethod, "StripeService") {
		msg, ok := req.(proto.Message); if !ok {
			return nil, errors.New("could not convert request to proto.Message")
		}

		reflection := msg.ProtoReflect()
		descriptor := reflection.Descriptor()
		
		// Set "Username" field
		usernameField := descriptor.Fields().ByName("username")
		if usernameField != nil {
			reflection.Set(usernameField, protoreflect.ValueOfString(claims.Username))
		}

		// Set "Bankname" field
		banknameField := descriptor.Fields().ByName("bankname")
		if banknameField != nil {
			reflection.Set(banknameField, protoreflect.ValueOfString(claims.Bankname))
		}
	}
	log.Printf("RBACUnaryInterceptor: user details overwritten with JWT claims")

	// Proceed with the RPC
	return handler(ctx, req)
}

func extractUserRole(ctx context.Context) (*Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata not found in context")
	}
	log.Printf("extractUserRole: metadata found")

	token, ok := md["authorization"]
	if !ok {
		return nil, errors.New("role not found in metadata")
	}
	log.Printf("extractUserRole: token found")

	// We extract the bearer token from the metadata
	// and then we extract the role from the token
	// the token is parsed by jwt then
	tokenTrimmed := strings.TrimPrefix(token[0], "Bearer ")
	claims, err := ParseJWT(tokenTrimmed)
	if err != nil {
		return nil, err
	}
	log.Printf("extractUserRole: claims found")

	return claims, nil
}
