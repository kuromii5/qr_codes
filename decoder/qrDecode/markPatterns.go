package qrdecode

import (
	"math"

	"github.com/kuromii5/qr_codes/models"
)

func MarkPatterns(v models.Version, matrix [][]byte) [][]byte {
	size := v.GetPixelSize()

	// Mark Finder Patterns
	markFinderPattern(matrix, 0, 0)
	markFinderPattern(matrix, 0, size-8)
	markFinderPattern(matrix, size-8, 0)

	if v > 6 {
		markVersionPattern(matrix)
	}

	// Mark Timing Patterns
	for i := 8; i < size-8; i++ {
		matrix[6][i] = 2 // value to mark patterns
		matrix[i][6] = 2
	}

	// Mark Alignment Patterns for versions 2 and above
	if v >= 2 {
		alignmentPatternCenters := getAlignmentPatternCenters(v)
		for _, centerX := range alignmentPatternCenters {
			for _, centerY := range alignmentPatternCenters {
				// Avoid overlapping with finder patterns
				if (centerX == 6 && (centerY == 6 || centerY == size-7)) || (centerY == 6 && centerX == size-7) {
					continue
				}
				markAlignmentPattern(matrix, centerX, centerY)
			}
		}
	}

	// Mark Format patterns
	markFormatPatterns(matrix)

	matrix[size-8][8] = 2 // Dark Module pixel

	return matrix
}

func markFormatPatterns(matrix [][]byte) {
	// horizontal pattern
	for y := 0; y < len(matrix); y++ {
		if y == 8 {
			for x := 0; x < 9; x++ {
				matrix[y][x] = 2
			}
			for x := len(matrix) - 1; x > len(matrix)-9; x-- {
				matrix[y][x] = 2
			}
		}
	}

	// vertical pattern
	for y := 0; y < 9; y++ {
		matrix[y][8] = 2
	}
	for y := len(matrix) - 1; y > len(matrix)-9; y-- {
		matrix[y][8] = 2
	}
}

func markVersionPattern(matrix [][]byte) {
	for y := 0; y < len(matrix); y++ {
		// top-right bits
		if y >= 0 && y < 6 {
			for x := len(matrix) - 11; x < len(matrix)-8; x++ {
				matrix[y][x] = 2
			}
		}

		// bottom-left bits
		if y == len(matrix)-11 {
			for x := 0; x < 6; x++ {
				matrix[y][x] = 2
				matrix[y+1][x] = 2
				matrix[y+2][x] = 2
			}
		}
	}
}

func markFinderPattern(matrix [][]byte, x, y int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			matrix[x+i][y+j] = 2 // value to mark patterns
		}
	}
}

func getAlignmentPatternCenters(v models.Version) []int {
	if v == 1 {
		return nil
	}
	numCenters := int(math.Floor(float64(v)/7.0) + 2)
	centers := make([]int, numCenters)
	centers[0] = 6
	centers[numCenters-1] = (4 * int(v)) + 10
	step := (centers[numCenters-1] - 6) / (numCenters - 1)
	for i := 1; i < numCenters-1; i++ {
		centers[i] = centers[i-1] + step
	}
	return centers
}

func markAlignmentPattern(matrix [][]byte, x, y int) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			matrix[x+i][y+j] = 2 // value to mark patterns
		}
	}
}
