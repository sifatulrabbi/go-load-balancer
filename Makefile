run:
	go run ./main.go -config=./tests/sample-config.json

test-servers:
	go run ./tests/testservers/main.go

build:
	mkdir -p ./build
	go build -o ./build/glb ./main.go

list:
	GOPROXY=proxy.golang.org go list -m github.com/sifatulrabbi/go-load-blancer@v0.1.0-beta.1

.PHONY: run, test-servers, build, list
