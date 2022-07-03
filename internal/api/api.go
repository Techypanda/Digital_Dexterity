package api

import (
	"fmt"
	"log"
	"net/http"
	"techytechster/digitaldexterity/internal/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

type (
	SimpleValidator struct {
		validator *validator.Validate
	}
	APIConfig struct {
		Port             string
		Database         *database.Database
		JWTSecret        []byte
		JWTRefreshSecret []byte
		SecretKey        string
		CORSList         []string
	}
)

// @title Digital Dexterity API
// @version 1.0
// @description This is an api for measuring users digital dexterity
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://techytechster.com
// @contact.email jonathan_wright@hotmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host noidea.com
// @BasePath /

func (cv *SimpleValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func NewAPI(config APIConfig) {
	e := echo.New()
	e.Validator = &SimpleValidator{validator: validator.New()}
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))) // 20 requests a second
	e.Use(middleware.Recover())                                             // Recover From Panic
	e.Use(middleware.Logger())                                              // Log basic HTTP traffic
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{                  // CORS
		AllowOrigins:     config.CORSList,
		AllowMethods:     []string{echo.POST, echo.GET},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// e.Use(middleware.CSRF()) // TODO: Use CSRF middleware
	e.GET("/", func(c echo.Context) error {
		log.Println("❤️❤️❤️❤️ THUMP ❤️❤️❤️❤️ (HeartBeat Request)")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"date": time.Now().Unix(),
		})
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	User(e, config.Database, config.JWTSecret, config.JWTRefreshSecret)
	refresh := e.Group("/refresh")
	refresh.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    config.JWTRefreshSecret,
		Claims:        &UserTokenClaims{},
		SigningMethod: "HS512",
	}))
	Refresh(refresh, config.Database, config.JWTSecret, config.JWTRefreshSecret)
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
