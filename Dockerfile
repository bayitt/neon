FROM golang:1.23-alpine

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./neon

EXPOSE 8080

CMD ["./neon"]