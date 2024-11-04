package model

import (
	"time"
	"tmr-backend/entity"

	"gorm.io/gorm"
)

type LabModel interface {
	CreateBreathingHistory(idForLogin string, averageVolume int, timestamp time.Time) error
	GetLabBySubjectIdForLogin(idForLogin string) (*entity.Lab, error)
}

type labModel struct {
	db *gorm.DB
}

func NewLabModel(db *gorm.DB) LabModel {
	return &labModel{
		db: db,
	}
}

func (m *labModel) CreateBreathingHistory(idForLogin string, averageVolume int, timestamp time.Time) error {
	lab, err := m.GetLabBySubjectIdForLogin(idForLogin)

	if err != nil {
		return err
	}

	breathingHistory := &entity.LabSleepHistory{
		LabID:         lab.ID,
		AverageVolume: averageVolume,
		Timestamp:     timestamp,
	}
	if err := m.db.Save(breathingHistory).Error; err != nil {
		return err
	}

	return nil
}

func (m *labModel) GetLabBySubjectIdForLogin(idForLogin string) (*entity.Lab, error) {
	lab := &entity.Lab{}

	if err := m.db.Joins("Subject", m.db.Where(&entity.Subject{
		IdForLogin: idForLogin,
	})).First(lab).Error; err != nil {
		return nil, err
	}

	return lab, nil
}
