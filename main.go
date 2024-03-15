package main

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}

	// Use CORS middleware
	router.Use(cors.New(config))

	configs.ConnectDatabase()

	router.GET("/wasteItems", controllers.GetWasteItems)
	router.POST("/wasteItem", controllers.CreateWasteItem)
	router.DELETE("/wasteItem/:id", controllers.DeleteWasteItem)

	router.Run("localhost:8080")
}
