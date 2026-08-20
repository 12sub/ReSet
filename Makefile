.PHONY: build run dev stop logs

build:
	docker compose build

run:
	docker compose up -d

dev:
	docker compose up

stop:
	docker compose down

logs:
	docker-compose logs -f api