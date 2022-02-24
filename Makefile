.PHONY: build clean run_fresh run

BINARY_NAME=group_d_go_app
CONTAINER_NAME=Minitwit-container

#all: run_fresh
all: help

build:
	./build_app.sh ${BINARY_NAME}

clean:
	go clean
	rm -f out/${BINARY_NAME}

clean-all: clean
	rm -f out/*

## Tests
test: ## Run Go tests (Not implemented)
	go test src/minitwit.go

test_coverage: ## Run Go tests with coverage (Not implemented)
	go test src/minitwit.go -coverprofile=coverage.out

## Install dependencies
deps:
	go install github.com/pilu/fresh
	go install github.com/mattn/go-sqlite3

dep:
	go mod download

# From https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705 ,
# https://golangdocs.com/makefiles-golang and

## Docker:
docker-build: ## Use the dockerfile to build the container
	docker build --rm --tag $(BINARY_NAME) .

docker-release: ## Release the container with tag latest and version
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

docker-run: docker-build ## Build and run the container locally with port 8080
	docker rm -f $(CONTAINER_NAME)
	docker run -d -p 8080:8080 --name=$(CONTAINER_NAME) $(BINARY_NAME)
	docker ps -l
	docker logs Minitwit-container

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make ${GREEN}<target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}${RESET}%s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)