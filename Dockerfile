FROM golang:latest

WORKDIR /app

COPY main.go .

RUN go get github.com/prometheus/client_golang/prometheus

CMD ["go", "run", "main.go"]
