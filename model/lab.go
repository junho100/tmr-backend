package model

import (
	"time"
	"tmr-backend/dto"
	"tmr-backend/entity"

	"gorm.io/gorm"
)

type LabModel interface {
	CreateBreathingHistory(idForLogin string, averageVolume int, timestamp time.Time) error
	GetLabBySubjectIdForLogin(idForLogin string) (*entity.Lab, error)
	CreateCueHistory(idForLogin string, timestamp time.Time, targetWord string) error
	CreatePreTest(labID uint) error
	GetLabTestByIdForLogin(idForLogin string) (*entity.LabTest, error)
	CreateTestHistory(createTestHistoryDto dto.CreateTestHistoryDto) error
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

func (m *labModel) CreateCueHistory(idForLogin string, timestamp time.Time, targetWord string) error {
	lab, err := m.GetLabBySubjectIdForLogin(idForLogin)

	if err != nil {
		return err
	}

	cueHistroy := &entity.LabCueHistory{
		LabID:      lab.ID,
		Timestamp:  timestamp,
		TargetWord: targetWord,
	}
	if err := m.db.Save(cueHistroy).Error; err != nil {
		return err
	}

	return nil
}

func (m *labModel) CreatePreTest(labID uint) error {
	preTest := &entity.LabTest{
		LabID:     labID,
		StartDate: time.Now(),
		LabType:   "pretest",
	}

	if err := m.db.Save(preTest).Error; err != nil {
		return err
	}

	return nil
}

func (m *labModel) GetLabTestByIdForLogin(idForLogin string) (*entity.LabTest, error) {
	lab, err := m.GetLabBySubjectIdForLogin(idForLogin)
	if err != nil {
		return nil, err
	}

	var labTest entity.LabTest

	if err := m.db.Where(&entity.LabTest{
		ID: lab.ID,
	}).First(&labTest).Error; err != nil {
		return nil, err
	}

	return &labTest, nil
}

func (m *labModel) CreateTestHistory(createTestHistoryDto dto.CreateTestHistoryDto) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where(&entity.LabTestHistory{
		LabTestID: createTestHistoryDto.LabTestID,
	}).Delete(&entity.LabTestHistory{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, result := range createTestHistoryDto.Results {
		testHistory := entity.LabTestHistory{
			LabTestID: createTestHistoryDto.LabTestID,
			Word:      result.Word,
			IsCorrect: result.IsCorrect,
		}

		if err := tx.Save(&testHistory).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
