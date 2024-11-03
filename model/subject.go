package model

import "gorm.io/gorm"

type SubjectModel interface {
}

type subjectModel struct {
	db *gorm.DB
}

func NewSubjectModel(db *gorm.DB) SubjectModel {
	return &subjectModel{
		db: db,
	}
}
