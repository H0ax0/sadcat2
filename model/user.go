package model

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

type User struct {
	ID          int    `gorm:"unique;primaryKey;not null"`
	Name        string `gorm:"not null"`
	AccessLevel int    `gorm:"not null"`
}

var AutorizedUsers []User

func Preload() {
	err := DB.Where("access_level > ?", 0).Find(&AutorizedUsers).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Loaded %d autorized users", len(AutorizedUsers))
}

func IdExist(id int) bool {
	if err := DB.Where(&User{ID: id}).First(&User{}).Error; err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func Elevate(id int) {
	u := &User{ID: id}
	err := DB.First(u).Error
	if err != nil {
		fmt.Println("Error elevating")
	}
	u.AccessLevel = 1
	err = DB.Save(u).Error
	if err != nil {
		fmt.Println("Error elevating 2")
	}
	AutorizedUsers = append(AutorizedUsers, *u)
}

func IsAuth(id int) bool {

	for _, usr := range AutorizedUsers {
		if usr.ID == id {
			return true
		}
	}

	return false
}

func NewClient(user *tb.User) error {
	u := &User{
		ID:          user.ID,
		Name:        user.Username,
		AccessLevel: 0,
	}

	if err := DB.Create(u).Error; err != nil {
		return err
	}

	return nil
}
