package qrdecode

import (
	"strings"

	"github.com/kuromii5/qr_codes/models"
)

type Alphanumeric struct{}

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

var alphaLen = [3]int{9, 11, 13}

// Decode takes the QR code data string and version, and decodes it to a readable string.
func (d Alphanumeric) Decode(bits *BitsReader, v models.Version) (string, error) {
	var result strings.Builder

	// Read character count indicator
	charCount := bits.Read(alphaLen[v.SizeClass()])

	for charCount > 1 {
		w := bits.Read(11)
		result.WriteByte(alphabet[w/45])
		result.WriteByte(alphabet[w%45])
		charCount -= 2
	}

	if charCount == 1 {
		w := bits.Read(6)
		result.WriteByte(alphabet[w])
	}

	return result.String(), nil
}
