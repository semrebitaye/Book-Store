package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"book-store/utilities"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var categoryRequest struct {
	Name       string `json:"name"`
	Role       string `json:"role"`
	CategoryID uint   `json:"category_id"`
}

func CreateCategory(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&categoryRequest)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}
	// create a post
	category := models.Category{
		Name:       categoryRequest.Name,
		Role:       models.Role(categoryRequest.Role),
		CategoryID: bookRequest.CategoryID,
	}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to create the category", Stack: result.Error})
		return
	}

	// return it
	SuccessResponse(c, http.StatusOK, category, &models.Metadata{})
}

func GetCategories(c *gin.Context) {
	var pgParam utilities.PaginationParam

	if err := c.BindQuery(&pgParam); err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the query parameters", Stack: err})
		return
	}

	fmt.Println("params:", pgParam)
	filterParam, err := utilities.ExtractPagination(pgParam)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, &models.Error{Message: "Failed to extract pagination", Stack: err})
		return
	}
	// Retrieve books with pagination and sorting
	var categories []models.Category
	db := initializers.DB
	if pgParam.Search != "" {
		db.Where("name LIKE %%?%%", pgParam.Search)
	} else {
		for _, filter := range filterParam.Filters {
			db = db.Where(fmt.Sprintf("%s %s %v", filter.ColumnName, filter.Operator, filter.Value))
		}
	}

	offset := (filterParam.Page - 1) * filterParam.PerPage
	result := db.Offset(offset).Limit(filterParam.PerPage).Order(filterParam.Sort.ColumnName + " " + filterParam.Sort.Value).Find(&categories)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: result.Error})
		return
	}

	SuccessResponse(c, http.StatusOK, categories, &models.Metadata{
		TotalCount: len(categories),
		Page:       filterParam.Page,
		PerPage:    filterParam.PerPage,
	})
}

func GetCategoryById(c *gin.Context) {
	// get id of url
	category_id := c.Param("category_id")

	// get the user
	var category models.Category
	result := initializers.DB.First(&category, category_id)
	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the data", Stack: result.Error})
		return
	}
	// respond with them
	SuccessResponse(c, http.StatusOK, category, &models.Metadata{})
}

func UpdateCategory(c *gin.Context) {
	// get the id of the url
	category_id := c.Param("category_id")

	//get the data of the req body
	err := c.Bind(&categoryRequest)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}

	//fined the user where updating
	var category models.User
	result := initializers.DB.First(&category, category_id)
	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to find the updating category", Stack: result.Error})
		return
	}

	// update it
	initializers.DB.Model(&category).Updates(models.Category{
		Name: categoryRequest.Name, 
		Role: category.Role, 
		CategoryID: bookRequest.CategoryID,
	})

	// respond it
	SuccessResponse(c, http.StatusOK, category, &models.Metadata{})
}

func DeleteCategory(c *gin.Context) {
	// get the url of the body
	category_id := c.Param("category_id")

	// delete it
	err := initializers.DB.Delete(&models.Category{}, category_id)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to delete the category", Stack: err.Error})
		return
	}

	// respond it
	SuccessResponse(c, http.StatusOK, "", &models.Metadata{})
}
