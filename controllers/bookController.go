package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"book-store/utilities"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var bookRequest struct {
	Title           string  `json:"title"`
	PublicationDate string  `json:"publication_date"`
	Price           float64 `json:"price"`
	Quantity        uint    `json:"quantity"`
	UserID          uint    `json:"user_id"`
	AuthorID        uint    `json:"author_id"`
	CategoryID      uint    `json:"category_id"`
	Role            string  `json:"role"`
	ImageName       string  `json:"image_name"`
}

func CreateBook(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&bookRequest)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body"})
		return
	}
	// create a post
	u := c.GetUint("user_id")

	if u == 0 {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the user id"})
		return
	}

	book := models.Book{
		UserID:          uint(u),
		Title:           bookRequest.Title,
		PublicationDate: bookRequest.PublicationDate,
		Price:           bookRequest.Price,
		Quantity:        bookRequest.Quantity,
		AuthorID:        bookRequest.AuthorID,
		CategoryID:      bookRequest.CategoryID,
		Role:            models.Role(bookRequest.Role),
		ImageName:       bookRequest.ImageName,
	}

	result := initializers.DB.Create(&book)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to create the book", Stack: result.Error})
		return
	}

	// return it
	SuccessResponse(c, http.StatusOK, book, &models.Metadata{})
}

func GetBooks(c *gin.Context) {
	var pgParam utilities.PaginationParam

	if err := c.BindQuery(&pgParam); err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to to bind the query", Stack: err})

		return
	}
	fmt.Println("params:", pgParam)
	filterParam, err := utilities.ExtractPagination(pgParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Retrieve books with pagination and sorting
	var books []models.Book
	db := initializers.DB
	if pgParam.Search != "" {
		db.Where("title LIKE %%?%% OR author LIKE %%?%% OR category LIKE %%?%%", pgParam.Search, pgParam.Search, pgParam.Search)
	} else {
		for _, filter := range filterParam.Filters {
			db = db.Where(fmt.Sprintf("%s %s %v", filter.ColumnName, filter.Operator, filter.Value))
		}
	}

	offset := (filterParam.Page - 1) * filterParam.PerPage
	result := db.Offset(offset).Limit(filterParam.PerPage).Order(filterParam.Sort.ColumnName + " " + filterParam.Sort.Value).Find(&books)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: result.Error})
		return
	}

	SuccessResponse(c, http.StatusOK, books, &models.Metadata{
		TotalCount: len(books),
		Page:       filterParam.Page,
		PerPage:    filterParam.PerPage,
	})
}

func GetBookByID(c *gin.Context) {
	// get id of url
	book_id := c.Param("book_id")

	//get the book
	var book models.Book
	result := initializers.DB.First(&book, book_id)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the book", Stack: result.Error})
	}

	// respond
	SuccessResponse(c, http.StatusOK, book, &models.Metadata{})
}

func UpdateBook(c *gin.Context) {
	// get the id of the url
	book_id := c.Param("book_id")

	//get the data of the req body
	err := c.Bind(&bookRequest)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}

	//fined the user where updating
	var book models.User
	result := initializers.DB.First(&book, book_id)
	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the book", Stack: result.Error})
		return
	}

	// update it
	initializers.DB.Model(&book).Updates(models.Book{
		Title:           bookRequest.Title,
		Role:            models.Role(bookRequest.Role),
		PublicationDate: bookRequest.PublicationDate,
		Price:           bookRequest.Price,
		Quantity:        bookRequest.Quantity})

	// respond it
	SuccessResponse(c, http.StatusOK, book, &models.Metadata{})
}

func DeleteBook(c *gin.Context) {
	// get the url of the body
	book_id := c.Param("book_id")

	// delete it
	err := initializers.DB.Delete(&models.Book{}, book_id)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to delete the book", Stack: err.Error})
		return
	}

	// respond it
	SuccessResponse(c, http.StatusOK, "", &models.Metadata{})
}

func UploadBookCover(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to parse the image", Stack: err})
		return
	}
	defer file.Close()

	// Define the maximum file size (32 MB)
	const maxFileSize = 32

	// Validate file size
	if header.Size > maxFileSize {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "File size exceeds the maximum limit of 32 MB"})
		return
	}

	// check if the file format is supported
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Unsupported file format"})
		return
	}

	// save the file to the server
	filename := uuid.NewString() + header.Filename
	destination := "./v1/image/" + filename

	if err := c.SaveUploadedFile(header, destination); err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to save file", Stack: err})
		return
	}

	// save the file path to the database
	image := models.Image{
		ImageName: c.PostForm("image_name"),
		CoverPath: destination,
	}
	result := initializers.DB.Create(&image)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to create the image", Stack: result.Error})
		return
	}

	// response it
	SuccessResponse(c, http.StatusOK, image, &models.Metadata{})
}

func GetBookCoverImage(c *gin.Context) {
	// get the book id from request url parameter
	book_id := c.Param("book_id")

	// querry database for books based on ids
	var images []models.Image
	err := initializers.DB.Where("id IN (?)", book_id).Find(&images)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the image", Stack: err.Error})
		return
	}

	// create a map to store book id and associated book cover image path
	coverImages := make(map[uint]string)
	for _, image := range images {
		coverImages[image.ID] = image.CoverPath
	}

	// return cover image path
	SuccessResponse(c, http.StatusOK, coverImages, &models.Metadata{})
}
