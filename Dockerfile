FROM golang:1.16.3-alpine3.13 as builder
EXPOSE 8080

ENV CGO_ENABLED=0

WORKDIR /
COPY . .
RUN go build -o app main.go

FROM scratch
COPY --from=builder /app .

ENTRYPOINT ["./app"]
