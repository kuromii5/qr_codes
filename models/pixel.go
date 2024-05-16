package models

// A Pixel describes a single pixel in a QR code.
type Pixel uint32

const (
	Black Pixel = 1 << iota
	White
)

func (p Pixel) Offset() uint {
	return uint(p >> 6)
}

func (p Pixel) Type() PixelType {
	return PixelType(p>>2) & 15
}

func OffsetPixel(o uint) Pixel {
	return Pixel(o << 6)
}

func Grid(size int) [][]Pixel {
	m := make([][]Pixel, size)
	pix := make([]Pixel, size*size)
	for i := range m {
		m[i], pix = pix[:size], pix[size:]
	}
	return m
}
