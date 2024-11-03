package entity

type Subject struct {
	ID           uint   `gorm:"primary_key"`
	IdForLogin   string `gorm:"id_for_login"`
	Age          uint   `gorm:"age"`
	EnglishLevel string `gorm:"english_level"`
	Detail       string `gorm:"detail;type:text"`
}
