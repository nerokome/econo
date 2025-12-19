package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nerokome/econo/controllers"
	"github.com/nerokome/econo/database"
	"github.com/nerokome/econo/routes"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	
	client := database.DBSet()

	
	app := controllers.NewApplication(
		database.Collection(client, "products"),
		database.Collection(client, "users"),
	)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	
	routes.UserRoutes(router, app)

	log.Fatal(router.Run(":" + port))
}
