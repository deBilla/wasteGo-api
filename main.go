package main

import (
	"billacode/wasteGo/configs"
	"billacode/wasteGo/controllers"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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

	router.GET("/wasteItems", controllers.GetWasteItems)
	router.POST("/wasteItem", controllers.CreateWasteItem)
	router.DELETE("/wasteItem/:id", controllers.DeleteWasteItem)

	router.POST("/wasteItem/uploadImage", uploadImage)
	router.POST("/wasteItem/upload-image-s3", uploadImageS3)

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

	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	bucketName := "your-bucket-name"
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

	c.JSON(http.StatusOK, gin.H{"url": publicURL})
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

	responseData, err := DetectLabels(out.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect labels"})
		fmt.Println("Error detecting labels:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "labels": string(responseData)})
}

func DetectLabels(imagePath string) ([]byte, error) {
	url := "https://api.edenai.run/v2/image/object_detection"
	apiKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRjNDA3NzUtMmVhOS00MTViLTk1YzEtYjYxYWM4ZWI0YTdkIiwidHlwZSI6ImFwaV90b2tlbiJ9.AaSjKlI6Ay4xwWc102wJifnnlGZqrIeaDHotjUlzIwc" // Replace YOUR_API_KEY with your actual API key

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a new form-data body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the image file to the form-data request
	part, err := writer.CreateFormFile("file", imagePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file to form:", err)
	}

	// Add JSON payload to form-data request
	_ = writer.WriteField("providers", "clarifai")
	_ = writer.WriteField("fallback_providers", "")

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer res.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	// Print the result
	// fmt.Println(string(responseData))
	return responseData, nil
}
