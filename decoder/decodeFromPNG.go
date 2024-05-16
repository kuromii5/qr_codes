package decoder

import (
	"bufio"
	"os"

	"github.com/tuotoo/qrcode"
)

func ReadQRCode(filename string) (*qrcode.Matrix, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	qr, err := qrcode.Decode(reader)
	if err != nil {
		return nil, err
	}

	return qr, nil
}
