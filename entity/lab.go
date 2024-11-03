package entity

import "time"

type Lab struct {
	ID        uint `gorm:"primary_key"`
	SubjectID uint
	Subject   Subject
	StartDate time.Time `gorm:"column:start_date;type:date"`
}

// type LabSleepHistory struct {
// 	ID uint `gorm`
// }
