FROM golang:1.22.2 
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o diver ./cmd/main.go
CMD ["./diver"]