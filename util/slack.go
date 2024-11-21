package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tmr-backend/dto"
)

type SlackUtil interface {
	SendPreTestStartMessage(labID string) error
	SendTestStartMessage(labID string) error
	SendTestResultMessage(labID string, correctCount int, wrongCount int, words []string) error
}

type slackUtil struct {
	fileUtil FileUtil
	baseURL  string
}

func NewSlackUtil(fileUtil FileUtil, baseURL string) SlackUtil {
	return &slackUtil{
		fileUtil: fileUtil,
		baseURL:  baseURL,
	}
}

func (u *slackUtil) SendPreTestStartMessage(labID string) error {
	url := os.Getenv("SLACK_WEBHOOK_URL")
	message := dto.SlackMessagePayload{
		Text: fmt.Sprintf("%s의 사전 테스트가 시작되었습니다. 아래 링크를 클릭해 테스트를 수행하세요.\n <https://junho100.github.io/tmr-lab-web/%s/pretest-result>", labID, labID),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("에러 응답: %d - %s\n", resp.StatusCode, string(body))
		return errors.New("status code of slack reuqest is not 200")
	}

	return nil
}

func (u *slackUtil) SendTestStartMessage(labID string) error {
	url := os.Getenv("SLACK_WEBHOOK_URL")
	message := dto.SlackMessagePayload{
		Text: fmt.Sprintf("%s의 사후 테스트가 시작되었습니다. 아래 링크를 클릭해 테스트를 수행하세요.\n <https://junho100.github.io/tmr-lab-web/%s/test-result>", labID, labID),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("에러 응답: %d - %s\n", resp.StatusCode, string(body))
		return errors.New("status code of slack reuqest is not 200")
	}

	return nil
}

func (u *slackUtil) SendTestResultMessage(labID string, correctCount int, wrongCount int, words []string) error {
	// CSV 파일 생성
	csvContent := "Word\n"
	for _, word := range words {
		csvContent += word + "\n"
	}

	// 임시 파일 생성
	filename, err := u.fileUtil.CreateTempCSVFile(csvContent)
	if err != nil {
		return err
	}

	// 다운로드 URL 생성
	fileUrl := fmt.Sprintf("%s/api/files/%s", u.baseURL, filename)

	url := os.Getenv("SLACK_WEBHOOK_URL")
	message := dto.SlackMessagePayload{
		Blocks: []dto.SlackMessageBlock{
			{
				Type: "section",
				Text: &dto.SlackMessageText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%s의 사전 테스트 결과가 저장되었습니다.*\n\n"+
						"• 맞은 단어 중 선정된 개수: %d개\n"+
						"• 틀린 단어 중 선정된 개수: %d개\n\n"+
						"아래 링크에서 60분 동안 선정된 단어 목록을 다운로드할 수 있습니다:\n%s",
						labID, correctCount, wrongCount, fileUrl),
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code of slack request is not 200: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}
