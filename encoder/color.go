package encoder

import (
	"image"
	"image/color"
)

var (
	whiteColor color.Color = color.Gray{0xFF}
	blackColor color.Color = color.Gray{0x00}
)

// codeImage implements image.Image
type codeImage struct {
	*QRCode
}

// Image returns an Image displaying the code.
func (c *QRCode) Image() image.Image {
	return &codeImage{c}
}

func (c *codeImage) Bounds() image.Rectangle {
	d := (c.Size + 8) * c.Scale
	return image.Rect(0, 0, d, d)
}

func (c *codeImage) At(x, y int) color.Color {
	if c.Black(x, y) {
		return blackColor
	}
	return whiteColor
}

func (c *codeImage) ColorModel() color.Model {
	return color.GrayModel
}
