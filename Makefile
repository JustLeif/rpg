GOPATH := $(shell go env GOPATH)

# For go v1.23+
air-install:
	go install github.com/air-verse/air@latest

run:
	$(GOPATH)/bin/air --build.cmd "go build -o bin/api main.go" --build.bin "./bin/api"

test:
	go test ./tests
