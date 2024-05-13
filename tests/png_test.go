package tests

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/kuromii5/qr_codes/encoder"
)

func TestPNG(t *testing.T) {
	c, err := encoder.Encode("https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley", encoder.L)
	if err != nil {
		t.Fatal(err)
	}

	pngdat := c.PNG()
	os.WriteFile("qr.png", pngdat, 0666)

	m, err := png.Decode(bytes.NewBuffer(pngdat))
	if err != nil {
		t.Fatal(err)
	}

	gm := m.(*image.Gray)
	scale := c.Scale
	siz := c.Size
	nbad := 0
	for y := 0; y < scale*(8+siz); y++ {
		for x := 0; x < scale*(8+siz); x++ {
			v := byte(255)
			if c.Black(x/scale-4, y/scale-4) {
				v = 0
			}
			if gv := gm.At(x, y).(color.Gray).Y; gv != v {
				t.Errorf("%d,%d = %d, want %d", x, y, gv, v)
				if nbad++; nbad >= 20 {
					t.Fatalf("too many bad pixels")
				}
			}
		}
	}
}
