package user

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	ID      uint
	Name    string
	Phone   string
	Address string
}

//export
