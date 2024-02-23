package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"log"

	"github.com/gin-gonic/gin"
)

var authorRequest struct {
	Name        string `json:"name"`
	Biography   string `json:"biography"`
	Nationality string `json:"nationality"`
}

func CreateAuthor(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&authorRequest)
	if err != nil {
		c.Status(400)
		return
	}
	// create a post
	author := models.Author{Name: authorRequest.Name, Biography: authorRequest.Biography, Nationality: authorRequest.Nationality}

	result := initializers.DB.Create(&author)

	if result.Error != nil {
		c.Status(500)
		return
	}
	// return it
	c.JSON(200, gin.H{"Author": author})
}

func GetAuthors(c *gin.Context) {
	// get the user
	var authors []models.Author
	result := initializers.DB.Find(&authors)
	if result.Error != nil {
		log.Fatal("Failed to get the author")
	}
	// respond with them
	c.IndentedJSON(200, gin.H{"Authors": authors})
}

func GetAuthorById(c *gin.Context) {
	// get id of url
	author_id := c.Param("author_id")

	// get the user
	var author models.Author
	result := initializers.DB.First(&author, author_id)
	if result.Error != nil {
		log.Fatal("Failed to get the author")
	}
	// respond with them
	c.IndentedJSON(200, gin.H{"Authors": author})
}

func UpdateAuthor(c *gin.Context) {
	// get the id of the url
	author_id := c.Param("author_id")

	//get the data of the req body
	c.Bind(&authorRequest)

	//fined the user where updating
	var author models.Author
	result := initializers.DB.First(&author, author_id)
	if result.Error != nil {
		log.Fatal("Failed to get the auther")
	}

	// update it
	initializers.DB.Model(&author).Updates(models.Author{Name: authorRequest.Name, Biography: authorRequest.Biography, Nationality: authorRequest.Nationality})

	// respond it
	c.JSON(200, gin.H{"Author": author})
}

func DeleteAuthor(c *gin.Context) {
	// get the url of the body
	author_id := c.Param("author_id")

	// delete it
	initializers.DB.Delete(&models.Author{}, author_id)

	// respond it
	c.Status(200)
}
