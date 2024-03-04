package controllers

import (
	"book-store/initializers"
	"book-store/models"
	"book-store/utilities"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var Body struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

func CreateUser(c *gin.Context) {
	// get data of the req body
	err := c.Bind(&Body)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}
	// hash the req password
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to hash the password", Stack: err})
		return
	}
	// create user
	user := models.User{
		UserName:  Body.UserName,
		Password:  string(hash),
		FirstName: Body.FirstName,
		LastName:  Body.LastName,
		Role:      models.Role(Body.Role),
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to create the user", Stack: result.Error})
		return
	}
	user.Password = ""
	// return it
	SuccessResponse(c, http.StatusOK, user, &models.Metadata{})
}

func GetUsers(c *gin.Context) {
	var pgParam utilities.PaginationParam

	if err := c.BindQuery(&pgParam); err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the query", Stack: err})
		return
	}

	filterParam, err := utilities.ExtractPagination(pgParam)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Retrieve users with pagination and sorting
	var users []models.User
	db := initializers.DB

	// search and filter
	if pgParam.Search != "" {
		db.Where("first_name LIKE %%?%% OR last_name LIKE %%?%% OR user_name LIKE %%?%%", pgParam.Search, pgParam.Search, pgParam.Search)
	} else if filterParam.Filters != nil {
		for _, filter := range filterParam.Filters {
			db = db.Where(fmt.Sprintf("%s %s %v", filter.ColumnName, filter.Operator, filter.Value))
		}
	}

	offset := (filterParam.Page - 1) * filterParam.PerPage

	result := db.Offset(offset).Limit(filterParam.PerPage).Order(filterParam.Sort.ColumnName + " " + filterParam.Sort.Value).Find(&users)

	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: result.Error})
		return
	}

	SuccessResponse(c, http.StatusOK, users, &models.Metadata{
		TotalCount: len(users),
		Page:       filterParam.Page,
		PerPage:    filterParam.PerPage,
	})
}

func GetUserById(c *gin.Context) {
	// get id of url
	user_id := c.Param("user_id")

	// get the user by the primary key
	var user models.User
	result := initializers.DB.First(&user, user_id)
	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the user by the given primary key", Stack: result.Error})
		return
	}
	// respond with them
	SuccessResponse(c, http.StatusOK, user, &models.Metadata{})
}

func UpdateUser(c *gin.Context) {
	// get the id of the url
	user_id := c.Param("user_id")

	//get the data of the req body
	err := c.Bind(&Body)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to bind the request body", Stack: err})
		return
	}

	//fined the user where updating
	var user models.User
	result := initializers.DB.First(&user, user_id)
	if result.Error != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to get the user", Stack: result.Error})
		return
	}

	// update it
	initializers.DB.Model(&user).Updates(models.User{
		UserName:  Body.UserName,
		Password:  Body.Password,
		FirstName: Body.FirstName,
		LastName:  Body.LastName,
	})

	// respond it
	SuccessResponse(c, http.StatusOK, user, &models.Metadata{})
}

func DeleteUser(c *gin.Context) {
	// get the url of the body
	user_id := c.Param("user_id")

	// delete it
	err := initializers.DB.Delete(&models.User{}, user_id)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to delete the user", Stack: err.Error})
		return
	}
	// respond it
	SuccessResponse(c, http.StatusOK, "", &models.Metadata{})
}

func Login(c *gin.Context) {
	var body struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to read the body"})
		return
	}

	// lock up requested user
	var user models.User
	initializers.DB.First(&user, "user_name = ?", body.UserName)
	fmt.Println(user)

	if user.ID == 0 {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Invalid user name"})
		return
	}

	// compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to compare the password", Stack: err})
		return
	}

	//generate a jwt tocken
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Failed to create tocken", Stack: err})
		return
	}

	// send it back
	SuccessResponse(c, http.StatusOK, tokenString, &models.Metadata{})
}

func Validate(c *gin.Context) {
	SuccessResponse(c, http.StatusOK, "I am logged in", &models.Metadata{})
}
