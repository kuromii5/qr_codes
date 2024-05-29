package qrdecode

import (
	"github.com/kuromii5/qr_codes/models"
	"github.com/kuromii5/qr_codes/utils"
)

func getFormatInfo(matrix [][]byte) (models.Level, models.Mask) {
	formatInfo1 := uint32(0)
	formatInfo2 := uint32(0)

	// read horizontal
	for y := 0; y < len(matrix); y++ {
		if y == 8 {
			for i := 0; i < 6; i++ {
				formatInfo1 |= (uint32(matrix[y][i] & 1)) << uint(14-i)
			}
			formatInfo1 |= (uint32(matrix[y][7] & 1)) << uint(8)
			for i := 0; i < 8; i++ {
				formatInfo1 |= (uint32(matrix[y][len(matrix)-8+i] & 1)) << uint(7-i)
			}
		}
	}

	// read vertical
	for y := len(matrix) - 1; y >= 0; y-- {
		switch {
		case y > len(matrix)-8:
			formatInfo2 |= (uint32(matrix[y][8] & 1)) << uint(14-(len(matrix)-1-y))
		case y < 9 && y > 6:
			formatInfo2 |= (uint32(matrix[y][8] & 1)) << uint(14-y)
		case y < 6:
			formatInfo2 |= (uint32(matrix[y][8] & 1)) << uint(y)
		}
	}

	// invert bits
	invert := uint32(0x5412) // format info mask
	for i := uint(0); i < 15; i++ {
		if invert&(1<<i) != 0 {
			formatInfo1 ^= 1 << i // use one of the read bit strings
		}
	}

	// return the error correction level and mask
	level := models.ToLevel(formatInfo1 >> 13)
	mask := models.Mask((formatInfo1 >> 10) & 7)
	return level, mask
}

func getVersionInfo(matrix [][]byte) models.Version {
	versionInfo1 := []byte{}
	versionInfo2 := []byte{}

	for y := 0; y < len(matrix); y++ {
		// read top-right version bytes
		if y >= 0 && y < 6 {
			for x := len(matrix) - 11; x < len(matrix)-8; x++ {
				versionInfo1 = append(versionInfo1, matrix[y][x])
			}
		}

		// read bottom-left version bytes
		if y == len(matrix)-11 {
			for x := 0; x < 6; x++ {
				versionInfo2 = append(versionInfo2, matrix[y][x])
				versionInfo2 = append(versionInfo2, matrix[y+1][x])
				versionInfo2 = append(versionInfo2, matrix[y+2][x])
			}
		}
	}

	utils.ReverseBytes(versionInfo1)
	utils.ReverseBytes(versionInfo2)
	versionNumber1 := utils.BytesToUInt32(versionInfo1)
	versionNumber2 := utils.BytesToUInt32(versionInfo2)
	if versionNumber2 != versionNumber1 {
		panic("encoding isn't working well")
	}

	// use only one of them as usually
	version := utils.FindVersion(versionNumber1)
	return models.Version(version)
}
