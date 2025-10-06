package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zpwoot/internal/core/ports/output"
)

type MediaProcessor struct {
	httpClient *http.Client
}

func NewMediaProcessor() *MediaProcessor {
	return &MediaProcessor{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (mp *MediaProcessor) ProcessMedia(file, mimeType, fileName string) (*output.MediaData, error) {
	var data []byte
	var err error
	var detectedMimeType string
	var detectedFileName string
	if mp.isURL(file) {
		data, detectedMimeType, detectedFileName, err = mp.downloadFromURL(file)
		if err != nil {
			return nil, fmt.Errorf("failed to download from URL: %w", err)
		}
	} else if mp.isFilePath(file) {
		data, detectedMimeType, detectedFileName, err = mp.readFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	} else if mp.isBase64(file) {
		data, err = mp.decodeBase64(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
	} else {
		return nil, fmt.Errorf("invalid file input: must be URL, file path, or base64")
	}
	finalMimeType := mimeType
	if finalMimeType == "" {
		finalMimeType = detectedMimeType
	}
	if finalMimeType == "" {
		finalMimeType = "application/octet-stream"
	}
	finalFileName := fileName
	if finalFileName == "" {
		finalFileName = detectedFileName
	}

	return &output.MediaData{
		MimeType: finalMimeType,
		Data:     data,
		FileName: finalFileName,
	}, nil
}
func (mp *MediaProcessor) isURL(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}
func (mp *MediaProcessor) isFilePath(input string) bool {

	if strings.Contains(input, "/") || strings.Contains(input, "\\") {
		if _, err := os.Stat(input); err == nil {
			return true
		}
	}
	return false
}
func (mp *MediaProcessor) isBase64(input string) bool {

	if strings.Contains(input, ",") {
		parts := strings.Split(input, ",")
		if len(parts) > 1 {
			input = parts[1]
		}
	}
	_, err := base64.StdEncoding.DecodeString(input)
	return err == nil && len(input) > 10
}
func (mp *MediaProcessor) downloadFromURL(url string) ([]byte, string, string, error) {
	resp, err := mp.httpClient.Get(url)
	if err != nil {
		return nil, "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", err
	}
	mimeType := resp.Header.Get("Content-Type")
	if mimeType != "" {

		if idx := strings.Index(mimeType, ";"); idx != -1 {
			mimeType = mimeType[:idx]
		}
		mimeType = strings.TrimSpace(mimeType)
	}
	fileName := filepath.Base(url)
	if strings.Contains(fileName, "?") {
		fileName = strings.Split(fileName, "?")[0]
	}

	return data, mimeType, fileName, nil
}
func (mp *MediaProcessor) readFromFile(filePath string) ([]byte, string, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", "", err
	}
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	fileName := filepath.Base(filePath)

	return data, mimeType, fileName, nil
}
func (mp *MediaProcessor) decodeBase64(input string) ([]byte, error) {

	base64Data := input
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	return base64.StdEncoding.DecodeString(base64Data)
}
func (mp *MediaProcessor) DetectMimeTypeFromData(data []byte) string {
	return http.DetectContentType(data)
}
