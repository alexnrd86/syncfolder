run:
	@go run ./cmd/app/main.go

runtest:
	@go test -v ./internal/synch
	@go test -v ./internal/logger
	@go test -v ./internal/utils

bench:
	@go test -bench=. -benchmem -benchtime 80x ./internal/synch 