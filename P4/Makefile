run-client:
	@go run services/client/*.go -N=$(N) -T=$(T)

run-server:
	@go run services/server/*.go -N=$(N) -T=$(T) -ID=$(ID) -PORT=$(PORT) -TYPE=$(TYPE)

gen:
	@protoc \
		--proto_path=protobuf "protobuf/service.proto" \
		--go_out=services/common/genproto/comms --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/comms --go-grpc_opt=paths=source_relative

clean:
	@echo "Killing processes holding ports 5000-5100..."
	@for port in $$(seq 5000 5100); do \
		PID=$$(lsof -t -i :$$port); \
		if [ -n "$$PID" ]; then \
			echo "Killing process $$PID holding port $$port"; \
			kill -9 $$PID; \
		fi; \
	done
	@rm *.out
	@echo "Cleanup complete."