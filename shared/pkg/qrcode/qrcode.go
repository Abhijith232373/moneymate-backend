// package qrcode

// import (
// 	"encoding/base64"
// 	"fmt"
	
// 	goqrcode "github.com/skip2/go-qrcode"
// )

// // GenerateBase64 generates a QR code for the given data and returns it as a base64-encoded PNG data URI.
// func GenerateBase64(data string) (string, error) {
// 	// Generate a 256x256 QR code.
// 	pngBytes, err := goqrcode.Encode(data, goqrcode.Medium, 256)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate QR code: %w", err)
// 	}
	
// 	// Convert to base64 data URI format so it can be directly used in <img src="..." />
// 	base64Encoded := base64.StdEncoding.EncodeToString(pngBytes)
// 	return fmt.Sprintf("data:image/png;base64,%s", base64Encoded), nil
// }


package qrcode

import (
	"encoding/base64"
	"fmt"
	"net/url"

	goqrcode "github.com/skip2/go-qrcode"
)

const paymentQRScheme = "moneymate://pay"

// BuildPaymentPayload builds the URI encoded into a payment QR code.
func BuildPaymentPayload(accountType, handle string) string {
	v := url.Values{}
	v.Set("type", accountType) // "user" or "merchant"
	v.Set("handle", handle)
	v.Set("v", "1")  
	return fmt.Sprintf("%s?%s", paymentQRScheme, v.Encode())
}

// GenerateBase64 generates a QR code for the given data and returns it as a base64-encoded PNG data URI.
func GenerateBase64(data string) (string, error) {
	pngBytes, err := goqrcode.Encode(data, goqrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}
	base64Encoded := base64.StdEncoding.EncodeToString(pngBytes)
	return fmt.Sprintf("data:image/png;base64,%s", base64Encoded), nil
}