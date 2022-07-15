package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
type StatePayload struct {
	State string `json:"state" form:"state" query:"state" validate:"required"`
}
type GithubLoginPayload struct {
	State       string `json:"state" form:"state" query:"state" validate:"required"`
	Code        string `json:"code" form:"code" query:"code" validate:"required"`
	RedirectURI string `json:"redirectURI" form:"redirectURI" query:"redirectURI" validate:"required"`
}
type GithubResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
type GithubMeResponse struct {
	Login string `json:"login"`
	ID    uint64 `json:"id"`
}
type UserTokenClaims struct {
	UserID   uint   `json:"userID"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var AccessTokenExpiryMinutes = 30
var RefreshTokenExpiryHours = 2

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

		if encryptedRes := encryption.Encrypt(payload.Password); !bytes.Equal(encryptedRes, u.EncryptedPassword) {
			log.Printf("Failed to authenticate, password incorrect\n")

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "incorrect username/password",
			})
		}

		claims := &UserTokenClaims{
			u.ID,
			u.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(AccessTokenExpiryMinutes)).Unix(),
			},
		}
		refreshClaims := &UserTokenClaims{
			u.ID,
			u.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(RefreshTokenExpiryHours)).Unix(),
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

func addNewState(stateStore *database.StateStore) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(StatePayload)
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to bind payload to state request, %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("Failed to validate payload for state request, %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		go stateStore.AddNewState(payload.State)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": "state has been queued",
		})
	}
}

func githubLogin(stateStore *database.StateStore, database *database.Database, githubClientID string, githubClientSecret string, jwtSecret []byte, refreshSecret []byte) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(GithubLoginPayload)
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to bind payload to github login request, %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("Failed to validate payload for github login request, %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		// Step 1: Check State

		// Step 2: Post Github
		req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
		if err != nil {
			log.Println("failed to create request to github - ", err.Error())

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "unable to request github login",
			})
		}

		queries := req.URL.Query()
		queries.Add("client_id", githubClientID)
		queries.Add("client_secret", githubClientSecret)
		queries.Add("code", payload.Code)
		queries.Add("redirect_uri", payload.RedirectURI)
		req.URL.RawQuery = queries.Encode()
		req.Header.Add("Accept", "application/json")

		client := http.Client{}
		resp, err := client.Do(req)

		log.Println(resp.StatusCode)

		if resp.StatusCode != http.StatusOK {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "bad code",
			})
		} else if err != nil {
			log.Println("failed to contact github - ", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "failed to contact github",
			})
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("failed to read body of github response - ", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "failed to read body of github response",
			})
		}
		var response GithubResponse
		json.Unmarshal(body, &response)
		githubResponseNotEmpty := response.AccessToken != "" && response.Scope != "" && response.TokenType != ""
		if !githubResponseNotEmpty {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "failed to get response from github, you probably did something wrong",
			})
		}

		// Step 3: Use Access Token To Construct A User
		req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			log.Println("failed to connect to api.github.com/user - ", err.Error())

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "failed to connect to api.github.com",
			})
		}
		req.Header.Add("Authorization", fmt.Sprintf("token %s", response.AccessToken))
		req.Header.Add("Accept", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Println("failed to auth to api.github.com/profile - ", err.Error())

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "failed to auth to api.github.com/profile",
			})
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("failed to read body of github response - ", err.Error())

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "failed to read body of github response",
			})
		}
		var userResponse GithubMeResponse
		json.Unmarshal(body, &userResponse)

		// Step 4: Check If User Exists, if so return the jwts
		user := database.GetGithubUser(userResponse.ID)
		if user != nil {
			tokenPair, err := generateTokenPair(user, jwtSecret, refreshSecret)
			if err != nil {
				log.Println(err.Error())

				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "failed to generate tokens",
				})
			}
			return c.JSON(http.StatusOK, tokenPair)
		}

		// Step 5: Create The User, Then Return JWTS for that user
		err = database.NewGithubUser(userResponse.Login, userResponse.ID)
		if err != nil {
			log.Println(err.Error())

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "failed to create github user",
			})
		}
		user = database.GetGithubUser(userResponse.ID)
		if user != nil {
			tokenPair, err := generateTokenPair(user, jwtSecret, refreshSecret)
			if err != nil {
				log.Println(err.Error())

				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "failed to generate tokens",
				})
			}
			return c.JSON(http.StatusOK, tokenPair)
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "user was created but failed to login, please try to relogin",
		})
	}
}

func GithubOAuth(e *echo.Echo, stateStore *database.StateStore, database *database.Database, githubClientID string, githubClientSecret string, jwtSecret []byte, refreshTokenSecret []byte) {
	e.POST("/github/state", addNewState(stateStore))
	e.POST("/github/login", githubLogin(stateStore, database, githubClientID, githubClientSecret, jwtSecret, refreshTokenSecret))
}

var ErrTokenGeneration = errors.New("failed to generate tokens")

func generateTokenPair(user *database.User, jwtSecret []byte, refreshSecret []byte) (*map[string]interface{}, error) {
	claims := &UserTokenClaims{
		user.ID,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(AccessTokenExpiryMinutes)).Unix(),
		},
	}
	refreshClaims := &UserTokenClaims{
		user.ID,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(RefreshTokenExpiryHours)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("failed to generate a token: %s\n", err.Error())

		return nil, ErrTokenGeneration
	}
	rT, err := rToken.SignedString(refreshSecret)
	if err != nil {
		log.Printf("failed to generate a token: %s\n", err.Error())

		return nil, ErrTokenGeneration
	}

	tokenPair := map[string]interface{}{
		"token":        t,
		"refreshToken": rT,
	}

	return &tokenPair, nil
}

func Refresh(e *echo.Group, db *database.Database, jwtSecret []byte, refreshSecret []byte) {
	e.GET("", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return AssertionError("*jwt.Token")
		}

		claims, ok := token.Claims.(*UserTokenClaims)
		if !ok {
			return AssertionError("*UserTokenClaims")
		}

		user := db.GetUser(claims.Username)
		claims = &UserTokenClaims{
			user.ID,
			user.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(AccessTokenExpiryMinutes)).Unix(),
			},
		}
		refreshClaims := &UserTokenClaims{
			user.ID,
			user.Username,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(RefreshTokenExpiryHours)).Unix(),
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
