package models

// A Level denotes a QR error correction level.
// From least to most tolerant of errors, they are L, M, Q, H.
type Level int

const (
	L Level = iota // 20% redundant 7%  recovery of the symbol codewords
	M              // 38% redundant 15% recovery of the symbol codewords
	Q              // 55% redundant 25% recovery of the symbol codewords
	H              // 65% redundant 30% recovery of the symbol codewords
)

func ToLevel(bits uint32) Level {
	switch bits {
	case 0:
		return M
	case 1:
		return L
	case 2:
		return H
	case 3:
		return Q
	}

	return L
}
