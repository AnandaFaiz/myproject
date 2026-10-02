package utils

import (
	"regexp"
	"strings"
)

func BuatSlug(teks string) string {
	//lowercase
	teks = strings.ToLower(teks)

	//trimspace
	teks = strings.TrimSpace(teks)

	//hapus char non-alfanumerik
	reg := regexp.MustCompile(`[^\w\s-]`)
	teks = reg.ReplaceAllString(teks, " ")

	//ganti spasi => tanda hubung
	teks = strings.ReplaceAll(teks, " ", "-")

	//rapikan tanda hubung berulang
	reg = regexp.MustCompile(`-+`)
	teks = reg.ReplaceAllString(teks, "-")

	return teks
}
