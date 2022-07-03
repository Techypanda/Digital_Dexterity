package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"techytechster/digitaldexterity/internal/database"
	"techytechster/digitaldexterity/internal/encryption"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type RegisterPayload struct {
	Username string `json:"username" form:"username" query:"username" validate:"min=1"`
	Password string `json:"password" form:"username" query:"username" validate:"min=8"`
}
type LoginPayload struct {
	Username string `json:"username" form:"username" query:"username" validate:"required"`
	Password string `json:"password" form:"username" query:"username" validate:"required"`
}
type UserTokenClaims struct {
	UserId   uint   `json:"userID"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func register(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(RegisterPayload)
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to create, %s\n", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   err.Error(),
			})
		} else if err = c.Validate(payload); err != nil {
			log.Printf("Failed to validate payload, %s\n", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   fmt.Sprintf("Expected Username (minimum 1 character) and Password (minimum 8 characters): %s\n", err.Error()),
			})
		}
		encryption := encryption.Encrypt(payload.Password)
		err := db.CreateUser(payload.Username, encryption)
		if err != nil {
			log.Printf("Failed to create user, %s", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   fmt.Sprintf("Failed to register: %s\n", err.Error()),
			})
		}
		log.Printf("Created a new user, %s", payload.Username)
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"created": true,
		})
	}
}

func login(db *database.Database, jwtSecret []byte, refreshSecret []byte) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(LoginPayload)
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to bind payload to login request, %s\n", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("Failed to validate payload for login request, %s\n", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		u := db.GetUser(payload.Username)
		if u == nil {
			log.Printf("Failed to get a user for requested username\n")
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "incorrect username/password",
			})
		}
		if encryptedRes := encryption.Encrypt(payload.Password); bytes.Compare(encryptedRes, u.EncryptedPassword) != 0 {
			log.Printf("Failed to authenticate, password incorrect\n")
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "incorrect username/password",
			})
		}
		claims := &UserTokenClaims{
			u.ID,
			u.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			},
		}
		refreshClaims := &UserTokenClaims{
			u.ID,
			u.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		rToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
		t, err := token.SignedString(jwtSecret)
		if err != nil {
			log.Printf("failed to generate a token: %s\n", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		rT, err := rToken.SignedString(refreshSecret)
		if err != nil {
			log.Printf("failed to generate a token: %s\n", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		log.Printf("Successfully logged in: %s\n", u.Username)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"username":     payload.Username,
			"loggedIn":     true,
			"token":        t,
			"refreshToken": rT,
		})
	}
}

func User(e *echo.Echo, db *database.Database, jwtSecret []byte, refreshSecret []byte) {
	e.POST("/register", register(db))
	e.POST("/login", login(db, jwtSecret, refreshSecret))
}

func Refresh(e *echo.Group, db *database.Database, jwtSecret []byte, refreshSecret []byte) {
	e.GET("", func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*UserTokenClaims)
		user := db.GetUser(claims.Username)
		claims = &UserTokenClaims{
			user.ID,
			user.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			},
		}
		refreshClaims := &UserTokenClaims{
			user.ID,
			user.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			},
		}
		token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		rToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
		t, err := token.SignedString(jwtSecret)
		if err != nil {
			log.Printf("failed to generate a token: %s\n", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		rT, err := rToken.SignedString(refreshSecret)
		if err != nil {
			log.Printf("failed to generate a token: %s\n", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		log.Printf("Successfully refreshed tokens for: %s\n", user.Username)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"token":        t,
			"refreshToken": rT,
		})
	})
}
