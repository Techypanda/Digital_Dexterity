package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"techytechster/digitaldexterity/internal/api"
	"techytechster/digitaldexterity/internal/database"
	"time"
)

const DEFAULT_PORT = "8080"
const DEFAULT_TLS = "true"

func main() {
	rand.Seed(time.Now().UnixNano())
	port, exists := os.LookupEnv("port")
	if !exists {
		port = DEFAULT_PORT
	}
	secretKey, exists := os.LookupEnv("secret_key")
	if !exists {
		panic("secret_key is not defined")
	}
	dbUsername, exists := os.LookupEnv("db_username")
	if !exists {
		panic("db_username is not defined")
	}
	dbPassword, exists := os.LookupEnv("db_password")
	if !exists {
		panic("db_password is not defined")
	}
	dbAddress, exists := os.LookupEnv("db_address")
	if !exists {
		panic("db_address is not defined")
	}
	log.Println("Initializing Database")
	db, err := database.NewDatabase(database.DatabaseConfig{
		Username: dbUsername,
		Password: dbPassword,
		IP:       dbAddress,
		TLS:      DEFAULT_TLS,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to setup db: %s", err.Error()))
	}
	log.Println("Migrating Database")
	if err = db.Migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate db: %s", err.Error()))
	}
	b := make([]byte, 1248)
	rand.Read(b)
	jwtSecret := []byte(fmt.Sprintf("%x", b)[:1248])
	b = make([]byte, 1248)
	rand.Read(b)
	refreshJwtSecret := []byte(fmt.Sprintf("%x", b)[:1248])
	api.NewAPI(api.APIConfig{
		Port:             port,
		Database:         db,
		JWTSecret:        jwtSecret,
		JWTRefreshSecret: refreshJwtSecret,
		SecretKey:        secretKey,
	})
}
