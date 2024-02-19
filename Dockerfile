FROM golang:alpine

WORKDIR /app

COPY go.mod  go.sum ./


ENV CGO_ENABLED=1

RUN go mod download

COPY . .

RUN go build main.go

EXPOSE 8000

CMD ["./main"]