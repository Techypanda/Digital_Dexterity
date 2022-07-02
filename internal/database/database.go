package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}
type DatabaseConfig struct {
	Username string
	Password string
	IP       string
	TLS      string
}

func NewDatabase(config DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/digitaldexterity?tls=%s&parseTime=true", config.Username, config.Password, config.IP, config.TLS)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return &Database{db: db}, err
}
func (db *Database) Migrate() error {
	return db.db.AutoMigrate(&User{})
}
