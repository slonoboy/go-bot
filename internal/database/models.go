package database

import (
	"gorm.io/gorm"
)

type Conversion struct {
	gorm.Model
	UserID   uint
	ChatID   int64
	Crypto   string
	Currency string
}

type ConversionCount struct {
	Crypto   string
	Currency string
	Count    int
}

type User struct {
	gorm.Model
	TGID        int64 `gorm:"uniqueIndex"`
	FirstName   string
	LastName    string
	Username    string
	Conversions []Conversion
}
