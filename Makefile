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

.PHONY: step0
step0:
	@docker compose exec tcpip-stack go run test/step0/main.go

.PHONY: step1
step1:
	@docker compose exec tcpip-stack go run test/step1/main.go
