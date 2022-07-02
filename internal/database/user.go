package database

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username            string               `gorm:"unique;column:username" json:"username"`
	EncryptedPassword   []byte               `gorm:"column:encrypted_password" json:"encrypted_password"`
	SelfAssessment      *SelfAssessment      // 1 -> 1 Relationship
	ExternalAssessments []ExternalAssessment // 1 -> Many
}

func NewUser(username string, password []byte) *User {
	return &User{
		Username:          username,
		EncryptedPassword: password,
	}
}

func (db *Database) CreateUser(username string, password []byte) error {
	u := new(User)
	db.db.Model(&User{}).Where(&User{Username: username}).First(u)
	if u.Username != "" {
		return errors.New("username is already taken")
	}
	u = NewUser(username, password)
	result := db.db.Create(u)
	return result.Error
}

func (db *Database) GetUser(username string) *User {
	u := new(User)
	db.db.Model(&User{}).Where(&User{Username: username}).First(u)
	if u.Username == "" {
		return nil
	}
	return u
}
