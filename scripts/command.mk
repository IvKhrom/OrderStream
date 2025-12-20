.PHONY: cov
cov:
	go test -cover ./...

.PHONY: mock
mock:
	mockery


