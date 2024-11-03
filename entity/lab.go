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
	AverageVolume int       `gorm:"average_volume"`
	Timestamp     time.Time `gorm:"column:timestamp;type:date"`
}

type LabCueHistory struct {
	ID        uint `gorm:"primary_key"`
	LabID     uint
	Lab       Lab
	Timestamp time.Time `gorm:"column:timestamp;type:date"`
}
