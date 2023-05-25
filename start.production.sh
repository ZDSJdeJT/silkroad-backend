#!/bin/sh
export APP_ENV=production
export CGO_ENABLED=1

go run main.go