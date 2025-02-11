DEFAULT_GOAL: build

.PHONY: fmt vet build clean run imports docker_up docker_down

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
	rm -rf api

docker_up:
	docker compose up -d

docker_down:
	docker down
