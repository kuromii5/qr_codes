package models

const FormatPoly = 0x537 // poly G(x) = x^10+x^8+x^5+x^4+x^2+x+1 used for format info

const VersionPoly = 0x1f25 // poly G(x) = x^12+x^11+x^10+x^9+x^8+x^5+x^2+1 used for version info

// GF(256) mutiplication: mul returns the product x×y mod poly.
func Mul(x, y, poly int) int {
	z := 0
	for x > 0 {
		if x&1 != 0 {
			z ^= y
		}
		x >>= 1
		y <<= 1
		if y&0x100 != 0 {
			y ^= poly
		}
	}
	return z
}

// BitCount returns the number of significant in p.
func BitCount(p int) uint {
	n := uint(0)
	for ; p > 0; p >>= 1 {
		n++
	}
	return n
}

// polyDiv divides the polynomial p by q and returns the remainder.
func PolyDiv(p, q int) int {
	np := BitCount(p)
	nq := BitCount(q)
	for ; np >= nq; np-- {
		if p&(1<<(np-1)) != 0 {
			p ^= q << (np - nq)
		}
	}
	return p
}

func Reducible(p int) bool {
	// Multiplying bitCount * bitCount produces (2n-1)-bit,
	// so if p is reducible, one of its factors must be
	// of np/2+1 bits or fewer.
	np := BitCount(p)
	for q := 2; q < 1<<(np/2+1); q++ {
		if PolyDiv(p, q) == 0 {
			return true
		}
	}
	return false
}
