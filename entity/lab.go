package entity

import "time"

type Lab struct {
	ID                uint `gorm:"primary_key"`
	SubjectID         uint
	Subject           Subject
	StartDate         time.Time `gorm:"column:start_date;type:date"`
	LabSleepHistories []LabSleepHistory
	LabCueHistories   []LabCueHistory
}

type LabSleepHistory struct {
	ID            uint `gorm:"primary_key"`
	LabID         uint
	Lab           Lab
	AverageVolume float64   `gorm:"average_volume"`
	Timestamp     time.Time `gorm:"column:timestamp;type:datetime(3)"`
}

type LabCueHistory struct {
	ID         uint `gorm:"primary_key"`
	LabID      uint
	Lab        Lab
	Timestamp  time.Time `gorm:"column:timestamp;type:datetime(3)"`
	TargetWord string    `gorm:"column:target_word"`
}

type LabTest struct {
	ID        uint `gorm:"primary_key"`
	LabID     uint
	Lab       Lab
	StartDate time.Time `gorm:"column:start_date;type:date"`
	LabType   string    `gorm:"column:lab_type"`
}

type LabTestHistory struct {
	ID          uint `gorm:"primary_key"`
	LabTestID   uint
	LabTest     LabTest
	Word        string `gorm:"column:word"`
	WriitenWord string `gorm:"column:wriiten_word"`
	IsCorrect   bool   `gorm:"column:is_correct"`
}

type LabCueTargetWord struct {
	ID    uint `gorm:"primary_key"`
	LabID uint
	Lab   Lab
	Word  string `gorm:"column:word"`
}
