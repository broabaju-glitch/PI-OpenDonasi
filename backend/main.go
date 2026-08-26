package main

import (
	"log"
	"os"

	"opendonasi-backend/config"
	"opendonasi-backend/routes"
	"opendonasi-backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	r := gin.Default()

	r.Static("/uploads", "./public/uploads")

	utils.SetupTemplates(r)

	r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
	}))

	jwtSecret := []byte("super-secret-opendonasi-key")
	routes.RegisterRoutes(r, db, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("OpenDonasi Backend running on http://localhost:%s", port)
	log.Printf("Admin Panel: http://localhost:%s/admin/login", port)
	r.Run(":" + port)
}