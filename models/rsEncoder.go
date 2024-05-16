package models

type RSEncoder struct {
	f             *Field // field used for operations
	ECBytesLength int    // error correction codewords length
	generator     []byte // generator poly
	logGenerator  []byte // logarithmic representation of gen poly
	poly          []byte // poly for error correction codewords
}

// NewRSEncoder returns a new Reed-Solomon encoder
// over the given field and number of error correction bytes.
func NewRSEncoder(f *Field, ECBytes int) *RSEncoder {
	gen, logGen := f.getGeneratorPoly(ECBytes)
	return &RSEncoder{
		f:             f,
		ECBytesLength: ECBytes,
		generator:     gen,
		logGenerator:  logGen,
	}
}

// AddECC writes to check the error correcting code bytes
// for data using the given Reed-Solomon parameters.
func (encoder *RSEncoder) AddECC(data []byte, ECBytes []byte) {
	if len(ECBytes) < encoder.ECBytesLength {
		panic("invalid check byte length")
	}

	// The error correction bytes are the remainder after dividing
	// data poly padded with ECC length zeros by the generator polynomial.

	// poly = data padded with ECC length zeros.
	var poly []byte
	n := len(data) + encoder.ECBytesLength
	if len(encoder.poly) >= n {
		poly = encoder.poly
	} else {
		poly = make([]byte, n)
	}
	copy(poly, data)
	for i := len(data); i < len(poly); i++ {
		poly[i] = 0
	}

	// Divide poly by gen, leaving the remainder in poly[len(data):].
	// poly[0] is the most significant term in poly, and
	// gen[0] is the most significant term in the generator,
	// which is always 1.
	f := encoder.f
	logGen := encoder.logGenerator[1:]
	for i := 0; i < len(data); i++ {
		coefficient := poly[i]
		if coefficient == 0 {
			continue
		}
		q := poly[i+1:]
		exp := f.exp[f.log[coefficient]:]
		for j, logGen := range logGen {
			if logGen != 255 { // logGen uses 255 for log 0
				q[j] ^= exp[logGen]
			}
		}
	}
	copy(ECBytes, poly[len(data):])
	encoder.poly = poly
}
