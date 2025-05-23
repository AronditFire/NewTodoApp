build:
	docker-compose build todo-app

run:
	docker-compose up todo-app

swag:
	swag init -g cmd/app/main.go

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

test:
	go test -v -count=1 ./...