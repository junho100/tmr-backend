package model

import (
	"testing"
	"time"

	"tmr-backend/entity" // 실제 프로젝트 경로로 수정 필요

	"github.com/stretchr/testify/assert"
)

// TestPickCueTargetWords tests the PickCueTargetWords method of labModel
func TestPickCueTargetWords(t *testing.T) {
	// 테스트용 모델 인스턴스 생성
	m := &labModel{}

	// 테스트 데이터 생성
	corrects := createTestHistories(5, true)
	wrongs := createTestHistories(5, false)

	t.Run("정상적인 케이스", func(t *testing.T) {
		numberOfCorrectWord := 2
		numberOfWrongWord := 3

		result := m.PickCueTargetWords(corrects, numberOfCorrectWord, wrongs, numberOfWrongWord)

		// 결과 검증
		assert.Equal(t, numberOfCorrectWord+numberOfWrongWord, len(result), "전체 결과 길이가 일치해야 함")

		// 정답과 오답 단어 개수 확인
		correctCount := countWordsByCorrect(result, true)
		wrongCount := countWordsByCorrect(result, false)

		assert.Equal(t, numberOfCorrectWord, correctCount, "정답 단어 개수가 일치해야 함")
		assert.Equal(t, numberOfWrongWord, wrongCount, "오답 단어 개수가 일치해야 함")
	})

	t.Run("랜덤성 검증", func(t *testing.T) {
		numberOfCorrectWord := 2
		numberOfWrongWord := 2

		resultSets := make(map[string]bool)

		for i := 0; i < 10; i++ {
			result := m.PickCueTargetWords(corrects, numberOfCorrectWord, wrongs, numberOfWrongWord)

			resultKey := ""
			for _, r := range result {
				resultKey += r.Word + ","
			}
			resultSets[resultKey] = true
		}
		assert.Greater(t, len(resultSets), 1, "여러 번 실행 시 서로 다른 결과가 나와야 함")
	})

	t.Run("순서 랜덤성 검증", func(t *testing.T) {
		numberOfCorrectWord := 3
		numberOfWrongWord := 2

		firstResult := m.PickCueTargetWords(corrects, numberOfCorrectWord, wrongs, numberOfWrongWord)

		isDifferentOrder := false
		for i := 0; i < 10; i++ {
			currentResult := m.PickCueTargetWords(corrects, numberOfCorrectWord, wrongs, numberOfWrongWord)

			if !isSameOrder(firstResult, currentResult) {
				isDifferentOrder = true
				break
			}
		}

		assert.True(t, isDifferentOrder, "여러 번 실행 시 순서가 다른 결과가 나와야 함")
	})
}

// createTestHistories creates test LabTestHistory array
func createTestHistories(count int, isCorrect bool) []*entity.LabTestHistory {
	histories := make([]*entity.LabTestHistory, count)
	for i := 0; i < count; i++ {
		// 테스트용 Lab 생성
		lab := &entity.Lab{
			ID:        uint(i + 1),
			SubjectID: 1,
			StartDate: time.Now(),
		}

		// 테스트용 LabTest 생성
		labTest := &entity.LabTest{
			ID:        uint(i + 1),
			LabID:     lab.ID,
			Lab:       *lab,
			StartDate: time.Now(),
			LabType:   "test",
		}

		histories[i] = &entity.LabTestHistory{
			ID:          uint(i + 1),
			LabTestID:   labTest.ID,
			LabTest:     *labTest,
			Word:        "word_" + string(rune(i+'1')),
			WriitenWord: "written_" + string(rune(i+'1')),
			IsCorrect:   isCorrect,
		}
	}
	return histories
}

// countWordsByCorrect counts correct/incorrect words in the result array
func countWordsByCorrect(histories []*entity.LabTestHistory, isCorrect bool) int {
	count := 0
	for _, history := range histories {
		if history.IsCorrect == isCorrect {
			count++
		}
	}
	return count
}

// isSameOrder checks if two result arrays have the same order
func isSameOrder(a, b []*entity.LabTestHistory) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i].Word != b[i].Word {
			return false
		}
	}
	return true
}
