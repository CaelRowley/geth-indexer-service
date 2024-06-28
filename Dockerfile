FROM golang:1.22.3

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN  go build -o main cmd/main.go

EXPOSE 8080

CMD ["./main", "-sync"]
