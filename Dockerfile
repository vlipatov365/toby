FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o bot .

ENV TELEGRAM_BOT_TOKEN=''

CMD ["sh", "-c", "./bot"]