run-client:
	@go run services/client/*.go -UNAME=$(UNAME) -BNAME=$(BNAME) -PASS=$(PASS)

run-server:
	@go run services/server/*.go -NAME=$(NAME)

run-gateway:
	@go run services/gateway/*.go

gen:
	@protoc \
		--proto_path=protobuf "protobuf/service.proto" \
		--go_out=services/common/genproto/comms --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/comms --go-grpc_opt=paths=source_relative

cert:
	cd services/client/cert/ && ./gen.sh && \
	cd ../../server/cert/ && ./gen.sh && \
	cd ../../gateway/cert/ && ./gen.sh