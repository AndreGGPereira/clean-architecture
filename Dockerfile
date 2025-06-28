FROM golang:1.23.4

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/cmd/ordersystem

CMD ["go", "run", "main.go", "wire_gen.go"]