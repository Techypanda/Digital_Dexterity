package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username            string               `gorm:"unique;column:username" json:"username"`
	EncryptedPassword   []byte               `gorm:"column:encrypted_password" json:"encrypted_password"`
	SelfAssessment      *SelfAssessment      // 1 -> 1 Relationship
	ExternalAssessments []ExternalAssessment // 1 -> Many
	GithubID            uint64               `gorm:"unique"` // NULLABLE
}

func NewUser(username string, password []byte) *User {
	return &User{
		Username:          username,
		EncryptedPassword: password,
	}
}

func (db *Database) GetGithubUser(githubID uint64) *User {
	user := new(User)

	db.db.Model(&User{}).Where(&User{GithubID: githubID}).First(user)

	if user.Username == "" {
		return nil
	}

	return user
}

func (db *Database) NewGithubUser(username string, githubID uint64) error {
	user := new(User)

	db.db.Model(&User{}).Where(&User{Username: username}).First(user)

	if user.Username != "" {
		username = fmt.Sprintf("%s-%s", username, uuid.NewString())
	}

	user = &User{
		Username:          username,
		EncryptedPassword: []byte{},
		GithubID:          githubID,
	}
	result := db.db.Create(user)

	return result.Error
}

var ErrUsernameTaken = errors.New("username is already taken")

func (db *Database) CreateUser(username string, password []byte) error {
	user := new(User)

	db.db.Model(&User{}).Where(&User{Username: username}).First(user)

	if user.Username != "" {
		return ErrUsernameTaken
	}

	user = NewUser(username, password)
	result := db.db.Create(user)

	return result.Error
}

func (db *Database) GetUser(username string) *User {
	user := new(User)
	db.db.Model(&User{}).Where(&User{Username: username}).First(user)

	if user.Username == "" {
		return nil
	}

	return user
}
