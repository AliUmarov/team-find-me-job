
run:
	go run ./cmd/find_me_job

dev:
	air

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...