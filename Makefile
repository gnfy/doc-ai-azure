bin_build:
	go build -ldflags="-s -w" -o ./app/main ./cmd/server/main.go