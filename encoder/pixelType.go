package encoder

import "strconv"

// A PixelType describes the type of a QR pixel.
type PixelType uint32

const (
	_         PixelType = iota
	Finder              // finder patterns (large)
	Alignment           // alignment squares (small)
	Timing              // timing strip between position squares
	Format              // format metadata
	PVersion            // version pattern
	Unused              // unused pixel
	Data                // data bit
	EC                  // error correction check bit
	Extra
)

var roles = []string{
	"",
	"finder pattern",
	"alignment",
	"timing",
	"format",
	"pversion",
	"unused",
	"data",
	"check",
	"extra",
}

func (r PixelType) Pixel() Pixel {
	return Pixel(r << 2)
}

func (r PixelType) String() string {
	if Finder <= r && r <= EC {
		return roles[r]
	}
	return strconv.Itoa(int(r))
}
