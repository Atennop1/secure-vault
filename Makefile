# general variables
OUT_DIR := ./out
COVER_FILE := $(OUT_DIR)/coverage.out

# testing parameters
VERBOSE ?= 0
COVERAGE ?= 0

# recipes list
.PHONY: proto fmt lint test 
.PHONY: build push run clean

# source files for tracking changes
SRC := $(shell find . -type f -name '*.go')

# recipe for generating Go code based on .proto files
proto:
	@protoc --go_out=. --go-grpc_out=. proto/generator.proto
	@protoc --go_out=. --go-grpc_out=. proto/storage.proto
	@echo ">> Go code generated."

# public recipe for formatting
fmt: $(OUT_DIR)/fmt.cache

# public recipe for linting
lint: $(OUT_DIR)/lint.cache

# public recipe for building
build: $(OUT_DIR)/build.cache

# public recipe for testing (with or without coverage)
test: $(OUT_DIR)/test.cache.$(COVERAGE)

# recipe for running app in the k8s cluster 
run: build
	@export $$(cat config/ports.env | xargs) && \
	for f in k8s/*.yaml; do \
		envsubst < $$f | kubectl apply -f -; \
	done

# recipe for tagging and pushing all containers to the repo
push: build
	@docker push atennop/secure-vault:generator
	@docker push atennop/secure-vault:storage
	@docker push atennop/secure-vault:encoder
	@docker push atennop/secure-vault:decoder

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

# recipe for building containers for every microservice
$(OUT_DIR)/build.cache: $(SRC) | $(OUT_DIR)
	@eval $(minikube docker-env) && \
	docker build . -t atennop/secure-vault:generator --target generator -f ./cmd/Dockerfile && \
	docker build . -t atennop/secure-vault:storage --target storage -f ./cmd/Dockerfile && \
	docker build . -t atennop/secure-vault:encoder --target encoder -f ./cmd/Dockerfile && \
	docker build . -t atennop/secure-vault:decoder --target decoder -f ./cmd/Dockerfile 
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
