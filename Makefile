BIN_NAME =kubedemo
PROJECT_NAME =kubedemo
VERSION :=latest
DOCKER_REPO=netyazilim

ifeq ($(OS), Windows_NT)
	PREFIX = env
	EXECUTABLE := $(BIN_NAME).exe
else
	PREFIX = 
	EXECUTABLE := BIN_NAME
endif

.PHONY: build
build:
	@echo "  >  Building binary..."
	statik -src=web
	$(PREFIX) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_NAME) ./cmd

.PHONY: clean
clean:
	rm -f $(BIN_NAME)
	rm -f statik/*

image:
	@echo "  >  Building docker..."
	statik -src=web
	docker build -t $(DOCKER_REPO)/$(PROJECT_NAME):$(VERSION) . -f Dockerfile

image-alpine:
	@echo "  >  Building docker..."
	docker build -t $(DOCKER_REPO)/$(PROJECT_NAME):alpine . -f alpine.Dockerfile

image-dev:
	@echo "  >  Building docker..."
	statik -src=web
	$(PREFIX) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_NAME) ./cmd
	docker build -t $(DOCKER_REPO)/$(PROJECT_NAME):dev . -f dev.Dockerfile
