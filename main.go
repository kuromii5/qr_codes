package main

import (
	"fmt"

	"github.com/kuromii5/qr_codes/decoder"
	"github.com/kuromii5/qr_codes/encoder"
	"github.com/kuromii5/qr_codes/models"
)

func main() {
	qrCode, err := encoder.Encode("HELLO WORLD", models.Q)
	if err != nil {
		fmt.Println(err)
	}

	decoded := decoder.Decode(qrCode)
	fmt.Println("The encoded data:", decoded)
}

func TestDecodePNG(filename string) {
	qrMatrix, err := decoder.ReadQRCode(filename)
	if err != nil {
		fmt.Println("image reading error:", err)
		return
	}

	fmt.Println(qrMatrix.Content)
}
