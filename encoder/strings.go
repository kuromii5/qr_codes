package encoder

import (
	"fmt"
)

// String is the encoding for 8-bit data.  All bytes are valid.
type String string

// number of bits in character count indicator for QR code
var stringLen = [3]int{8, 16, 16}

func (s String) String() string {
	return fmt.Sprintf("String(%#q)", string(s))
}

func (s String) Check() error {
	return nil
}

func (s String) Bits(v Version) int {
	return 4 + stringLen[v.SizeClass()] + 8*len(s)
}

func (s String) Encode(b *Bits, v Version) {
	b.Write(4, 4)
	b.Write(uint(len(s)), stringLen[v.SizeClass()])
	for i := 0; i < len(s); i++ {
		b.Write(uint(s[i]), 8)
	}
}
