.PHONY: run-api
run-api:
	cd services/api_service/cmd/api_service && go run .

.PHONY: run-worker
run-worker:
	cd services/worker/cmd/worker && go run .

.PHONY: swagger
swagger:
	docker run --rm -v "$(CURDIR):/workspace" -w /workspace golang:1.24.11-alpine sh -c "\
		apk add --no-cache protobuf protobuf-dev git >/dev/null && \
		/usr/local/go/bin/go env -w GOPROXY=https://proxy.golang.org,direct && \
		/usr/local/go/bin/go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.3 && \
		protoc -I services/api_service/api -I /usr/include \
			--openapiv2_out=services/api_service/internal/pb/swagger/orders_api \
			--openapiv2_opt=logtostderr=true,allow_merge=true,merge_file_name=orders \
			services/api_service/api/orders_api/orders.proto \
	"
