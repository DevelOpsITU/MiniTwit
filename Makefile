.PHONY: build clean run_fresh run

BINARY_NAME=gominitwit
CONTAINER_NAME=Minitwit

VERSION?=0.0.0
SERVICE_PORT?=8080
DOCKER_REGISTRY?=groupddevops/ #if set it should finished by /
COMMIT_SHA=$(git rev-parse --short HEAD)


#all: run_fresh
all: help

build:
	./scripts/build_app.sh ${BINARY_NAME}

clean:
	go clean
	rm -f out/${BINARY_NAME}

clean-all: clean
	rm -f out/*

## Run
run: ## Run go application
	go run main.go



run_fresh: ## Run go application with fresh (Auto-reloading)
	fresh -c ./fresh/my_fresh_runner.conf



## Linters
docker_lint: ## Run docker linting script
	./scripts/docker-lint.sh


## Tests
test: ## Run Go tests
	 go test -v ./...

test_coverage: ## Run Go tests with coverage
	go test ./main.go -coverprofile=coverage.out

go_lint: ## Lint all go files
	 golint  ./...


deps: ## Install dependencies
	go mod tidy
	go get -u golang.org/x/lint/golint
	go get github.com/pilu/fresh
	go install github.com/pilu/fresh
	go install gorm.io/gorm
	go install gorm.io/driver/sqlite
	go install gorm.io/driver/postgres
	go install golang.org/x/lint/golint


# From https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705 ,
# https://golangdocs.com/makefiles-golang and

## Docker:
docker-build: ## Use the dockerfile to build the image
	docker build --rm --tag $(BINARY_NAME):latest .

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

docker-scan: ## Scan the image built
	sudo apt-get update && apt-get install docker-scan-plugin -y
	docker scan $(BINARY_NAME):latest

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