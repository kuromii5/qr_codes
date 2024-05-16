package models

import "strconv"

type Field struct {
	log [256]byte
	exp [510]byte
}

// QR-Code Standard uses byte-wise modulo 100011101 == 285 arithmetic
var GField = NewField(0x11d, 2)

// NewField returns a new field corresponding to
// the given polynomial and generator.

// The Reed-Solomon encoding in QR codes uses
// polynomial 0x11d with generator 2.
func NewField(poly, generator int) *Field {
	if poly < 0x100 || poly >= 0x200 || Reducible(poly) {
		panic("gf256: invalid polynomial: " + strconv.Itoa(poly))
	}

	var f Field
	x := 1
	for i := 0; i < 255; i++ {
		if x == 1 && i != 0 {
			panic("gf256: invalid generator " + strconv.Itoa(generator) +
				" for polynomial " + strconv.Itoa(poly))
		}
		f.exp[i] = byte(x)
		f.exp[i+255] = byte(x)
		f.log[x] = byte(i)
		x = Mul(x, generator, poly)
	}
	f.log[0] = 255
	for i := 0; i < 255; i++ {
		if f.log[f.exp[i]] != byte(i) {
			panic("bad log")
		}
		if f.log[f.exp[i+255]] != byte(i) {
			panic("bad log")
		}
	}
	for i := 1; i < 256; i++ {
		if f.exp[f.log[i]] != byte(i) {
			panic("bad log")
		}
	}

	return &f
}

// Add returns the sum of x and y in the field.
func (f *Field) Add(x, y byte) byte {
	return x ^ y
}

// Exp returns the base 2 exponential of e in the field.
// If e < 0, Exp returns 0.
func (f *Field) Exp(e int) byte {
	if e < 0 {
		return 0
	}
	return f.exp[e%255]
}

// Log returns the base 2 logarithm of x in the field.
// If x == 0, Log returns -1.
func (f *Field) Log(x byte) int {
	if x == 0 {
		return -1
	}
	return int(f.log[x])
}

// Mul returns the product of x and y in the field.
func (f *Field) Mul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	return f.exp[int(f.log[x])+int(f.log[y])]
}

// Inv returns the multiplicative inverse of x in the field.
// If x == 0, Inv returns 0.
func (f *Field) Inv(x byte) byte {
	if x == 0 {
		return 0
	}
	return f.exp[255-f.log[x]]
}

// returns generator polynomial for given exponent
func (f *Field) getGeneratorPoly(exp int) (gen, logGen []byte) {
	poly := make([]byte, exp+1)
	poly[exp] = 1

	for i := 0; i < exp; i++ {
		currentExp := f.Exp(i)
		for j := 0; j < exp; j++ {
			poly[j] = f.Mul(poly[j], currentExp) ^ poly[j+1]
		}
		poly[exp] = f.Mul(poly[exp], currentExp)
	}

	// logPoly = log p.
	logPoly := make([]byte, exp+1)
	for i, coefficient := range poly {
		if coefficient == 0 {
			logPoly[i] = 255
		} else {
			logPoly[i] = byte(f.Log(coefficient))
		}
	}

	return poly, logPoly
}
