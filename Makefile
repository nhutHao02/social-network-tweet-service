run:
	go run main.go --config=./config/local/config.yaml
wire:
	wire ./internal
migrate-file:
	migrate create -ext sql -dir migrations -seq ${filename}