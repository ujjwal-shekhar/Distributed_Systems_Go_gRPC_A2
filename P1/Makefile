run-client:
	@go run services/client/*.go $(ARGS)

run-server:
	@go run services/server/*.go

run-lb:
	@go run services/lb_server/*.go

gen:
	@protoc \
		--proto_path=protobuf "protobuf/service.proto" \
		--go_out=services/common/genproto/comms --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/comms --go-grpc_opt=paths=source_relative