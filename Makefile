DEFAULT_GOAL: build

.PHONY: fmt vet build clean run imports docker_up docker_down test test-report

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build ./cmd/app

run: vet
	go run ./cmd/app

imports:
	find . -name \*.go -exec goimports -w -l {} \;

clean:
	rm -rf api coverage.out

docker_up:
	docker compose up -d

docker_down:
	docker down

test:
	go test ./... -v --cover

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out
