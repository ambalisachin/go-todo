package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	
)

type todo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type user struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// Compares the provided password with the hashed password stored in the user struct.
// It uses the bcrypt library to compare the two passwords, returning an error if they do not match.
// If the passwords match, it returns nil.
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
