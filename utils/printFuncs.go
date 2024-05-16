package utils

import (
	"fmt"
	"strconv"

	"github.com/kuromii5/qr_codes/models"
)

// print formatData
func PrintFormatInfo(formatInfo uint32) {
	for i := 14; i >= 0; i-- {
		bit := (formatInfo >> i) & 1
		fmt.Print(bit)
	}
	fmt.Println()
}

func PrintVersionInfo(versionInfo uint32) {
	for i := 17; i >= 0; i-- {
		bit := (versionInfo >> i) & 1
		fmt.Print(bit)
	}
	fmt.Println()
}

func PrintGrid(grid [][]models.Pixel) {
	for _, row := range grid {
		for _, pixel := range row {
			if pixel&models.Black != 0 {
				fmt.Printf("1 ")
			} else {
				fmt.Printf("0 ")
			}
		}
		fmt.Println()
	}

	fmt.Println("-------------------------------------------")
}

func PrintQRBitmap(qr *models.QRCode) {
	matrix := qr.CreateBitmapMatrix()
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			if matrix[i][j] == 1 {
				fmt.Print("■ ") // Черный пиксель
			} else {
				fmt.Print("□ ") // Белый пиксель
			}
		}
		fmt.Println()
	}
}

func PrintLevel(l models.Level) string {
	if models.L <= l && l <= models.H {
		return "LMQH"[l : l+1]
	}
	return strconv.Itoa(int(l))
}
