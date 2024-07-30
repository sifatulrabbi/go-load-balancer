run:
	go run ./main.go -config=./tests/sample-config.json
test-servers:
	go run ./tests/testservers/main.go

.PHONY: run, run-test-servers
