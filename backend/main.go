package main

import (
	"log"
	"student-dormitory-management/config"
	"student-dormitory-management/database"
	"student-dormitory-management/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	database.InitDB()

	r := gin.Default()

	routers.InitRoutes(r)

	port := config.AppConfig.Server.Port
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
