.PHONY: build
build:
	docker compose build

.PHONY: start
start:
	docker compose up -d

.PHONY: stop
stop:
	docker compose stop

.PHONY: attach
attach:
	docker compose exec tcpip-stack bash

.PHONY: status
status:
	docker compose ps
