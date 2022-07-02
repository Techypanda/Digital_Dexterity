package api

import (
	"fmt"
	"net/http"
	"techytechster/digitaldexterity/internal/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	SimpleValidator struct {
		validator *validator.Validate
	}
	APIConfig struct {
		Port      string
		Database  *database.Database
		JWTSecret []byte
		SecretKey string
	}
)

func (cv *SimpleValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func NewAPI(config APIConfig) {
	e := echo.New()
	e.Validator = &SimpleValidator{validator: validator.New()}
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"date": time.Now().Unix(),
		})
	})
	User(e, config.Database, config.JWTSecret)
	r := e.Group("/api/v1")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    config.JWTSecret,
		Claims:        &UserTokenClaims{},
		SigningMethod: "HS512",
	}))
	r.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"date": time.Now().Unix(),
			"auth": true,
		})
	})
	Assessment(r, config.Database, config.JWTSecret)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
