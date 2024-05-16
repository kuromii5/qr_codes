package utils

import "github.com/kuromii5/qr_codes/models"

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func FindVersion(vPattern uint32) models.Version {
	table := models.VTable
	for i := 1; i < len(table); i++ {
		if uint32(table[i].Pattern) == vPattern {
			return models.Version(i)
		}
	}

	return -1
}
