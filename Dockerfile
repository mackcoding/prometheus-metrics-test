FROM golang:alpine

WORKDIR /app

COPY go/main.go .

RUN go mod init app && \
    go mod tidy && \
    go build -o main main.go

CMD ["./main"]
