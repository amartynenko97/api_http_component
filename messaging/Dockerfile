FROM golang:1.16

WORKDIR /app

COPY . .

RUN go build -o messaging-app .

EXPOSE 8080

CMD ["./messaging-app"]