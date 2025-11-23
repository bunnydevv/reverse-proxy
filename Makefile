.PHONY: build run clean test fmt vet

BINARY_NAME=reverse-proxy
CONFIG_FILE=config.yaml

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME) -config $(CONFIG_FILE)

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t reverse-proxy:latest .

docker-run:
	docker run -p 8080:8080 -v $(PWD)/config.yaml:/app/config.yaml reverse-proxy:latest
