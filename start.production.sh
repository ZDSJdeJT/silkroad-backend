#!/bin/sh
export APP_MODE=production
export CGO_ENABLED=1

go run main.go