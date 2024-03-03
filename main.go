package main

import (
	"book-store/controllers"
	"book-store/initializers"
	"book-store/middleware"
	"book-store/migrate"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("connecting to db.........")

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	MigrateUp()
}

func MigrateUp() {
	migrate.MigrateUpModels()
}
func MigrateDown() {
	migrate.MigrateDownModels()
}
func main() {
	r := gin.Default()
	r.POST("/create", controllers.CreateUser)
	r.POST("/login", controllers.Login)

	r.Use(middleware.Authentication(), middleware.Authorization(), middleware.TimeoutMiddleware(5*time.Second))

	r.GET("/get", controllers.GetUsers)
	r.GET("/get/:user_id", controllers.GetUserById)
	r.PUT("/update/:user_id", controllers.UpdateUser)
	r.DELETE("/delete/:user_id", controllers.DeleteUser)

	r.POST("/createBook", controllers.CreateBook)
	r.GET("/getBook", controllers.GetBooks)
	r.GET("/getBook/:book_id", controllers.GetBookByID)
	r.PUT("/updateBook/:book_id", controllers.UpdateBook)
	r.DELETE("/deleteBook/:book_id", controllers.DeleteBook)

	r.POST("/createAuth", controllers.CreateAuthor)
	r.GET("/getAuth", controllers.GetAuthors)
	r.GET("/getAuth/:author_id", controllers.GetBookByID)
	r.PUT("/updateAuth/:author_id", controllers.UpdateBook)
	r.DELETE("/deleteAuth/:author_id", controllers.DeleteBook)

	r.POST("/createCategory", controllers.CreateCategory)
	r.GET("/getCategory", controllers.GetCategories)
	r.GET("/getCategory/:category_id", controllers.GetCategoryById)
	r.PUT("/updateCategory/:category_id", controllers.UpdateCategory)
	r.DELETE("/deleteCategory/:category_id", controllers.DeleteCategory)

	r.GET("/validate", controllers.Validate)

	r.POST("/upload/book-cover", controllers.UploadBookCover)
	r.GET("/get-book-cover/:book_id/cover", controllers.GetBookCoverImage)

	r.Run()
}
