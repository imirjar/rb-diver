FROM golang:1.22.2 
WORKDIR /app
COPY . /app/
# COPY db/reports /db
# COPY bin/diver .
# COPY config/diver .
RUN CGO_ENABLED=0 GOOS=linux go build -o diver ./cmd/main.go
CMD ["./diver"]