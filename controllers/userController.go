package controllers

import (
	"book-store/initializers"
	"book-store/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get the request body"})
	}
	// hash the req password
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to hash password"})
	}
	// create user
	user := models.User{UserName: Body.UserName, Password: string(hash), FirstName: Body.FirstName, LastName: Body.LastName, Role: models.Role(Body.Role)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create user"})
	}
	user.Password = ""
	// return it
	c.JSON(200, gin.H{"User": user})
}

func GetUsers(c *gin.Context) {
	// get the user from the model database
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get the user"})
	}
	// respond with them
	c.IndentedJSON(200, gin.H{"Users": users})
}

func GetUserById(c *gin.Context) {
	// get id of url
	user_id := c.Param("user_id")

	// get the user by the primary key
	var user models.User
	result := initializers.DB.First(&user, user_id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to the user by the given primary key"})
	}
	// respond with them
	c.IndentedJSON(200, gin.H{"Users": user})
}

func UpdateUser(c *gin.Context) {
	// get the id of the url
	user_id := c.Param("user_id")

	//get the data of the req body
	c.Bind(&Body)

	//fined the user where updating
	var user models.User
	result := initializers.DB.First(&user, user_id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to get user"})
	}

	// update it
	initializers.DB.Model(&user).Updates(models.User{UserName: Body.UserName, Password: Body.Password, FirstName: Body.FirstName, LastName: Body.LastName})

	// respond it
	c.JSON(200, gin.H{"User": user})
}

func DeleteUser(c *gin.Context) {
	// get the url of the body
	user_id := c.Param("user_id")

	// delete it
	initializers.DB.Delete(&models.User{}, user_id)

	// respond it
	c.Status(200)
}

func Login(c *gin.Context) {
	var body struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read body"})
		return
	}

	// lock up requested user
	var user models.User
	initializers.DB.First(&user, "user_name = ?", body.UserName)
	fmt.Println(user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid username"})
		return
	}

	// compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid password"})
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
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create tocken"})
		return
	}

	// send it back
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am logged in"})
}

//  Role: models.Role(Body.Role),=========== for line 35
