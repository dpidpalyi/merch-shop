DEFAULT_GOAL: build

.PHONY: fmt vet build clean run imports docker_up docker_down test test-report

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: build
build: vet
	go build ./cmd/app

.PHONY: run
run: vet
	go run ./cmd/app

.PHONY: imports
imports:
	find . -name \*.go -exec goimports -w -l {} \;

.PHONY: clean
clean:
	rm -rf app coverage.out

.PHONY: docker_up
docker_up:
	docker compose -f docker-compose.yaml up -d

.PHONY: docker_down
docker_down:
	docker compose -f docker-compose.yaml down

.PHONY: docker_down_with_volumes
docker_down_with_volumes:
	docker compose -f docker-compose.yaml down -v

.PHONY: test
test:
	go test ./... -v --cover

.PHONY: test-report
test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: docker_test_up
docker_test_up:
	docker compose -f docker-compose.test.yaml up -d -V

.PHONY: docker_test_down
docker_test_down:
	docker compose -f docker-compose.test.yaml down -v

.PHONY: e2e
e2e: docker_test_down docker_test_up
	sleep 3
	go test -v -count=1 -tags=e2e ./...
	make docker_test_down
