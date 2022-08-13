PRIVATE_IP := 192.0.2.1

.PHONY: build
build:
	docker compose build

.PHONY: start
start:
	docker compose up -d

.PHONY: stop
stop:
	docker compose stop

.PHONY: clean
clean:
	docker compose down

.PHONY: attach
attach:
	docker compose exec tcpip-stack bash

.PHONY: gen-tap0
gen-tap0:
	docker compose exec tcpip-stack ./bin/gen-tap0.sh

.PHONY: show-tap0
show-tap0:
	docker compose exec tcpip-stack ip addr show tap0

.PHONY: ping
ping:
	docker compose exec tcpip-stack ping ${PRIVATE_IP}

.PHONY: status
status:
	docker compose ps

.PHONY: step0
step0:
	@docker compose exec tcpip-stack go run test/log/main.go

.PHONY: step1
step1:
	@docker compose exec tcpip-stack go run test/dummy/main.go

.PHONY: step2
step2:
	@docker compose exec tcpip-stack go run test/dummy/main.go

.PHONY: step3
step3:
	@docker compose exec tcpip-stack go run test/loopback/without_interface/main.go

.PHONY: step4
step4:
	@docker compose exec tcpip-stack go run test/loopback/without_interface/main.go

.PHONY: step5
step5:
	@docker compose exec tcpip-stack go run test/loopback/without_interface/main.go

.PHONY: step6
step6:
	@docker compose exec tcpip-stack go run test/loopback/without_interface/main.go

.PHONY: step7
step7:
	@docker compose exec tcpip-stack go run test/loopback/with_interface/main.go

.PHONY: step8
step8:
	@docker compose exec tcpip-stack go run test/ip/main.go

.PHONY: step9
step9:
	@docker compose exec tcpip-stack go run test/ip/main.go

.PHONY: step10
step10:
	@docker compose exec tcpip-stack go run test/ip/main.go

.PHONY: step11
step11:
	@docker compose exec tcpip-stack go run test/icmp/main.go
