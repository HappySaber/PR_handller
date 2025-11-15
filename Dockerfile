FROM golang:1.25.3 AS builder

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin ./cmd/PR/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin /app/bin

COPY ./internal/migrations/ /docker-entrypoint-initdb.d/

COPY .env .env 

ENV PORT=8080

EXPOSE 8080

CMD ["/app/bin"]


