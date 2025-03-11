run-client:
	@go run services/client/*.go -M=$(M) -R=$(R) -T=$(T)

run-server:
	@go run services/server/*.go -TYPE=$(TYPE) -PORT=$(PORT)

gen:
	@protoc \
		--proto_path=protobuf "protobuf/service.proto" \
		--go_out=services/common/genproto/comms --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/comms --go-grpc_opt=paths=source_relative