.PHONY: run test build

run:
	go run ./cmd/securetalon

build:
	go build -o securetalon.exe ./cmd/securetalon

test:
	go test ./...
