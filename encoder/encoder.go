package encoder

import (
	"errors"

	"github.com/kuromii5/qr_codes/models"
)

// Encoder implements a QR data encoding scheme.
// Modes - Numeric, Alphanumeric, and String (Byte mode)
// the character set and the mapping from UTF-8 to code bits.
// The more restrictive the mode, the fewer code bits are needed.
type Encoder interface {
	Check() error
	Bits(v models.Version) int
	Encode(b *Bits, v models.Version)
}

// Encode returns an encoding of text at the given error correction level.
func Encode(text string, level models.Level) (*models.QRCode, error) {
	// Pick data encoding, smallest first.
	var enc Encoder
	switch {
	case Numeric(text).Check() == nil:
		enc = Numeric(text)
	case Alphanumeric(text).Check() == nil:
		enc = Alphanumeric(text)
	default:
		enc = Byte(text)
	}

	// Pick size.
	l := models.Level(level)
	var v models.Version
	for v = models.MinVersion; ; v++ {
		if v > models.MaxVersion {
			return nil, errors.New("text too long to encode as QR")
		}
		if enc.Bits(v) <= v.DataBytes(l)*8 {
			break
		}
	}

	// Build and execute template.
	template := NewTemplate(v, l, 0)

	// Build actual QR code and encode the given data
	qrCode, err := template.Encode(enc)
	if err != nil {
		return nil, err
	}

	return &models.QRCode{
		Bitmap: qrCode.Bitmap,
		Size:   qrCode.Size,
		Stride: qrCode.Stride,
		Scale:  8,
	}, nil
}
