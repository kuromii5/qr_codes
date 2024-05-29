package qrdecode

import (
	"fmt"
	"strings"

	"github.com/kuromii5/qr_codes/models"
)

type Numeric struct{}

var numLen = [3]int{10, 12, 14}

func (d Numeric) Decode(bits *BitsReader, version models.Version) (string, error) {
	var result strings.Builder
	charCount := bits.Read(numLen[version.SizeClass()])

	for charCount >= 3 {
		w := bits.Read(10)
		result.WriteString(fmt.Sprintf("%03d", w))
		charCount -= 3
	}

	if charCount == 2 {
		w := bits.Read(7)
		result.WriteString(fmt.Sprintf("%02d", w))
	} else if charCount == 1 {
		w := bits.Read(4)
		result.WriteString(fmt.Sprintf("%1d", w))
	}

	return result.String(), nil
}
