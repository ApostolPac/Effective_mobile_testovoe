FROM golang:1.23.8-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o applic ./cmd/app/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/applic .

COPY --from=builder /app/internal/migrations internal/migrations

CMD [ "./applic" ]