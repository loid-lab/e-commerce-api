FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download


COPY . .


RUN go build -o server main.go

# Final stage
FROM alpine:latest
WORKDIR /app


RUN apk --no-cache add ca-certificates


COPY --from=builder /app/server .


EXPOSE 8080
CMD ["./server"]