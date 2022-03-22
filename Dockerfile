FROM golang:1.18-alpine3.15 as builder

COPY . /github.com/kostyawhite/telegram-bot/
WORKDIR /github.com/kostyawhite/telegram-bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/kostyawhite/telegram-bot/bin/bot .
COPY --from=0 /github.com/kostyawhite/telegram-bot/configs configs/

EXPOSE 80

CMD ["./bot"]
