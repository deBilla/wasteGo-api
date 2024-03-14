package controllers

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWasteItems(c *gin.Context) {
	var wasteItems []models.WasteItem
	configs.DB.Find(&wasteItems)

	c.JSON(http.StatusOK, gin.H{"data": wasteItems})
}

func CreateWasteItem(c *gin.Context) {
	var input models.WasteItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := models.WasteItem{ID: input.ID, Name: input.Name, Type: input.Type}
	configs.DB.Create(&customer)

	c.JSON(http.StatusOK, gin.H{"data": customer})
}
