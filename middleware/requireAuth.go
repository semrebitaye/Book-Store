package middleware

import (
	"book-store/initializers"
	"book-store/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication() gin.HandlerFunc {
	//get the cookie of the req body
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if !(strings.HasPrefix(tokenString, "Bearer ")) {
			fmt.Println("bearer not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(tokenString, "Bearer ")

		//Decode/validateit
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			fmt.Print("Error validation", err)

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Print("Error expired", err)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			// fined the user with token sub
			var user models.User
			initializers.DB.First(&user, claims["sub"])

			if user.ID == 0 {

				err := errors.New("user not found")
				fmt.Print("Error not found", err)
				c.AbortWithError(http.StatusNotFound, err)
				return
			}

			// attach the req
			c.Set("user_id", user.ID)
			c.Set("role", user.Role)
			fmt.Println("user data: ", user.Role)
			//continue

			c.Next()
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.Value("role").(models.Role)

		if userRole == "" {
			err := errors.New("role not found")
			fmt.Println("role not found")
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		method := c.Request.Method
		if userRole == (models.UserRole) {
			if method != http.MethodGet {
				err := errors.New("unauthorized user")
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

		}
		c.Next()
	}

}
