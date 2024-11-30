FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates


WORKDIR /root/


COPY --from=builder /app/main .


EXPOSE 8080


CMD ["./main"]