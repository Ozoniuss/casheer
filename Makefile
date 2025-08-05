.PHONY: run-local stop-local run-local-persistent stop-local-persistent run-ui-only trigger-e2e-tests

run-local:
	docker compose -f docker-compose.yml up -d --build`

stop-local:
	docker compose -f docker-compose.yml down

run-local-persistent:
	docker compose -f docker-compose.yml -f docker-compose.persistent.yml up -d --build

stop-local-persistent:
	docker compose -f docker-compose.yml -f docker-compose.persistent.yml down

run-ui-only:
	go build -o ui/main ui/*.go && ./ui/main

trigger-e2e-tests:
	chmod +x ./scripts/test_e2e.sh
	bash ./scripts/test_e2e.sh
