package models

import "gorm.io/gorm"

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	gorm.Model
	UserName  string `gorm:"not null"`
	Password  string `json:"password,omitempty" gorm:"not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Books []Book `gorm:"many2many:user_books"`

	Role Role `gorm:"not null"`
}

type Book struct {
	gorm.Model
	Title           string  `gorm:"not null"`
	PublicationDate string  `gorm:"not null"`
	Price           float64 `gorm:"not null"`
	Quantity        uint    `gorm:"not null"`
	UserID          uint    `gorm:"references:id;not null"`
	User            *User
	AuthorID        uint `gorm:"references:id;not null"`
	Author          *Author
	CategoryID      uint `gorm:"references:id;not null"`
	Category        *Category
	CoverPath       string `gorm:"not null"`
}

type Author struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Biography   string `gorm:"not null"`
	Nationality string `gorm:"not null"`
}

type Category struct {
	gorm.Model
	Name string `gorm:"not null"`
}
