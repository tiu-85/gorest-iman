GOBIN=$(GOPATH)/bin
GOBUILD=go build
GOTEST=go test
GOGET=go get
GOTOOL=go tool
GOLIST=go list
XML=xml
MIN_COVERAGE=80.0

-include .env
-include .env.version
-include .env.custom
export

DOCKER_PROJECT=gorest-iman
DOCKER_DIR=$(PROJECT_DIR)/deployments
DOCKER_COMPOSE=$(DOCKER_DIR)/docker-compose.yml
DOCKER_COMPOSE_INFRA=$(DOCKER_DIR)/docker-compose.infra.yml

PROJECT_DIR=$(shell pwd)
PROTO_GEN_DIR=$(PROJECT_DIR)/pkg/common/gen/pbv1
PROTO_DIR=$(PROJECT_DIR)/api/proto
PROTO_FILES=$(shell ls $(PROTO_DIR))

# protoc plugin versions
VERSION_GOGO_PROTOBUF=v1.3.1
VERSION_GOOGLE_PROTOBUF=v1.4.3
VERSION_GRPC_GATEWAY=v1.15.2
VERSION_SETTER=latest

define build_proto_file
    cat $(PROTO_DIR)/$(1) | \
	sed 's/\/\/\-\-404_RESPONSE\-\-\/\//$(PROTO_404_RESPONSE)/g' | \
	sed 's/\/\/\-\-409_RESPONSE\-\-\/\//$(PROTO_409_RESPONSE)/g' | \
	sed 's/\/\/\-\-PARTNERS_API\-\-\/\//$(PROTO_PARTNERS_API)/g' | \
	sed 's/\/\/\-\-SUBSCRIPTIONS_API\-\-\/\//$(PROTO_SUBSCRIPTIONS_API)/g' | \
	sed 's/\/\/\-\-VIRTUAL_CARDS_API\-\-\/\//$(PROTO_VIRTUAL_CARDS_API)/g' | \
	sed 's/\/\*\-\-PROP_UINT32\-\-\*\//$(PROTO_PROP_UINT32)/g' | \
	sed 's/\/\*\-\-PROP_TIME\-\-\*\//$(PROTO_PROP_TIME)/g' | \
	sed 's/\/\*\-\-PROP_DURATION\-\-\*\//$(PROTO_PROP_DURATION)/g' | \
	sed 's/\/\*\-\-PROP_NOT_NULL\-\-\*\//$(PROTO_PROP_NOT_NULL)/g' | \
	sed 's/\/\*\-\-PROP_SETTER_INCLUDE\-\-\*\//$(PROTO_PROP_SETTER_INCLUDE)/g' | \
	sed 's/\/\*\-\-PROP_PTR\-\-\*\//$(PROTO_PROP_PTR)/g' | \
	sed 's/\/\*\-\-PROTO_FILE_OPTIONS\-\-\*\//$(PROTO_FILE_OPTIONS)/g' | \
	sed 's/\/\/\-\-X_EXAMPLES\[\([a-z\/\.\_]*\)\]\-\-\/\//$(PROTO_X_EXAMPLE_START)"\1"$(PROTO_X_EXAMPLE_END)/g' | \
	sed 's/\/\/\-\-METHOD_SECURITY\-\-\/\//$(PROTO_METHOD_SECURITY)/g' > $(PROJECT_DIR)/build/api/$(1)
endef

define build_dir
	$(foreach BUILD_APP, $(shell ls $(PROJECT_DIR)/cmd/$(1)), $(call build_app,$(1),$(BUILD_APP));)
endef

define build_app
	if [ -f "$(PROJECT_DIR)/cmd/$(1)/$(2)/main.go" ]; then \
		echo build: $(1)/$(2); \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PROJECT_DIR)/build/$(2) $(PROJECT_DIR)/cmd/$(1)/$(2)/main.go; \
    fi
endef

buf_process:
	rm -rf $(PROJECT_DIR)/build/api
	mkdir -p $(PROJECT_DIR)/build/api
	$(foreach PROTO_FILE, $(PROTO_FILES), $(call build_proto_file,$(PROTO_FILE)))
	buf generate --path $(PROJECT_DIR)/build/api

buf_version:
	buf --version

buf_post_process:
	mv $(PROJECT_DIR)/pkg/common/gen/build/api/* $(PROTO_GEN_DIR)
	rm -rf $(PROJECT_DIR)/pkg/common/gen/build
	$(foreach PROTO_FILE, $(shell ls $(PROTO_GEN_DIR)), $(call post_process_proto_file,$(PROTO_FILE)))
	find $(PROTO_GEN_DIR) -empty -type f -delete

buf_generate: buf_process buf_post_process

install_protoc_plugins:
	go install github.com/gogo/protobuf/protoc-gen-gogoslick@$(VERSION_GOGO_PROTOBUF)
	go install github.com/golang/protobuf/protoc-gen-go@$(VERSION_GOOGLE_PROTOBUF)
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@$(VERSION_GRPC_GATEWAY)
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@$(VERSION_GRPC_GATEWAY)
	go install github.com/mikekonan/protoc-gen-setter@$(VERSION_SETTER)

pull_images:
	docker pull postgres:latest

up-infra: pull_images
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) up -d --force-recreate

build_apps:
	@echo "Wait please! The build is running... "
	@$(foreach BUILD_APP, $(shell ls $(PROJECT_DIR)/cmd), $(call build_dir,$(BUILD_APP)))

up-services: build_apps
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE) up -d --build --force-recreate

down-infra:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) down --remove-orphans

down-services:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE) down