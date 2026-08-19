FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with embedded templates
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/reset ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/reset .

EXPOSE 8080
CMD ["./reset"]