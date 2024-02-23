package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var bookRequest struct {
	Title           string  `json:"title"`
	PublicationDate string  `json:"publication_date"`
	Price           float64 `json:"price"`
	Quantity        uint    `json:"quantity"`
	Cover           string  `json:"cover"`
	AuthorID        uint    `json:"author_id"`
	CategoryID      uint    `json:"category_id"`
}

func CreateBook(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&bookRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind body"})
		return
	}
	// create a post
	u := c.GetUint("user_id")

	if u == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user id"})
		return
	}

	book := models.Book{UserID: uint(u), Title: bookRequest.Title, PublicationDate: bookRequest.PublicationDate, Price: bookRequest.Price, Quantity: bookRequest.Quantity, AuthorID: bookRequest.AuthorID, CategoryID: bookRequest.CategoryID}

	result := initializers.DB.Create(&book)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": result.Error})
		return
	}

	// return it
	c.JSON(200, gin.H{"User": book})
}

func GetBooks(c *gin.Context) {
	// get the book
	var books []models.Book
	result := initializers.DB.Find(&books)
	if result.Error != nil {
		log.Fatal("Failed to get the books")
		return
	}

	// for pagination and sorting books
	// parse pagination parameteres
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// parse sorting criteria
	sortBy := c.DefaultQuery("sort_by", "title")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	//construct database querry
	dbOffset := (page - 1) * pageSize
	err := initializers.DB.Order(sortBy + " " + sortOrder).Offset(dbOffset).Limit(pageSize).Find(&books)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the book"})
	}

	// respond
	c.IndentedJSON(200, books)
}

func GetBookByID(c *gin.Context) {
	// get id of url
	book_id := c.Param("book_id")

	//get the book
	var book models.Book
	initializers.DB.First(&book, book_id)

	// respond
	c.JSON(200, gin.H{"Book": book})

}

func UpdateBook(c *gin.Context) {
	// get the id of the url
	book_id := c.Param("book_id")

	//get the data of the req body
	c.Bind(&bookRequest)

	//fined the user where updating
	var book models.User
	result := initializers.DB.First(&book, book_id)
	if result.Error != nil {
		log.Fatal("Failed to get the user")
	}

	// update it
	initializers.DB.Model(&book).Updates(models.Book{Title: bookRequest.Title, PublicationDate: bookRequest.PublicationDate, Price: bookRequest.Price, Quantity: bookRequest.Quantity})

	// respond it
	c.JSON(200, gin.H{"Book": book})
}

func DeleteBook(c *gin.Context) {
	// get the url of the body
	book_id := c.Param("book_id")

	// delete it
	initializers.DB.Delete(&models.Book{}, book_id)

	// respond it
	c.Status(200)
}

func UploadBookCover(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Eroor": err.Error()})
		return
	}
	defer file.Close()

	// check if the file format is supported
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format"})
		return
	}

	// save the file to the server
	filename := uuid.NewString() + header.Filename
	destination := "./image/" + filename

	if err := c.SaveUploadedFile(header, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save file"})
		return
	}

	// save the file path to the database
	book := models.Book{
		Title:     c.PostForm("title"),
		CoverPath: destination,
	}
	initializers.DB.Create(&book)

	c.JSON(http.StatusOK, gin.H{"message": "Image Uploaded successfully", "book": book})
}

func GetBookCoverImage(c *gin.Context) {
	// get the book id from request url parameter
	book_id := c.Param("book_id")

	// querry database for books based on ids
	var books []models.Book
	initializers.DB.Where("id IN (?)", book_id).Find(&books)

	// create a map to store book id and associated book cover image path
	coverImages := make(map[uint]string)
	for _, book := range books {
		coverImages[book.ID] = book.CoverPath
	}

	// return cover image path
	c.JSON(http.StatusOK, coverImages)
}

func SearchBooks(c *gin.Context) {
	// get query parameters(title, author, category) from the request url
	title := c.Query("title")
	author := c.Query("author")
	category := c.Query("category")

	// query database for books based on the specified criteria
	var books []models.Book
	query := initializers.DB.Model(&books)

	if title != "" {
		query = query.Where("title LIKE ?", title)
	}
	if author != "" {
		query = query.Where("author= ?", author)
	}
	if category != "" {
		query = query.Where("category= ?", category)
	}
	query.Find(&books)
	// return filtered books
	c.JSON(http.StatusOK, books)

}
