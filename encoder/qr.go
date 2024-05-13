package encoder

import (
	"errors"
)

// A QRCode is a square pixel grid.
type QRCode struct {
	Bitmap []byte // 1 is black, 0 is white
	Size   int    // number of pixels on a side
	Stride int    // number of bytes per row
	Scale  int    // number of image pixels per QR pixel
}

var GField = NewField(0x11d, 2)

// Black returns true if the pixel at (x,y) is black.
func (c *QRCode) Black(x, y int) bool {
	return 0 <= x && x < c.Size && 0 <= y && y < c.Size &&
		c.Bitmap[y*c.Stride+x/8]&(1<<uint(7-x&7)) != 0
}

// Encode returns an encoding of text at the given error correction level.
func Encode(text string, level Level) (*QRCode, error) {
	// Pick data encoding, smallest first.
	var enc Encoding
	switch {
	case Num(text).Check() == nil:
		enc = Num(text)
	case Alphanumeric(text).Check() == nil:
		enc = Alphanumeric(text)
	default:
		enc = String(text)
	}

	// Pick size.
	l := Level(level)
	var v Version
	for v = MinVersion; ; v++ {
		if v > MaxVersion {
			return nil, errors.New("text too long to encode as QR")
		}
		if enc.Bits(v) <= v.DataBytes(l)*8 {
			break
		}
	}

	// Build and execute template.
	template := NewTemplate(v, l, 0)

	qrCode, err := template.Encode(enc)
	if err != nil {
		return nil, err
	}

	// TODO: Pick appropriate mask.

	return &QRCode{qrCode.Bitmap, qrCode.Size, qrCode.Stride, 8}, nil
}

// PNG returns a PNG image displaying the code.
//
// PNG uses a custom encoder tailored to QR codes.
// Its compressed size is about 2x away from optimal
func (c *QRCode) PNG() []byte {
	var p PNGWriter
	return p.Encode(c)
}
