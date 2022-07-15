package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}
type Config struct {
	Username string
	Password string
	IP       string
	TLS      string
}

func NewDatabase(config Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/digitaldexterity?tls=%s&parseTime=true",
		config.Username,
		config.Password,
		config.IP,
		config.TLS,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	return &Database{db: db}, fmt.Errorf("failed to initialize database: %w", err)
}
func (db *Database) Migrate() error {
	if err := db.db.AutoMigrate(&User{}, &SelfAssessment{}, &ExternalAssessment{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}
