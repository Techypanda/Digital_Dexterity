package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"techytechster/digitaldexterity/internal/database"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AssessmentPayload struct {
	WillingnessToLearn      uint `json:"willingnessToLearn" form:"willingnessToLearn" query:"willingnessToLearn" validate:"max=100"`
	SelfSufficientLearning  uint `json:"selfSufficientLearning" form:"selfSufficientLearning" query:"selfSufficientLearning" validate:"max=100"`
	ImprovingCapability     uint `json:"improvingCapability" form:"improvingCapability" query:"improvingCapability" validate:"max=100"`
	InnovativeThinking      uint `json:"innovativeThinking" form:"innovativeThinking" query:"innovativeThinking" validate:"max=100"`
	GrowthMindset           uint `json:"growthMindset" form:"growthMindset" query:"growthMindset" validate:"max=100"`
	AwarenessOfSelfEfficacy uint `json:"awarenessOfSelfEfficacy" form:"awarenessOfSelfEfficacy" query:"awarenessOfSelfEfficacy" validate:"max=100"`
	ApplyingWhatTheyLearn   uint `json:"applyingWhatTheyLearn" form:"applyingWhatTheyLearn" query:"applyingWhatTheyLearn" validate:"max=100"`
	Adaptability            uint `json:"adaptability" form:"adaptability" query:"adaptability" validate:"max=100"`
}
type ExternalAssessment struct {
	WillingnessToLearn      uint   `json:"willingnessToLearn" form:"willingnessToLearn" query:"willingnessToLearn" validate:"max=100"`
	SelfSufficientLearning  uint   `json:"selfSufficientLearning" form:"selfSufficientLearning" query:"selfSufficientLearning" validate:"max=100"`
	ImprovingCapability     uint   `json:"improvingCapability" form:"improvingCapability" query:"improvingCapability" validate:"max=100"`
	InnovativeThinking      uint   `json:"innovativeThinking" form:"innovativeThinking" query:"innovativeThinking" validate:"max=100"`
	GrowthMindset           uint   `json:"growthMindset" form:"growthMindset" query:"growthMindset" validate:"max=100"`
	AwarenessOfSelfEfficacy uint   `json:"awarenessOfSelfEfficacy" form:"awarenessOfSelfEfficacy" query:"awarenessOfSelfEfficacy" validate:"max=100"`
	ApplyingWhatTheyLearn   uint   `json:"applyingWhatTheyLearn" form:"applyingWhatTheyLearn" query:"applyingWhatTheyLearn" validate:"max=100"`
	Adaptability            uint   `json:"adaptability" form:"adaptability" query:"adaptability" validate:"max=100"`
	Assessing               string `json:"assessing" form:"assessing" query:"assessing" validate:"required"`
}

func selfAssess(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return AssertionError("*jwt.Token")
		}

		payload := new(AssessmentPayload)

		if err := c.Bind(payload); err != nil {
			log.Printf("failed to bind payload to self assessment: %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"asssessedSelf": false,
				"error":         err.Error(),
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("failed to validate to self assessment: %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"asssessedSelf": false,
				"error":         fmt.Sprintf("Expected WillingnessToLearn SelfSufficientLearning ImprovingCapability InnovativeThinking GrowthMindset AwarenessOfSelfEfficacy ApplyingWhatTheyLearn Adaptability: %s", err.Error()),
			})
		}

		claims, ok := token.Claims.(*UserTokenClaims)
		if !ok {
			return AssertionError("*UserTokenClaims")
		}

		user := db.GetUser(claims.Username)

		if user == nil {
			log.Printf("fatal! a username is missing from db: %s\n", claims.Username)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":         "please contact admin, your username is gone",
				"asssessedSelf": false,
			})
		}

		log.Printf("adding self assessment to db\n")

		err := db.AddSelfAssessment(*user, database.DigitalDexterityAssessment{
			WillingnessToLearn:      payload.WillingnessToLearn,
			SelfSufficientLearning:  payload.SelfSufficientLearning,
			ImprovingCapability:     payload.ImprovingCapability,
			InnovativeThinking:      payload.InnovativeThinking,
			GrowthMindset:           payload.GrowthMindset,
			AwarenessOfSelfEfficacy: payload.AwarenessOfSelfEfficacy,
			ApplyingWhatTheyLearn:   payload.ApplyingWhatTheyLearn,
			Adaptability:            payload.Adaptability,
		})

		if err != nil {
			log.Printf("failed to create self assessment, %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":         fmt.Sprintf("failed to create self assessment: %s", err.Error()),
				"asssessedSelf": false,
			})
		}

		log.Printf("successfully created self assessment\n")

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"asssessedSelf": true,
		})
	}
}

func externalAssess(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		var payload ExternalAssessment
		if err := c.Bind(&payload); err != nil {
			log.Printf("failed to bind payload to external assessment: %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"assessed": false,
				"error":    err.Error(),
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("failed to validate to external assessment: %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"assessed": false,
				"error":    fmt.Sprintf("Expected WillingnessToLearn SelfSufficientLearning ImprovingCapability InnovativeThinking GrowthMindset AwarenessOfSelfEfficacy ApplyingWhatTheyLearn Adaptability Assessed: %s", err.Error()),
			})
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return AssertionError("*jwt.Token")
		}

		claims, ok := token.Claims.(*UserTokenClaims)
		if !ok {
			return AssertionError("*UserTokenClaims")
		}

		user := db.GetUser(claims.Username)

		if user == nil {
			log.Printf("fatal! username is missing from db: %s\n", claims.Username)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":    "please contact admin, your username is gone",
				"assessed": false,
			})
		}

		assessingUser := db.GetUser(payload.Assessing)
		if assessingUser == nil {
			log.Printf("assessed user does not exist: %s\n", payload.Assessing)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":    "unknown user",
				"assessed": false,
			})
		}

		if assessingUser.ID == user.ID {
			log.Printf("someone tried to assess themselves: %s\n", payload.Assessing)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":    "you cannot assess yourself",
				"assessed": false,
			})
		}

		err := db.AddExternalAssessment(*user, *assessingUser, database.DigitalDexterityAssessment{
			WillingnessToLearn:      payload.WillingnessToLearn,
			SelfSufficientLearning:  payload.SelfSufficientLearning,
			ImprovingCapability:     payload.ImprovingCapability,
			InnovativeThinking:      payload.InnovativeThinking,
			GrowthMindset:           payload.GrowthMindset,
			AwarenessOfSelfEfficacy: payload.AwarenessOfSelfEfficacy,
			ApplyingWhatTheyLearn:   payload.ApplyingWhatTheyLearn,
			Adaptability:            payload.Adaptability,
		})
		if err != nil {
			log.Printf("failed to external assess: %s\n", err.Error())

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":    fmt.Sprintf("failed to create external assessment: %s", err.Error()),
				"assessed": false,
			})
		}

		log.Printf("externally assessed: %s\n", assessingUser.Username)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"assessed": true,
		})
	}
}

var ErrAssertion = errors.New("failed to assert type")

func AssertionError(expected string) error {
	return fmt.Errorf("failed to assert type %s : %w", expected, ErrAssertion)
}

func getAssessments(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return AssertionError("*jwt.Token")
		}

		claims, ok := token.Claims.(*UserTokenClaims)
		if !ok {
			return AssertionError("*UserTokenClaims")
		}

		user := db.GetUser(claims.Username)
		if user == nil {
			log.Printf("fatal! username is missing from db: %s\n", claims.Username)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":         "please contact admin, your username is gone",
				"asssessedSelf": false,
			})
		}

		selfAssessment, _ := db.GetSelfAssessment(*user)
		externalAssessments := db.GetExternalAssessments(*user)

		if selfAssessment != nil {
			log.Printf("returning self assessments and external assessments for user: %s\n", claims.Username)

			return c.JSON(http.StatusOK, map[string]interface{}{
				"selfAssessment":      selfAssessment.DigitalDexterityAssessment,
				"externalAssessments": externalAssessments,
			})
		}

		log.Printf("returning external assessments for user: %s\n", claims.Username)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"externalAssessments": externalAssessments,
		})
	}
}

func Assessment(e *echo.Group, db *database.Database, jwtSecret []byte) {
	e.POST("/assess/me", selfAssess(db))
	e.POST("/assess", externalAssess(db))
	e.GET("/assess", getAssessments(db))
}
