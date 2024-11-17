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
	SendTestStartMessage(labID string) error
}

type slackUtil struct {
}

func NewSlackUtil() SlackUtil {
	return &slackUtil{}
}

func (u *slackUtil) SendTestStartMessage(labID string) error {
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
