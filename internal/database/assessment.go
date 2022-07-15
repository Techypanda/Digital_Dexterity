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

type ExternalAssessment struct {
	gorm.Model
	DigitalDexterityAssessment DigitalDexterityAssessment `gorm:"embedded"`
	AssessedBy                 uint
	UserID                     uint
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

func (db *Database) AddExternalAssessment(user User, assess User, assessment DigitalDexterityAssessment) error {
	externalAssessment := ExternalAssessment{
		DigitalDexterityAssessment: assessment,
		AssessedBy:                 user.ID,
		UserID:                     assess.ID,
	}
	db.db.Where("assessed_by = ? and user_id = ?", user.ID, assess.ID).Delete(&ExternalAssessment{})
	result := db.db.Create(&externalAssessment)
	return result.Error
}

func (db *Database) GetExternalAssessments(user User) []DigitalDexterityAssessment {
	var externalAssessments []ExternalAssessment
	var assessments []DigitalDexterityAssessment
	db.db.Where("user_id = ?", user.ID).Find(&externalAssessments)
	for _, assessment := range externalAssessments {
		assessments = append(assessments, assessment.DigitalDexterityAssessment)
	}
	return assessments
}

func (db *Database) GetSelfAssessment(user User) (*SelfAssessment, error) {
	var assessment SelfAssessment
	db.db.Where("user_id = ?", user.ID).First(&assessment)
	if assessment.UserID != user.ID {
		return nil, errors.New("there is no self assessment")
	}
	return &assessment, nil
}
