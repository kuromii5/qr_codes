package qrdecode

import (
	"strings"

	"github.com/kuromii5/qr_codes/models"
)

type Byte struct{}

var byteLen = [3]int{8, 16, 16}

func (d Byte) Decode(bits *BitsReader, version models.Version) (string, error) {
	var result strings.Builder
	charCount := int(bits.Read(byteLen[version.SizeClass()]))

	for i := 0; i < charCount; i++ {
		b := bits.Read(8)
		result.WriteByte(byte(b))
	}

	return result.String(), nil
}
