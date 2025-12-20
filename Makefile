include ./scripts/command.mk

.PHONY: run-api
run-api:
	cd cmd/api_service && go run .

.PHONY: run-worker
run-worker:
	cd cmd/worker && go run .
