FROM golang:latest

WORKDIR /app/messaging

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o messaging

EXPOSE 8080

CMD ["./messaging"]