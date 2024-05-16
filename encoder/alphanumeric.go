package encoder

import (
	"fmt"
	"strings"

	"github.com/kuromii5/qr_codes/models"
)

// Alphanumeric is the encoding for alphanumeric data.
type Alphanumeric string

// The valid characters are 0-9A-Z$%*+-./: and space.
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

// number of bits in character count indicator for QR code
var alphaLen = [3]int{9, 11, 13}

func (s Alphanumeric) Check() error {
	for _, c := range s {
		if !strings.ContainsRune(alphabet, c) {
			return fmt.Errorf("non-alphanumeric string %#q", string(s))
		}
	}
	return nil
}

func (s Alphanumeric) Bits(v models.Version) int {
	return 4 + alphaLen[v.SizeClass()] + (11*len(s)+1)/2
}

func (s Alphanumeric) Encode(b *Bits, v models.Version) {
	b.Write(2, 4)                                  // 0010
	b.Write(uint(len(s)), alphaLen[v.SizeClass()]) // write character count indicator
	var i int
	for i = 0; i+2 <= len(s); i += 2 {
		w := uint(strings.IndexRune(alphabet, rune(s[i])))*45 +
			uint(strings.IndexRune(alphabet, rune(s[i+1])))
		b.Write(w, 11)
	}

	if i < len(s) {
		w := uint(strings.IndexRune(alphabet, rune(s[i])))
		b.Write(w, 6)
	}
}
