all: build

build:
	CGO_ENABLED="1" go build -o bin/go-gl-cube  ./src

errcheck:
	errcheck ./src

test:
	go test ./src

.PHONY: all build test errcheck
