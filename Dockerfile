FROM golang:1.22.1 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main .


EXPOSE 8080


CMD ["/app/main"]
