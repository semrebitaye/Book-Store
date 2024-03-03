package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"book-store/utilities"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var authorRequest struct {
	Name        string `json:"name"`
	Biography   string `json:"biography"`
	Nationality string `json:"nationality"`
	Role        string `json:"role"`
	AuthorID    uint   `json:"author_id"`
}

func CreateAuthor(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&authorRequest)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}
	// create a post
	author := models.Author{
		Model:       gorm.Model{},
		AuthorID:    authorRequest.AuthorID,
		Name:        authorRequest.Name,
		Biography:   authorRequest.Biography,
		Nationality: authorRequest.Nationality,
		Role:        models.Role(categoryRequest.Role),
	}

	result := initializers.DB.Create(&author)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to create the author", Stack: result.Error})
		return
	}
	// return it
	SuccessResponse(c, http.StatusOK, author, &models.Metadata{})
}

func GetAuthors(c *gin.Context) {
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
	var authors []models.Author
	db := initializers.DB
	if pgParam.Search != "" {
		db.Where("name LIKE %%?%% OR biography LIKE %%?%% OR nationality LIKE %%?%%", pgParam.Search, pgParam.Search, pgParam.Search)
	} else {
		for _, filter := range filterParam.Filters {
			db = db.Where(fmt.Sprintf("%s %s %v", filter.ColumnName, filter.Operator, filter.Value))
		}
	}

	offset := (filterParam.Page - 1) * filterParam.PerPage
	result := db.Offset(offset).Limit(filterParam.PerPage).Order(filterParam.Sort.ColumnName + " " + filterParam.Sort.Value).Find(&authors)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: result.Error})
		return
	}

	SuccessResponse(c, http.StatusOK, authors, &models.Metadata{
		TotalCount: len(authors),
		Page:       filterParam.Page,
		PerPage:    filterParam.PerPage,
	})
}

func GetAuthorById(c *gin.Context) {
	// get id of url
	author_id := c.Param("author_id")

	// get the user by the primary key
	var author models.Author
	err := initializers.DB.First(&author, author_id)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the user by the given primary key", Stack: err.Error})
		return
	}
	// respond with them
	SuccessResponse(c, http.StatusOK, author, &models.Metadata{})
}

func UpdateAuthor(c *gin.Context) {
	// get the id of the url
	author_id := c.Param("author_id")

	//get the data of the req body
	err := c.Bind(&authorRequest)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}
	//fined the user where updating
	var author models.Author
	if initializers.DB.First(&author, author_id) != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the author", Stack: err})
		return
	}

	// update it
	if err := initializers.DB.Model(&author).Updates(models.Author{
		Name:        authorRequest.Name,
		Biography:   authorRequest.Biography,
		Role:        models.Role(categoryRequest.Role),
		AuthorID:    authorRequest.AuthorID,
		Nationality: authorRequest.Nationality,
	}); err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to update the author", Stack: err.Error})
		return
	}

	// respond it
	SuccessResponse(c, http.StatusOK, author, &models.Metadata{})
}

func DeleteAuthor(c *gin.Context) {
	// get the url of the body
	author_id := c.Param("author_id")

	// delete it
	err := initializers.DB.Delete(&models.Author{}, author_id)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to delete the author", Stack: err.Error})
		return
	}

	// respond it
	SuccessResponse(c, http.StatusOK, "", &models.Metadata{})
}
