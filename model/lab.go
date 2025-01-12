package model

import (
	"time"
	"tmr-backend/dto"
	"tmr-backend/entity"

	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type LabModel interface {
	CreateBreathingHistory(idForLogin string, averageVolume int, timestamp time.Time) error
	GetLabBySubjectIdForLogin(idForLogin string) (*entity.Lab, error)
	CreateCueHistory(idForLogin string, timestamp time.Time, targetWord string) error
	CreatePreTest(labID uint, results []dto.StartLabRequestResult) error
	CreateTest(labID uint) error
	GetLabTestByIdForLogin(idForLogin string, labType string) (*entity.LabTest, error)
	CreatePreTestHistory(createPreTestHistoryDto dto.CreatePreTestHistoryDto) ([]string, int, int, error)
	CreateTestHistory(createPreTestHistoryDto dto.CreatePreTestHistoryDto) error
	GetTargetWordsByLabId(labID uint) ([]string, error)
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
	lab := entity.Lab{}

	if err := m.db.
		Preload("Subject").
		Joins("Subject").
		Where("Subject.id_for_login = ?", idForLogin).
		First(&lab).Error; err != nil {
		return nil, err
	}

	return &lab, nil
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

func (m *labModel) CreatePreTest(labID uint, results []dto.StartLabRequestResult) error {
	// 트랜잭션 시작
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// pretest 존재 여부 확인
	var labTest entity.LabTest
	err := tx.Where("lab_id = ? AND lab_type = ?", labID, "pretest").
		First(&labTest).Error

	if err == gorm.ErrRecordNotFound {
		// pretest가 없으면 새로 생성
		labTest = entity.LabTest{
			LabID:     labID,
			StartDate: time.Now(),
			LabType:   "pretest",
		}
		if err := tx.Create(&labTest).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		tx.Rollback()
		return err
	}

	// 기존 LabTestHistory 삭제
	if err := tx.Where("lab_test_id = ?", labTest.ID).
		Delete(&entity.LabTestHistory{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 새로운 LabTestHistory 생성
	for _, result := range results {
		testHistory := entity.LabTestHistory{
			LabTestID:   labTest.ID,
			Word:        result.Word,
			WriitenWord: result.WrittenWord,
		}
		if err := tx.Create(&testHistory).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 트랜잭션 커밋
	return tx.Commit().Error
}

func (m *labModel) CreateTest(labID uint) error {
	var exists bool
	err := m.db.Model(&entity.LabTest{}).
		Select("1").
		Where("lab_id = ? AND lab_type = ?", labID, "test").
		Limit(1).
		Find(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		test := &entity.LabTest{
			LabID:     labID,
			StartDate: time.Now(),
			LabType:   "test",
		}

		if err := m.db.Save(test).Error; err != nil {
			return err
		}
	}

	return nil
}

func (m *labModel) GetLabTestByIdForLogin(idForLogin string, labType string) (*entity.LabTest, error) {
	lab, err := m.GetLabBySubjectIdForLogin(idForLogin)
	if err != nil {
		return nil, err
	}

	var labTest entity.LabTest

	if err := m.db.Where(&entity.LabTest{
		LabID:   lab.ID,
		LabType: labType,
	}).First(&labTest).Error; err != nil {
		return nil, err
	}

	return &labTest, nil
}

func (m *labModel) CreatePreTestHistory(createPreTestHistoryDto dto.CreatePreTestHistoryDto) ([]string, int, int, error) {
	tx := m.db.Begin()
	if tx.Error != nil {
		return nil, 0, 0, tx.Error
	}

	if err := tx.Where(&entity.LabTestHistory{
		LabTestID: createPreTestHistoryDto.LabTestID,
	}).Delete(&entity.LabTestHistory{}).Error; err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}

	correctTestHistories := []*entity.LabTestHistory{}
	wrongTestHistories := []*entity.LabTestHistory{}
	for _, result := range createPreTestHistoryDto.Results {
		testHistory := entity.LabTestHistory{
			LabTestID: createPreTestHistoryDto.LabTestID,
			Word:      result.Word,
			IsCorrect: result.IsCorrect,
		}

		if err := tx.Save(&testHistory).Error; err != nil {
			tx.Rollback()
			return nil, 0, 0, err
		}

		if result.IsCorrect {
			correctTestHistories = append(correctTestHistories, &testHistory)
		} else {
			wrongTestHistories = append(wrongTestHistories, &testHistory)
		}
	}

	labTest := entity.LabTest{}
	if err := tx.Where(&entity.LabTest{
		ID: createPreTestHistoryDto.LabTestID,
	}).First(&labTest).Error; err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}

	if err := tx.Where(&entity.LabCueTargetWord{
		LabID: labTest.LabID,
	}).Delete(&entity.LabCueTargetWord{}).Error; err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}

	numberOfCorrectWord := len(correctTestHistories)
	var numberOfWrongWord int
	if len(correctTestHistories)%3 != 0 {
		for {
			numberOfCorrectWord--

			if numberOfCorrectWord%3 == 0 {
				numberOfCorrectWord /= 3
				numberOfCorrectWord *= 2
				break
			}
		}
	} else {
		numberOfCorrectWord /= 3
		numberOfCorrectWord *= 2
	}
	numberOfWrongWord = 80 - numberOfCorrectWord

	targetWords := m.PickCueTargetWords(correctTestHistories, numberOfCorrectWord, wrongTestHistories, numberOfWrongWord)
	selectedWords := make([]string, len(targetWords))

	for i, targetWord := range targetWords {
		cueWord := entity.LabCueTargetWord{
			LabID: labTest.LabID,
			Word:  targetWord.Word,
		}
		selectedWords[i] = targetWord.Word

		if err := tx.Save(&cueWord).Error; err != nil {
			tx.Rollback()
			return nil, 0, 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, 0, err
	}

	return selectedWords, numberOfCorrectWord, numberOfWrongWord, nil
}

func (m *labModel) CreateTestHistory(createPreTestHistoryDto dto.CreatePreTestHistoryDto) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where(&entity.LabTestHistory{
		LabTestID: createPreTestHistoryDto.LabTestID,
	}).Delete(&entity.LabTestHistory{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, result := range createPreTestHistoryDto.Results {
		testHistory := entity.LabTestHistory{
			LabTestID: createPreTestHistoryDto.LabTestID,
			Word:      result.Word,
			IsCorrect: result.IsCorrect,
		}

		if err := tx.Save(&testHistory).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (m *labModel) PickCueTargetWords(corrects []*entity.LabTestHistory, numberOfCorrectWord int, wrongs []*entity.LabTestHistory, numberOfWrongWord int) []*entity.LabTestHistory {
	results := make([]*entity.LabTestHistory, numberOfCorrectWord+numberOfWrongWord)

	// 각 배열에서 랜덤하게 선택
	rand.Seed(uint64(time.Now().UnixNano()))

	// 정답 단어 선택
	correctIndices := rand.Perm(len(corrects))[:numberOfCorrectWord]
	for i, idx := range correctIndices {
		results[i] = corrects[idx]
	}

	// 오답 단어 선택
	wrongIndices := rand.Perm(len(wrongs))[:numberOfWrongWord]
	for i, idx := range wrongIndices {
		results[numberOfCorrectWord+i] = wrongs[idx]
	}

	// 결과 배열을 섞음
	rand.Shuffle(len(results), func(i, j int) {
		results[i], results[j] = results[j], results[i]
	})

	return results
}

func (m *labModel) GetTargetWordsByLabId(labID uint) ([]string, error) {
	var targetWords []entity.LabCueTargetWord

	if err := m.db.Where(&entity.LabCueTargetWord{
		LabID: labID,
	}).Find(&targetWords).Error; err != nil {
		return nil, err
	}

	targetWordsString := make([]string, len(targetWords))
	for i, v := range targetWords {
		targetWordsString[i] = v.Word
	}

	return targetWordsString, nil
}
