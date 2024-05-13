package encoder

import (
	"fmt"
	"strings"
)

// Alphanumeric is the encoding for alphanumeric data.
type Alphanumeric string

// The valid characters are 0-9A-Z$%*+-./: and space.
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

// number of bits in character count indicator for QR code
var alphaLen = [3]int{9, 11, 13}

func (s Alphanumeric) String() string {
	return fmt.Sprintf("Alpha(%#q)", string(s))
}

func (s Alphanumeric) Check() error {
	for _, c := range s {
		if !strings.ContainsRune(alphabet, c) {
			return fmt.Errorf("non-alphanumeric string %#q", string(s))
		}
	}
	return nil
}

func (s Alphanumeric) Bits(v Version) int {
	return 4 + alphaLen[v.SizeClass()] + (11*len(s)+1)/2
}

func (s Alphanumeric) Encode(b *Bits, v Version) {
	b.Write(2, 4)
	b.Write(uint(len(s)), alphaLen[v.SizeClass()])
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
