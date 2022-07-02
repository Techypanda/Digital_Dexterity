package database

import (
	"errors"

	"gorm.io/gorm"
)

type DigitalDexterityAssessment struct {
	WillingnessToLearn      uint
	SelfSufficientLearning  uint
	ImprovingCapability     uint
	InnovativeThinking      uint
	GrowthMindset           uint
	AwarenessOfSelfEfficacy uint
	ApplyingWhatTheyLearn   uint
	Adaptability            uint
}

type SelfAssessment struct {
	gorm.Model
	DigitalDexterityAssessment DigitalDexterityAssessment `gorm:"embedded"`
	UserID                     uint
}

func (db *Database) AddSelfAssessment(user User, assessment DigitalDexterityAssessment) error {
	selfAssessment := SelfAssessment{
		DigitalDexterityAssessment: assessment,
		UserID:                     user.ID,
	}
	db.db.Where("user_id = ?", user.ID).Delete(&SelfAssessment{})
	result := db.db.Create(&selfAssessment)
	return result.Error
}

func (db *Database) GetSelfAssessment(user User) (*SelfAssessment, error) {
	var assessment SelfAssessment
	db.db.Where("user_id = ?", user.ID).First(&assessment)
	if assessment.UserID != user.ID {
		return nil, errors.New("there is no self assessment")
	} else {
		return &assessment, nil
	}
}
