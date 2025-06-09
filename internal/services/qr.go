package services

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

type QRService struct{}

func NewQRService() *QRService {
	return &QRService{}
}

// GenerateQRCode generates a QR code for the given URL and returns it as base64 data URL
func (s *QRService) GenerateQRCode(url string, size int) (string, error) {
	// Generate QR code as PNG bytes
	qrBytes, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert to base64 data URL
	base64Str := base64.StdEncoding.EncodeToString(qrBytes)
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Str)

	return dataURL, nil
}

// GenerateQRCodeBytes generates a QR code and returns raw PNG bytes for downloads
func (s *QRService) GenerateQRCodeBytes(url string, size int) ([]byte, error) {
	qrBytes, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrBytes, nil
}