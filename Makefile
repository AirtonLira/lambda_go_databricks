build:
	go build -o ./cmd/main ./cmd/main.go 
	cd ./infra && terraform apply --auto-approve