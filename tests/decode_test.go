package tests

import (
	"fmt"
	"testing"

	pngdecode "github.com/kuromii5/qr_codes/decoder/pngDecode"
	qrdecode "github.com/kuromii5/qr_codes/decoder/qrDecode"
	"github.com/kuromii5/qr_codes/encoder"
	"github.com/kuromii5/qr_codes/models"
	"github.com/stretchr/testify/assert"
)

func TestDecodePNG(t *testing.T) {
	qrMatrix, err := pngdecode.ReadQRCode("qr.png")
	if err != nil {
		fmt.Println("image reading error:", err)
		return
	}

	fmt.Println("extracted data:", qrMatrix.Content)
}

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		mode     encoder.Encoder
		level    models.Level
		expected string
	}{
		{
			name:     "Alphanumeric L",
			input:    "HELLO WORLD",
			mode:     encoder.Alphanumeric("HELLO WORLD"),
			level:    models.L,
			expected: "HELLO WORLD",
		},
		{
			name:     "Byte M",
			input:    "hello, world!",
			mode:     encoder.Byte("hello, world!"),
			level:    models.M,
			expected: "hello, world!",
		},
		{
			name:     "Numeric Q",
			input:    "1234567890",
			mode:     encoder.Numeric("1234567890"),
			level:    models.Q,
			expected: "1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qrCode, err := encoder.Encode(tt.input, tt.level)
			assert.NoError(t, err)

			decoded := qrdecode.ReadQR(qrCode)
			assert.Equal(t, tt.expected, decoded)
		})
	}
}
