package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileUtil interface {
	CreateTempCSVFile(content string) (string, error)
	StartCleanupRoutine()
}

type fileUtil struct {
	tempDir     string
	fileExpiry  time.Duration
	cleanupLock sync.Mutex
}

func NewFileUtil() FileUtil {
	util := &fileUtil{
		tempDir:    "./temp_files",
		fileExpiry: 24 * time.Hour,
	}

	// 임시 디렉토리 생성
	if err := os.MkdirAll(util.tempDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create temp directory: %v", err))
	}

	return util
}

func (u *fileUtil) CreateTempCSVFile(content string) (string, error) {
	// 현재 시간을 파일명에 포함
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("target_words_%s.csv", timestamp)
	filepath := filepath.Join(u.tempDir, filename)

	// UTF-8 BOM 추가
	bom := []byte{0xEF, 0xBB, 0xBF}
	data := append(bom, []byte(content)...)

	// 파일 생성 및 내용 작성
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (u *fileUtil) StartCleanupRoutine() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			u.cleanupExpiredFiles()
		}
	}()
}

func (u *fileUtil) cleanupExpiredFiles() {
	u.cleanupLock.Lock()
	defer u.cleanupLock.Unlock()

	files, err := os.ReadDir(u.tempDir)
	if err != nil {
		fmt.Printf("Error reading temp directory: %v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(u.tempDir, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		// 파일이 만료되었는지 확인
		if time.Since(fileInfo.ModTime()) > u.fileExpiry {
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Error removing expired file %s: %v\n", filePath, err)
			}
		}
	}
}
