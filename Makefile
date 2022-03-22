.PHONY:
.SILENT:

build:
	go build -o bin/bot cmd/bot/main.go

run: build
	./bin/bot

build-container:
	docker build -t telegram-bot:v0.1 .

run-container:
	docker run -p 80:80 --env-file .env telegram-bot:v0.1
