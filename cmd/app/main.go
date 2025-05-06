package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	server "github.com/AronditFire/todo-app"
	db "github.com/AronditFire/todo-app/internal/database"
	"github.com/AronditFire/todo-app/internal/handlers"
	"github.com/AronditFire/todo-app/internal/repository"
	"github.com/AronditFire/todo-app/internal/service"
	"github.com/joho/godotenv"
)

// @title Swagger Example API
// @version 1.0
// @description This is a my first rest api.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}

	database, err := db.InitDB() // connect to DB
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	repo := repository.NewRepository(database)
	srv := service.NewService(repo)
	handler := handlers.NewHander(srv)

	server := new(server.Server)
	go func() {
		log.Printf("Starting server at port: %s", string(os.Getenv("PORT")))
		if err := server.Run(string(os.Getenv("PORT")), handler.InitRoutes()); err != nil {
			log.Fatalf("Could not run server : %v", err)
		}
	}()

	// Ловим сигнал остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown the server: %v", err)
	}
	log.Println("HTTP server stopped")

	if err := db.CloseConnection(database); err != nil {
		log.Fatalf("Failed to close connection with database: %v", err)
	}

	log.Println("Database connection successfully stopped")
}
