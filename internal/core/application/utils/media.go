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

// MediaProcessor handles different types of media input (URL, file path, base64)
type MediaProcessor struct {
	httpClient *http.Client
}

// NewMediaProcessor creates a new media processor
func NewMediaProcessor() *MediaProcessor {
	return &MediaProcessor{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProcessMedia processes a file input and returns MediaData
// file can be: URL, local file path, or base64 string
func (mp *MediaProcessor) ProcessMedia(file, mimeType, fileName string) (*output.MediaData, error) {
	var data []byte
	var err error
	var detectedMimeType string
	var detectedFileName string

	// Determine the type of input and process accordingly
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

	// Use provided mimeType if available, otherwise use detected
	finalMimeType := mimeType
	if finalMimeType == "" {
		finalMimeType = detectedMimeType
	}
	if finalMimeType == "" {
		finalMimeType = "application/octet-stream"
	}

	// Use provided fileName if available, otherwise use detected
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

// isURL checks if the input is a URL
func (mp *MediaProcessor) isURL(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

// isFilePath checks if the input is a file path
func (mp *MediaProcessor) isFilePath(input string) bool {
	// Check if it's a valid file path and the file exists
	if strings.Contains(input, "/") || strings.Contains(input, "\\") {
		if _, err := os.Stat(input); err == nil {
			return true
		}
	}
	return false
}

// isBase64 checks if the input is base64 encoded
func (mp *MediaProcessor) isBase64(input string) bool {
	// Remove data URL prefix if present
	if strings.Contains(input, ",") {
		parts := strings.Split(input, ",")
		if len(parts) > 1 {
			input = parts[1]
		}
	}

	// Try to decode as base64
	_, err := base64.StdEncoding.DecodeString(input)
	return err == nil && len(input) > 10 // Minimum reasonable base64 length
}

// downloadFromURL downloads content from a URL
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

	// Get MIME type from Content-Type header
	mimeType := resp.Header.Get("Content-Type")
	if mimeType != "" {
		// Remove charset and other parameters
		if idx := strings.Index(mimeType, ";"); idx != -1 {
			mimeType = mimeType[:idx]
		}
		mimeType = strings.TrimSpace(mimeType)
	}

	// Extract filename from URL
	fileName := filepath.Base(url)
	if strings.Contains(fileName, "?") {
		fileName = strings.Split(fileName, "?")[0]
	}

	return data, mimeType, fileName, nil
}

// readFromFile reads content from a local file
func (mp *MediaProcessor) readFromFile(filePath string) ([]byte, string, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", "", err
	}

	// Detect MIME type from file extension
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))

	// Get filename
	fileName := filepath.Base(filePath)

	return data, mimeType, fileName, nil
}

// decodeBase64 decodes base64 content
func (mp *MediaProcessor) decodeBase64(input string) ([]byte, error) {
	// Remove data URL prefix if present (e.g., "data:image/jpeg;base64,")
	base64Data := input
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	return base64.StdEncoding.DecodeString(base64Data)
}

// DetectMimeTypeFromData detects MIME type from file data
func (mp *MediaProcessor) DetectMimeTypeFromData(data []byte) string {
	return http.DetectContentType(data)
}
