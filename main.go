package main

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/controllers"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	router.GET("/wasteItems/:userId", controllers.GetWasteItems)
	router.POST("/wasteItem", controllers.CreateWasteItem)
	router.DELETE("/wasteItem/:id", controllers.DeleteWasteItem)
	router.POST("/wasteItem/uploadImage", uploadImageS3)

	router.Run(":80")
}

func uploadImageS3(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		fmt.Println("Error retrieving file:", err)
		return
	}
	defer file.Close()

	awsRegion := "us-east-1"

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	svc := s3.New(sess)

	bucketName := "wastego"
	objectKey := header.Filename

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
		ACL:    aws.String("public-read"), // Set ACL to allow public read access
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file"})
		fmt.Println("Error uploading file to S3:", err)
		return
	}

	publicURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey)

	responseData, err := DetectS3Label(publicURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect labels"})
		fmt.Println("Error detecting labels:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "labels": responseData, "img_url": publicURL})
}

func DetectS3Label(imagePath string) (any, error) {
	url := "https://api.edenai.run/v2/image/object_detection"
	payload := map[string]interface{}{
		"providers":          "clarifai",
		"file_url":           imagePath,
		"fallback_providers": "",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding payload:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRjNDA3NzUtMmVhOS00MTViLTk1YzEtYjYxYWM4ZWI0YTdkIiwidHlwZSI6ImFwaV90b2tlbiJ9.AaSjKlI6Ay4xwWc102wJifnnlGZqrIeaDHotjUlzIwc")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	return string(body), nil
}
