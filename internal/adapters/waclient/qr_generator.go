package waclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
	"zpwoot/internal/core/session"
	"zpwoot/platform/logger"
)


type QRGenerator struct {
	logger *logger.Logger
}


func NewQRGenerator(logger *logger.Logger) session.QRCodeGenerator {
	return &QRGenerator{
		logger: logger,
	}
}


func (qr *QRGenerator) Generate(ctx context.Context, sessionName string) (*session.QRCodeResponse, error) {
	qr.logger.InfoWithFields("Generating QR code", map[string]interface{}{
		"session_name": sessionName,
	})



	expiresAt := time.Now().Add(2 * time.Minute)
	
	return &session.QRCodeResponse{
		QRCode:    "",
		ExpiresAt: expiresAt,
		Timeout:   120,
	}, nil
}


func (qr *QRGenerator) GenerateImage(ctx context.Context, qrCode string) ([]byte, error) {
	if qrCode == "" {
		return nil, fmt.Errorf("QR code string cannot be empty")
	}

	qr.logger.DebugWithFields("Generating QR code image", map[string]interface{}{
		"qr_code_length": len(qrCode),
	})


	image, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		qr.logger.ErrorWithFields("Failed to encode QR code", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to encode QR code: %w", err)
	}

	qr.logger.DebugWithFields("QR code image generated successfully", map[string]interface{}{
		"image_size": len(image),
	})

	return image, nil
}


func (qr *QRGenerator) GenerateBase64Image(ctx context.Context, qrCode string) (string, error) {
	image, err := qr.GenerateImage(ctx, qrCode)
	if err != nil {
		return "", err
	}


	base64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)

	qr.logger.DebugWithFields("Base64 QR code image generated", map[string]interface{}{
		"base64_length": len(base64Image),
	})

	return base64Image, nil
}


func (qr *QRGenerator) IsExpired(expiresAt time.Time) bool {
	expired := time.Now().After(expiresAt)
	
	qr.logger.DebugWithFields("Checking QR code expiration", map[string]interface{}{
		"expires_at": expiresAt,
		"now":        time.Now(),
		"expired":    expired,
	})

	return expired
}


func (qr *QRGenerator) GetExpirationTime() time.Duration {
	return 2 * time.Minute
}


func (qr *QRGenerator) ValidateQRCode(qrCode string) error {
	if qrCode == "" {
		return fmt.Errorf("QR code cannot be empty")
	}


	if len(qrCode) < 10 {
		return fmt.Errorf("QR code too short")
	}

	if len(qrCode) > 2048 {
		return fmt.Errorf("QR code too long")
	}

	qr.logger.DebugWithFields("QR code validated", map[string]interface{}{
		"qr_code_length": len(qrCode),
	})

	return nil
}


func (qr *QRGenerator) CreateQRResponse(qrCode string, base64Image string) *session.QRCodeResponse {
	expiresAt := time.Now().Add(qr.GetExpirationTime())
	
	response := &session.QRCodeResponse{
		QRCode:      qrCode,
		QRCodeImage: base64Image,
		ExpiresAt:   expiresAt,
		Timeout:     int(qr.GetExpirationTime().Seconds()),
	}

	qr.logger.DebugWithFields("QR response created", map[string]interface{}{
		"expires_at": expiresAt,
		"timeout":    response.Timeout,
	})

	return response
}


func (qr *QRGenerator) ProcessQRCodeFromWhatsApp(ctx context.Context, qrCode string) (*session.QRCodeResponse, error) {
	qr.logger.InfoWithFields("Processing QR code from WhatsApp", map[string]interface{}{
		"qr_code_length": len(qrCode),
	})


	err := qr.ValidateQRCode(qrCode)
	if err != nil {
		qr.logger.ErrorWithFields("Invalid QR code received", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}


	base64Image, err := qr.GenerateBase64Image(ctx, qrCode)
	if err != nil {
		qr.logger.ErrorWithFields("Failed to generate base64 image", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}


	response := qr.CreateQRResponse(qrCode, base64Image)

	qr.logger.InfoWithFields("QR code processed successfully", map[string]interface{}{
		"expires_at": response.ExpiresAt,
		"timeout":    response.Timeout,
	})

	return response, nil
}


func (qr *QRGenerator) CleanupExpiredQRCodes(ctx context.Context, db interface{}) error {


	qr.logger.DebugWithFields("Cleaning up expired QR codes", map[string]interface{}{
		"timestamp": time.Now(),
	})


	return nil
}


func (qr *QRGenerator) GetQRCodeStats() map[string]interface{} {
	return map[string]interface{}{
		"default_expiration_minutes": int(qr.GetExpirationTime().Minutes()),
		"default_timeout_seconds":    int(qr.GetExpirationTime().Seconds()),
		"image_format":              "PNG",
		"image_size":                "256x256",
		"encoding":                  "base64",
	}
}
