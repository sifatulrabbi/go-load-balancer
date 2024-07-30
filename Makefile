run:
	go run ./main.go
test-servers:
	go run ./tests/testservers/main.go

.PHONY: run, run-test-servers
