build:
	docker-compose build todo-app

run:
	docker-compose up todo-app

swag:
	swag init -g cmd/app/main.go