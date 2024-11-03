package model

import (
	"fmt"
	"strconv"
	"strings"
	"tmr-backend/entity"

	"gorm.io/gorm"
)

type SubjectModel interface {
	CreateSubject(age uint, englishLevel string, detail string) (string, error)
	FindSubjectByIdForLogin(idForLogin string) (*entity.Subject, error)
}

type subjectModel struct {
	db *gorm.DB
}

func NewSubjectModel(db *gorm.DB) SubjectModel {
	return &subjectModel{
		db: db,
	}
}

func (m *subjectModel) CreateSubject(age uint, englishLevel string, detail string) (string, error) {
	tx := m.db.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	subject := &entity.Subject{
		Age:          age,
		EnglishLevel: strings.ToLower(englishLevel),
		Detail:       detail,
	}
	if err := tx.Save(subject).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	idForLogin := fmt.Sprintf("%s%s", "lab", strconv.Itoa(int(subject.ID)))
	subject.IdForLogin = idForLogin
	if err := tx.Save(subject).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	return idForLogin, tx.Commit().Error
}

func (m *subjectModel) FindSubjectByIdForLogin(idForLogin string) (*entity.Subject, error) {
	var subject entity.Subject

	if err := m.db.Where(&entity.Subject{
		IdForLogin: idForLogin,
	}).First(&subject).Error; err != nil {
		return nil, err
	}

	return &subject, nil
}
