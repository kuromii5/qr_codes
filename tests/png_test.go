package tests

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/kuromii5/qr_codes/encoder"
	"github.com/kuromii5/qr_codes/models"
	"github.com/kuromii5/qr_codes/utils"
)

func TestEncodePNG(t *testing.T) {
	qrCode, err := encoder.Encode("HELLO WORLD", models.L)
	if err != nil {
		t.Fatal(err)
	}

	pngdat := utils.PNG(qrCode)
	os.WriteFile("qr.png", pngdat, 0666)

	m, err := png.Decode(bytes.NewBuffer(pngdat))
	if err != nil {
		t.Fatal(err)
	}

	gm := m.(*image.Gray)
	scale := qrCode.Scale
	size := qrCode.Size
	nbad := 0
	for y := 0; y < scale*(8+size); y++ {
		for x := 0; x < scale*(8+size); x++ {
			v := byte(255)
			if qrCode.Black(x/scale-4, y/scale-4) {
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
