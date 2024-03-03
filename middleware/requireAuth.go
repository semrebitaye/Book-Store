package middleware

import (
	"book-store/controllers"
	"book-store/initializers"
	"book-store/models"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication() gin.HandlerFunc {
	//get the Bearer of the req body
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if !(strings.HasPrefix(tokenString, "Bearer ")) {
			controllers.ErrorResponse(c, http.StatusUnauthorized, &models.Error{Message: "Bearer not found"})
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
			controllers.ErrorResponse(c, http.StatusUnauthorized, &models.Error{Message: "Validation error", Stack: err})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				controllers.ErrorResponse(c, http.StatusUnauthorized, &models.Error{Message: "Tocken Expired", Stack: err})
				return
			}
			// fined the user with token sub
			var user models.User
			initializers.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				controllers.ErrorResponse(c, http.StatusNotFound, &models.Error{Message: "User not found", Stack: err})
				return
			}

			// attach the req
			c.Set("user_id", user.ID)
			c.Set("role", user.Role)
			fmt.Println("user data: ", user.Role)
			//continue

			c.Next()
		} else {
			controllers.ErrorResponse(c, http.StatusInternalServerError, &models.Error{Message: "Tocken not found", Stack: err})
			return
		}
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.Value("role").(models.Role)

		if userRole == "" {
			controllers.ErrorResponse(c, http.StatusUnauthorized, &models.Error{Message: "role not found"})
			return
		}

		method := c.Request.Method
		if userRole == (models.UserRole) {
			if method != http.MethodGet {
				controllers.ErrorResponse(c, http.StatusBadRequest, &models.Error{Message: "Unauthorized user"})
				return
			}

		}
		c.Next()
	}

}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Assign the new context with timeout to the request
		c.Request = c.Request.WithContext(ctx)

		// Call the next handler
		c.Next()

		// If the request context has been canceled, respond with a timeout error
		if ctx.Err() == context.DeadlineExceeded {
			controllers.ErrorResponse(c, http.StatusRequestTimeout, &models.Error{Message: "Request timeout"})
			return
		}
	}
}
