.PHONY: run test build test-all

run:
	go run ./cmd/securetalon

build:
	go build -o securetalon.exe ./cmd/securetalon

test:
	go test ./...

# Backend tests + UI build + screenshot capture. Run: .\scripts\test-all.ps1  (from repo root)
test-all: test
	cd ui && npm run build
