package models

import (
	"gorm.io/gorm"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	gorm.Model
	UserID    uint   `gorm:"primary key"`
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
	Role            Role   `gorm:"not null"`
	ImageName       string `gorm:"image_name"`
}

type Author struct {
	gorm.Model
	AuthorID    uint   `gorm:"primary key"`
	Name        string `gorm:"not null"`
	Biography   string `gorm:"not null"`
	Nationality string `gorm:"not null"`
	Role        Role   `gorm:"not null"`
}

type Category struct {
	gorm.Model
	CategoryID uint   `gorm:"primary key"`
	Name       string `gorm:"not null"`
	Role       Role   `gorm:"not null"`
}

type Image struct {
	gorm.Model
	ImageID   uint   `gorm:"primary key"`
	ImageName string `gorm:"not null"`
	CoverPath string `gorm:"not null"`
}

type Response struct {
	Meta *Metadata   `json:"meta_data,omitempty"`
	Data interface{} `json:"data"`
	Ok   bool        `json:"ok"`
	Err  *Error      `json:"error"`
}

type Metadata struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}

type Error struct {
	Message string `json:"message"`
	Stack   error  `json:"stack"`
}
