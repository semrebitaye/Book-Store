package migrate

import (
	"book-store/initializers"
	"book-store/models"
	"fmt"
)

func MigrateUpModels() {
	fmt.Println("migrating up.............")
	initializers.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Author{}, &models.Category{})
}

func MigrateDownModels() {
	fmt.Println("migrating down.............")
	initializers.DB.Migrator().DropTable(&models.User{})
	initializers.DB.Migrator().DropTable(&models.Book{})
	initializers.DB.Migrator().DropTable(&models.Author{})
	initializers.DB.Migrator().DropTable(&models.Category{})
}
