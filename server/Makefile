GOPATH := $(shell go env GOPATH)

# For go v1.23+
air-install:
	go install github.com/air-verse/air@latest

serve:
	$(GOPATH)/bin/air --build.cmd "go build -o ./bin/app main.go" --build.bin "./bin/app"

test:
	go test ./tests
