package api

import (
	"fmt"
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

func selfAssess(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		payload := new(AssessmentPayload)
		if err := c.Bind(payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   err.Error(),
			})
		} else if err := c.Validate(payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"created": false,
				"error":   fmt.Sprintf("Expected WillingnessToLearn SelfSufficientLearning ImprovingCapability InnovativeThinking GrowthMindset AwarenessOfSelfEfficacy ApplyingWhatTheyLearn Adaptability: %s", err.Error()),
			})
		}
		claims := token.Claims.(*UserTokenClaims)
		user := db.GetUser(claims.Username)
		if user == nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":         "please contact admin, your username is gone",
				"asssessedSelf": false,
			})
		} else {
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
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":         fmt.Sprintf("failed to create self assessment: %s", err.Error()),
					"asssessedSelf": false,
				})
			} else {
				return c.JSON(http.StatusCreated, map[string]interface{}{
					"asssessedSelf": true,
				})
			}
		}
	}
}

func getAssessments(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*UserTokenClaims)
		user := db.GetUser(claims.Username)
		if user == nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":         "please contact admin, your username is gone",
				"asssessedSelf": false,
			})
		} else {
			selfAssessment, _ := db.GetSelfAssessment(*user)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"selfAssessment": selfAssessment.DigitalDexterityAssessment,
			})
		}
	}
}

func Assessment(e *echo.Group, db *database.Database, jwtSecret []byte) {
	e.POST("/assess/me", selfAssess(db))
	e.GET("/assess", getAssessments(db))
}
