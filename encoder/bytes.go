package encoder

import (
	"github.com/kuromii5/qr_codes/models"
)

// Byte is the encoding for 8-bit data. All bytes are valid.
type Byte string

// number of bits in character count indicator for QR code
var stringLen = [3]int{8, 16, 16}

func (s Byte) Check() error {
	return nil
}

func (s Byte) Bits(v models.Version) int {
	return 4 + stringLen[v.SizeClass()] + 8*len(s)
}

func (s Byte) Encode(b *Bits, v models.Version) {
	b.Write(4, 4)                                   // 0100
	b.Write(uint(len(s)), stringLen[v.SizeClass()]) // write character count indicator
	for i := 0; i < len(s); i++ {
		b.Write(uint(s[i]), 8)
	}
}
