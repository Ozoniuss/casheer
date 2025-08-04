.PHONY: unit-tests e2e-tests

unit-tests:
	go test -v ./...

e2e-tests:
	chmod +x ./test_e2e.sh
	bash ./test_e2e.sh

run-local:
	docker compose -f docker-compose.yml up -d

run-local-persistent:
	docker compose -f docker-compose.yml -f docker-compose.persistent.yml up -d
