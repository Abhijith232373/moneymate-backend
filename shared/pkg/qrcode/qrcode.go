package qrcode

import (
	"encoding/base64"
	"fmt"
	
	goqrcode "github.com/skip2/go-qrcode"
)

// GenerateBase64 generates a QR code for the given data and returns it as a base64-encoded PNG data URI.
func GenerateBase64(data string) (string, error) {
	// Generate a 256x256 QR code.
	pngBytes, err := goqrcode.Encode(data, goqrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	
	// Convert to base64 data URI format so it can be directly used in <img src="..." />
	base64Encoded := base64.StdEncoding.EncodeToString(pngBytes)
	return fmt.Sprintf("data:image/png;base64,%s", base64Encoded), nil
}
