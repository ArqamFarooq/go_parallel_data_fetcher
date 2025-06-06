# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o fetcher main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fetcher .
EXPOSE 8080
CMD ["./fetcher"]