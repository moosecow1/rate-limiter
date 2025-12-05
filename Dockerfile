FROM golang:1.25-alpine AS builder

WORKDIR /app

# Go modules
COPY go.mod go.sum ./
RUN go mod download

# Resources
ADD quotes.json /app/quotes.json

# App files
COPY . .

# Build
RUN go build -o api ./cmd/api 

# Alpie
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/api .

# Expose API port
EXPOSE 8080

# Other
COPY quotes.json .


CMD ["./api"]