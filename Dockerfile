FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o chat-service .
EXPOSE 10201
CMD ["./chat-service"]
