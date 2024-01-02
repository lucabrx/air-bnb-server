# Build stage
FROM golang:latest as builder
WORKDIR /app
COPY app.env app.env
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .


EXPOSE 8080
CMD ["app/main"]

