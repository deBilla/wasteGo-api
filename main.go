package main

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	configs.ConnectDatabase()

	router.GET("/wasteItems", controllers.GetWasteItems)
	router.POST("/wasteItem", controllers.CreateWasteItem)

	router.Run("localhost:8080")
}
