FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o server ./cmd/api

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache sqlite-libs
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
