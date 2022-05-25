build-prod:
	docker build -t openslides-search-service -f Dockerfile .

build-dev:
	docker build -t openslides-search-service-dev -f Dockerfile.dev .

run-dev: | build-dev
	docker-compose -f docker-compose.dev.yml up
	stop-dev

run-pre-test: | build-dev
	docker-compose -f docker-compose.dev.yml up -d
	docker-compose -f docker-compose.dev.yml exec -T search-service ./wait-for.sh redis:6379
	docker-compose -f docker-compose.dev.yml exec -T search-service ./wait-for.sh search-service:9022

run-bash: | run-pre-test
	docker-compose -f docker-compose.dev.yml exec search-service sh
	docker-compose -f docker-compose.dev.yml down

run-check-lint:
	docker-compose -f docker-compose.dev.yml exec -T search-service npm run lint-check

run-check-prettify:
	docker-compose -f docker-compose.dev.yml exec -T search-service npm run prettify-check

run-test: | run-pre-test
	@echo "########################################################################"
	@echo "###################### Start full system tests #########################"
	@echo "########################################################################"
	docker-compose -f docker-compose.dev.yml exec -T search-service npm run test

run-cleanup: | build-dev
	docker-compose -f docker-compose.dev.yml up -d
	docker-compose -f docker-compose.dev.yml exec search-service ./wait-for.sh search-service:9022
	docker-compose -f docker-compose.dev.yml exec search-service npm run cleanup
	docker-compose -f docker-compose.dev.yml down

run-test-and-stop: | run-test
	stop-dev

run-test-prod: | build-prod
	docker-compose -f .github/startup-test/docker-compose.yml up -d
	docker-compose -f .github/startup-test/docker-compose.yml exec -T search-service ./wait-for.sh search-service:9022
	docker-compose -f .github/startup-test/docker-compose.yml down

stop-dev:
	docker-compose -f docker-compose.dev.yml down