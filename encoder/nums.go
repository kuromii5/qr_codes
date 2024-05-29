package encoder

import (
	"fmt"

	"github.com/kuromii5/qr_codes/models"
)

// Numeric is the encoding for numeric data.
// The only valid characters are the decimal digits 0 through 9.
type Numeric string

// number of bits in character count indicator for QR code
var numLen = [3]int{10, 12, 14}

func (s Numeric) Check() error {
	for _, c := range s {
		if c < '0' || '9' < c {
			return fmt.Errorf("non-numeric string %#q", string(s))
		}
	}
	return nil
}

func (s Numeric) Bits(v models.Version) int {
	return 4 + numLen[v.SizeClass()] + (10*len(s)+2)/3
}

func (s Numeric) Encode(b *Bits, v models.Version) {
	b.Write(1, 4)                                // 0001
	b.Write(uint(len(s)), numLen[v.SizeClass()]) // write character count indicator
	var i int
	for i = 0; i+3 <= len(s); i += 3 {
		w := uint(s[i]-'0')*100 + uint(s[i+1]-'0')*10 + uint(s[i+2]-'0')
		b.Write(w, 10)
	}
	switch len(s) - i {
	case 1:
		w := uint(s[i] - '0')
		b.Write(w, 4)
	case 2:
		w := uint(s[i]-'0')*10 + uint(s[i+1]-'0')
		b.Write(w, 7)
	}
}
