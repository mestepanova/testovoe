gen-api:
	go tool oapi-codegen \
		-generate gin-server \
		-package api ./api/openapi.yml > ./internal/generated/api/server.gen.go
	go tool oapi-codegen \
		-generate types \
		-package api ./api/openapi.yml > ./internal/generated/api/types.gen.go
	go tool oapi-codegen \
		-generate client \
		-package api ./api/openapi.yml > ./internal/generated/api/client.gen.go

test-int:
	go test -tags=integration -race -v ./tests/...

test-unit:
	go test -race -v ./...

lint:
	go tool golangci-lint run

fmt:
	go tool golangci-lint run --fix
	go fmt ./...