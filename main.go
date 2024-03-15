package main

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/controllers"
	"fmt"
	"io"
	"net/http"
	"os"

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

	router.POST("/wasteItem/uploadImage", uploadImage)

	router.Run("localhost:8080")
}

func uploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		fmt.Println("Error retrieving file:", err)
		return
	}
	defer file.Close()

	// Create a new file in the uploads directory
	out, err := os.Create("uploads/" + header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	// Copy the file to the new destination
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		fmt.Println("Error copying file:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
