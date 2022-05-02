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

.PHONY: step2
step2:
	@docker compose exec tcpip-stack go run test/step2/main.go

.PHONY: step3
step3:
	@docker compose exec tcpip-stack go run test/step3/main.go

.PHONY: step4
step4:
	@docker compose exec tcpip-stack go run test/step4/main.go
