build:
	docker-compose build todo-rest-api

run:
	docker-compose up todo-rest-api

swag:
	swag init -g cmd/app/main.go