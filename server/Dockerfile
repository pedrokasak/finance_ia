FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build for linux amd64
RUN GOARCH=amd64 GOOS=linux go build -o main .

# Stage 2: runtime
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]