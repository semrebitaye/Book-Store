package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"log"

	"github.com/gin-gonic/gin"
)

var categoryRequest struct {
	Name string `json:"name"`
}

func CreateCategory(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&categoryRequest)
	if err != nil {
		c.Status(400)
		return
	}
	// create a post
	category := models.Category{Name: categoryRequest.Name}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		c.Status(500)
		return
	}

	// return it
	c.JSON(200, gin.H{"Category": category})
}

func GetCategory(c *gin.Context) {
	// get the user
	var categories []models.Category
	result := initializers.DB.Find(&categories)
	if result.Error != nil {
		log.Fatal("Failed to get the user")
	}
	// respond with them
	c.IndentedJSON(200, gin.H{"Category": categories})
}

func GetCategoryById(c *gin.Context) {
	// get id of url
	category_id := c.Param("category_id")

	// get the user
	var category models.Category
	result := initializers.DB.First(&category, category_id)
	if result.Error != nil {
		log.Fatal("Failed to get the user")
	}
	// respond with them
	c.IndentedJSON(200, gin.H{"Category": category})
}

func UpdateCategory(c *gin.Context) {
	// get the id of the url
	category_id := c.Param("category_id")

	//get the data of the req body
	c.Bind(&categoryRequest)

	//fined the user where updating
	var category models.User
	result := initializers.DB.First(&category, category_id)
	if result.Error != nil {
		log.Fatal("Failed to get the user")
	}

	// update it
	initializers.DB.Model(&category).Updates(models.Category{Name: categoryRequest.Name})

	// respond it
	c.JSON(200, gin.H{"User": category})
}

func DeleteCategory(c *gin.Context) {
	// get the url of the body
	category_id := c.Param("category_id")

	// delete it
	initializers.DB.Delete(&models.Category{}, category_id)

	// respond it
	c.Status(200)
}
