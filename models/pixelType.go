package models

// A PixelType describes the type of a QR pixel.
type PixelType uint32

const (
	_          PixelType = iota
	Finder               // finder patterns (large)
	Alignment            // alignment squares (small)
	Timing               // timing strip between position squares
	Format               // format metadata
	PVersion             // version pattern
	DarkModule           // unused dark module pixel
	Data                 // data bit
	EC                   // error correction check bit
	Extra                // extra bits placed after error correction bits
)

func (r PixelType) Pixel() Pixel {
	return Pixel(r << 2)
}
