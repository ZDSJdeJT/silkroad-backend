SET CGO_ENABLED=1

swag init -g main.go --output docs

go run -v main.go
