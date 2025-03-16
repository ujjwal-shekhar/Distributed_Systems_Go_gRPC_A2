run-client:
	@go run services/client/*.go -R=$(R) -T=$(T)

run-server:
	@go run services/server/*.go -TYPE=$(TYPE) -PORT=$(PORT) -NUM_REDUCERS=$(NUM_REDUCERS) -TASK=$(TASK)

gen:
	@protoc \
		--proto_path=protobuf "protobuf/service.proto" \
		--go_out=services/common/genproto/comms --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/comms --go-grpc_opt=paths=source_relative

clean:
	@echo "Cleaning up mapResults directory..."
	@rm -rf ./mapResults/*
	@rm -rf ./reducerResults/*
	@echo "Killing processes holding ports 5000-5100 and 6000-6100..."
	@for port in $$(seq 5000 5100) $$(seq 6000 6100); do \
		PID=$$(lsof -t -i :$$port); \
		if [ -n "$$PID" ]; then \
			echo "Killing process $$PID holding port $$port"; \
			kill -9 $$PID; \
		fi; \
	done
	@echo "Cleanup complete."