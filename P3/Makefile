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

clean:
	@echo "Killing processe holding port 8080..."
	@for port in $$(seq 8080 8081); do \
		PID=$$(lsof -t -i :$$port); \
		if [ -n "$$PID" ]; then \
			echo "Killing process $$PID holding port $$port"; \
			kill -9 $$PID; \
		fi; \
	done
	@echo "Cleanup complete."