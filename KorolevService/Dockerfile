FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

VOLUME /app/logs:/logs
VOLUME /app/data:/data
VOLUME /app/KorolevService:/KorolevService
VOLUME /app/directories:/app/directories

CMD ["go", "run", "main.go"]
