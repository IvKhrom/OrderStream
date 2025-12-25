module github.com/ivkhr/orderstream/services/api_service

go 1.23.0

toolchain go1.23.4

require (
	github.com/ivkhr/orderstream/shared v0.0.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/google/uuid v1.6.0
)

replace github.com/ivkhr/orderstream/shared => ../../shared


