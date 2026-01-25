# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /category-api main.go

# Final stage
FROM alpine:edge

WORKDIR /

COPY --from=builder /category-api /category-api

EXPOSE 8080

CMD ["/category-api"]
