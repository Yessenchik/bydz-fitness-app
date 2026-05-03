.PHONY: proto generate test test-cover lint migrate-up migrate-down

proto:
	protoc \
	  -I=proto \
	  --go_out=./gen/userauth/v1 --go_opt=paths=source_relative \
	  --go-grpc_out=./gen/userauth/v1 --go-grpc_opt=paths=source_relative \
	  user_auth.proto

generate:
	go generate ./...

test:
	go test ./... -race -count=1

test-cover:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path internal/infrastructure/postgres/migrations \
	        -database "$$POSTGRES_DSN" up

migrate-down:
	migrate -path internal/infrastructure/postgres/migrations \
	        -database "$$POSTGRES_DSN" down 1