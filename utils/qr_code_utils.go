package utils

import (
	"github.com/skip2/go-qrcode"
)

// GenerateQRCodePNG returns raw PNG bytes instead of base64 string
func GenerateQRCodePNG(data string) ([]byte, error) {
	return qrcode.Encode(data, qrcode.Medium, 256)
}
