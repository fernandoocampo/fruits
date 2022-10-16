# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
SRC_FOLDER=cmd/fruitsd
BINARY_NAME=bin/fruits
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
DOCKER_REPO=fdocampo
DOCKER_CONTAINER=frutal

all: build build-linux

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v ./$(SRC_FOLDER)

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

tidy:
	$(GOCMD) mod tidy


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./$(SRC_FOLDER)
docker-build:
	DOCKER_BUILDKIT=0 docker build --no-cache -t $(DOCKER_REPO)/$(DOCKER_CONTAINER) .
docker-push:
	docker push $(DOCKER_REPO)/$(DOCKER_CONTAINER)
lint-docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.48-alpine golangci-lint run
run-local:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./$(SRC_FOLDER)
	docker-compose up --build -d
run-docker-local:
	docker run --rm -it -p 8080:8080 vivekteam/fruits
clean-local:
	docker-compose down
helm-lint:
	helm lint k8s-v2/fruits
helm-install:
	helm install --name fruits ./k8s-v2/fruits
create-topic:
	aws sns create-topic --name fruits --endpoint-url http://localhost:4566 --region us-east-1
list-sns:
	aws sns list-topics --endpoint-url http://localhost:4566 --region us-east-1
test:
	go test -race ./...