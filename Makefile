.PHONY: build
build:
	go build -o ./build/server ./cmd/server/main.go
	go build -o ./build/client ./cmd/client/main.go

.PHONY: dc
dc:
	docker-compose up --remove-orphans --build

.PHONY: test
test:
	go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

.PHONY: lint
lint:
	golangci-lint run