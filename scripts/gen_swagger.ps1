Param(
  [string]$GoImage = "golang:1.23.4-alpine"
)

$ErrorActionPreference = "Stop"

# Генерация Swagger/OpenAPI из proto через protoc + protoc-gen-openapiv2.
# Требование: swagger.json лежит в services/api_service/internal/api/swagger/orders.swagger.json
# 1) Пытаемся через Docker (без установки protoc/buf локально)
# 2) Если Docker недоступен — пытаемся через локальные protoc + protoc-gen-openapiv2 (если установлены)

Write-Host "Generating Swagger via protoc-gen-openapiv2..."

$root = Split-Path -Parent $PSScriptRoot

function Test-Docker {
  try {
    docker version | Out-Null
    return $true
  } catch {
    return $false
  }
}

if (Test-Docker) {
  Write-Host "Docker detected, generating in container..."
  docker run --rm -v "${root}:/workspace" -w /workspace $GoImage sh -c `
    "apk add --no-cache protobuf protobuf-dev git >/dev/null && \
     /usr/local/go/bin/go env -w GOPROXY=https://proxy.golang.org,direct && \
     /usr/local/go/bin/go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.3 && \
     protoc -I services/api_service/api -I /usr/include \
       --openapiv2_out=services/api_service/internal/api/swagger \
       --openapiv2_opt=logtostderr=true,allow_merge=true,merge_file_name=orders \
       services/api_service/api/orders_api/orders.proto"
  Write-Host "Done. Swagger: services/api_service/internal/api/swagger/orders.swagger.json"
  exit 0
}

try {
  protoc --version | Out-Null
} catch {
  throw "Neither Docker nor local protoc found. Install protoc + protoc-gen-openapiv2 or start Docker Desktop."
}

Push-Location $root
try {
  protoc -I services/api_service/api `
    --openapiv2_out=services/api_service/internal/api/swagger `
    --openapiv2_opt=logtostderr=true,allow_merge=true,merge_file_name=orders `
    services/api_service/api/orders_api/orders.proto
  Write-Host "Done. Swagger: services/api_service/internal/api/swagger/orders.swagger.json"
} finally {
  Pop-Location
}


