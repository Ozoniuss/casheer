.PHONY: run-dev run-dev-backend stop-dev-backend run-ui-only trigger-e2e-tests

run-dev: run-dev-backend run-ui-only

install-dev:
	go get -tool github.com/Ozoniuss/genconfig

run-dev-backend:
	docker compose -f docker-compose.yml -f docker-compose.persistent.yml up -d --build

stop-dev-backend:
	docker compose -f docker-compose.yml -f docker-compose.persistent.yml down

run-ui-only:
	cd ./ui && go build -o main *.go && ./main

trigger-e2e-tests:
	chmod +x ./scripts/test_e2e.sh
	bash ./scripts/test_e2e.sh

trigger-e2e-tests-build:
	chmod +x ./scripts/test_e2e.sh
	bash ./scripts/test_e2e.sh --build

# todo: make to create new db