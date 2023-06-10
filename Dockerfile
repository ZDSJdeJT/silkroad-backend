FROM golang:1.20

WORKDIR /app
COPY . .

RUN go build -o bin/server main.go

ENV APP_MODE=production

CMD ["./bin/server"]
