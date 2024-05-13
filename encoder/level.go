package encoder

import "strconv"

// A Level denotes a QR error correction level.
// From least to most tolerant of errors, they are L, M, Q, H.
type Level int

const (
	L Level = iota // 20% redundant 7%  recovery of the symbol codewords
	M              // 38% redundant 15% recovery of the symbol codewords
	Q              // 55% redundant 25% recovery of the symbol codewords
	H              // 65% redundant 30% recovery of the symbol codewords
)

func (l Level) String() string {
	if L <= l && l <= H {
		return "LMQH"[l : l+1]
	}
	return strconv.Itoa(int(l))
}
