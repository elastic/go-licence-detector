.PHONY: generate
generate:
	@ go generate ./...

.PHONY: build
build: generate
	@ go build -o bin/go-licence-detector

