# Quiz Backend (AWS Lambda + Twirp + DynamoDB)

This is a starter template for your quiz / rewards backend.

Tech:
- Go 1.21
- Twirp (RPC framework)
- AWS Lambda + API Gateway
- DynamoDB

## Layout

- `cmd/server/main.go` – Lambda entry + local HTTP server
- `proto/*.proto` – Twirp service definitions
- `internal/models` – domain models
- `internal/repository` – interfaces + DynamoDB stubs
- `internal/service` – business logic layer (skeleton)
- `internal/handlers` – Twirp handlers that call services

## Generate Twirp Code

You must install `protoc` and the Twirp Go plugin locally, then run:

```bash
protoc   --go_out=./rpc --go_opt=paths=source_relative   --twirp_out=./rpc --twirp_opt=paths=source_relative   proto/auth.proto proto/user.proto proto/question.proto proto/reward.proto
```

That will create Go packages under `rpc/` which are imported by the handlers and `main.go`.

## Run locally

```bash
go mod tidy
go run ./cmd/server
```

Server listens on `:8080`.

## Deploy to Lambda

Build for Linux and upload the binary, then wire it behind API Gateway (HTTP API):

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/server
zip function.zip bootstrap
```

Create a Lambda with runtime "provided.al2023" (or use container image) and upload `function.zip`.

You'll then configure API Gateway to proxy requests to Lambda.

## Next steps

- Implement real logic in `internal/service/*`
- Fill DynamoDB code in `internal/repository/dynamo_*`
- Add authentication (JWT) and admin-only endpoints
