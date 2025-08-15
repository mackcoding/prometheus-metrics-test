FROM golang:alpine

WORKDIR /app

COPY go/go.mod go/go.sum ./
RUN go mod download
COPY go/main.go ./
RUN go build -o main main.go
CMD ["./main"]
