package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nerokome/econo/controllers"
	"github.com/nerokome/econo/database"
	"github.com/nerokome/econo/routes"
)

func main() {
	// Load port from env or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Initialize MongoDB client
	dbClient := database.DBSet()
	productCollection := database.ProductData(dbClient, "Products")
	userCollection := database.UserData(dbClient, "Users")

	// Initialize controllers with DB collections
	app := controllers.NewApplication(productCollection, userCollection)

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with allowed origins in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	routes.RegisterRoutes(router, app)

	// Start server with graceful shutdown
	srvAddr := ":" + port
	server := &http.Server{
		Addr:    srvAddr,
		Handler: router,
	}

	go func() {
		log.Printf("Server running on %s", srvAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
