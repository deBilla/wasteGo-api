package controllers

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWasteItems(c *gin.Context) {
	var wasteItems []models.WasteItem
	if err := configs.DB.Where("user_id = ?", c.Param("userId")).First(&wasteItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wasteItems})
}

func CreateWasteItem(c *gin.Context) {
	var input models.WasteItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wasteItem := models.WasteItem{Name: input.Name, Type: input.Type, Quantity: input.Quantity, UserID: input.UserID}
	configs.DB.Create(&wasteItem)

	c.JSON(http.StatusOK, gin.H{"data": wasteItem})
}

func DeleteWasteItem(c *gin.Context) {
	var wasteItem models.WasteItem
	if err := configs.DB.Where("id = ?", c.Param("id")).First(&wasteItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	configs.DB.Delete(&wasteItem)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
