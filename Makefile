# general variables
OUT_DIR := ./out
COVER_FILE := $(OUT_DIR)/coverage.out

# testing parameters
VERBOSE ?= 0
COVERAGE ?= 0

# recipes list
.PHONY: proto test lint fmt clean
.PHONY: build push run 

# source files for tracking changes
SRC := $(shell find . -type f -name '*.go')

# recipe for building containers for every microservice
build:
	@docker build . -t vault-base -f ./cmd/Dockerfile
	@docker build . -t vault-generator -f ./cmd/generator/Dockerfile
	@docker build . -t vault-storage -f ./cmd/storage/Dockerfile
	@docker build . -t vault-encoder -f ./cmd/encoder/Dockerfile
	@docker build . -t vault-decoder -f ./cmd/decoder/Dockerfile

# recipe for tagging and pushing all containers to the repo
push: build
	@docker tag vault-generator:latest atennop/secure-vault:generator
	@docker push atennop/secure-vault:generator
	@docker tag vault-storage:latest atennop/secure-vault:storage
	@docker push atennop/secure-vault:storage
	@docker tag vault-encoder:latest atennop/secure-vault:encoder
	@docker push atennop/secure-vault:encoder
	@docker tag vault-decoder:latest atennop/secure-vault:decoder
	@docker push atennop/secure-vault:decoder

# recipe for running app in the k8s cluster
run: build
	minikube image load vault-generator:latest
	minikube image load vault-storage:latest
	minikube image load vault-encoder:latest
	minikube image load vault-decoder:latest

	@export $$(cat config/ports.env | xargs) && \
	for f in k8s/*.yaml; do \
		envsubst < $$f | kubectl apply -f -; \
	done

# recipe for generating Go code based on .proto files
proto:
	@protoc --go_out=. --go-grpc_out=. proto/generator.proto
	@protoc --go_out=. --go-grpc_out=. proto/storage.proto
	@echo ">> Go code generated."

# public recipe for testing (with or without coverage)
test: $(OUT_DIR)/test.cache.$(COVERAGE)

# public recipe for formatting
fmt: $(OUT_DIR)/fmt.cache

# public recipe for linting
lint: $(OUT_DIR)/lint.cache

# recipe for cleaning up garbage
clean:
	@echo ">> Cleaning up..."
	@rm -rf $(BIN_DIR) $(OUT_DIR)
	@kubectl delete -f k8s/
	@echo ">> Cleaned."

# creating output directory if it's missing
$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

# formatting source code
$(OUT_DIR)/fmt.cache: $(SRC) | $(OUT_DIR)
	@echo ">> Formatting..."
	@gofmt -s -w .
	@echo ">> Formatted."
	@touch $@

# linting source code
$(OUT_DIR)/lint.cache: $(SRC) | $(OUT_DIR)
	@echo ">> Linting..."
	@golangci-lint run ./... --timeout=5m > $(OUT_DIR)/lint.log || { \
		cat $(OUT_DIR)/lint.log; \
		exit 1; \
	};
	@echo ">> Linted."
	@touch $@

# testing source code (with or without coverage)
$(OUT_DIR)/test.cache.$(COVERAGE): $(SRC) | $(OUT_DIR)
	@set -e; \
	if [ "$(COVERAGE)" = "1" ]; then \
		echo ">> Testing with coverage..."; \
		go test ./... -coverprofile=$(COVER_FILE) -covermode=atomic > $(OUT_DIR)/test.log || { \
			cat $(OUT_DIR)/test.log; \
			exit 1; \
		}; \
		COVERAGE_OUTPUT=$$(go tool cover -func=$(COVER_FILE)); \
		if [ "$(VERBOSE)" = "1" ]; then \
			echo "$$COVERAGE_OUTPUT" | grep -v "total:"; \
		fi; \
		echo "$$COVERAGE_OUTPUT" | grep "total:" | awk '{print ">> Test coverage:", $$3}'; \
	else \
		echo ">> Testing..."; \
		go test ./... > $(OUT_DIR)/test.log || { \
			cat $(OUT_DIR)/test.log; \
			exit 1; \
		}; \
		if [ "$(VERBOSE)" = "1" ]; then \
			cat $(OUT_DIR)/test.log; \
		fi; \
		echo ">> All tests passed."; \
	fi;
	@touch $@
