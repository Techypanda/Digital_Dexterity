package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"
	"techytechster/digitaldexterity/internal/api"
	"techytechster/digitaldexterity/internal/database"
)

const DefaultPort = "8080"
const DefaultTLS = "true"

const SecretLength = 1248

func main() {
	port, exists := os.LookupEnv("port")
	if !exists {
		port = DefaultPort
	}

	secretKey, exists := os.LookupEnv("secret_key")
	if !exists {
		panic("secret_key is not defined")
	}

	corsList, exists := os.LookupEnv("cors_list")
	if !exists {
		panic("cors_list is not defined")
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

	db, err := database.NewDatabase(database.Config{
		Username: dbUsername,
		Password: dbPassword,
		IP:       dbAddress,
		TLS:      DefaultTLS,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to setup db: %s", err.Error()))
	}

	log.Println("Migrating Database")

	if err = db.Migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate db: %s", err.Error()))
	}

	b := make([]byte, SecretLength)
	if _, err = rand.Read(b); err != nil {
		panic(err.Error())
	}

	jwtSecret := []byte(fmt.Sprintf("%x", b)[:SecretLength])
	b = make([]byte, SecretLength)

	if _, err = rand.Read(b); err != nil {
		panic(err.Error())
	}

	refreshJwtSecret := []byte(fmt.Sprintf("%x", b)[:1248])

	api.NewAPI(api.Config{
		Port:             port,
		Database:         db,
		JWTSecret:        jwtSecret,
		JWTRefreshSecret: refreshJwtSecret,
		SecretKey:        secretKey,
		CORSList:         strings.Split(corsList, ","),
	})
}
