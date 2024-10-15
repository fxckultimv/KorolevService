FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o korolev_service main.go

CMD ["./korolev_service"]
