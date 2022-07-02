package api

import (
	"bytes"
	"fmt"
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
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   err.Error(),
			})
		} else if err = c.Validate(payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   fmt.Sprintf("Expected Username (minimum 1 character) and Password (minimum 8 characters): %s", err.Error()),
			})
		}
		encryption := encryption.Encrypt(payload.Password)
		err := db.CreateUser(payload.Username, encryption)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   fmt.Sprintf("Failed to register: %s", err.Error()),
			})
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"created": true,
		})
	}
}

func login(db *database.Database, jwtSecret []byte) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(LoginPayload)
		if err := c.Bind(payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		} else if err = c.Validate(payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		u := db.GetUser(payload.Username)
		if u == nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "incorrect username/password",
			})
		}
		if encryptedRes := encryption.Encrypt(payload.Password); bytes.Compare(encryptedRes, u.EncryptedPassword) != 0 {
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
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		t, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Username": payload.Username,
			"LoggedIn": true,
			"Token":    t,
		})
	}
}

func User(e *echo.Echo, db *database.Database, jwtSecret []byte) {
	e.POST("/api/v1/register", register(db))
	e.POST("/api/v1/login", login(db, jwtSecret))
}
