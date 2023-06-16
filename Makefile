GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=app
DOCKER_IMAGE=releases-monitoring

all: test build
build:
	$(GOMOD) tidy
	$(GOMOD) vendor
	$(GOBUILD) -mod vendor -a -o $(BINARY_NAME) ./cmd/app
test:
	$(GOTEST) -v -count=1 ./...
clean:
	$(GOCLEAN) -cache -modcache # optional
	rm -f $(BINARY_NAME)
run:
	./$(BINARY_NAME)
deps:
	$(GOMOD) tidy
	$(GOMOD) vendor


build-linux:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -mod vendor -a -o $(BINARY_NAME) ./cmd/app
docker-build:
	docker build -t $(DOCKER_IMAGE):test .
	docker run --rm -p 8080:8080 --net host --env-file ./docker/env.list $(DOCKER_IMAGE):test
