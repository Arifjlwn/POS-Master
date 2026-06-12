package utils

import (
	"fmt"
	"strings"
)

// ==========================================
// 💸 FORMAT RUPIAH HELPER (ANTI-COMPILER ERROR)
// ==========================================
// FormatRupiah mengubah angka int64 menjadi string berformat titik ribuan indonesia 
func FormatRupiah(amount int64) string {
	// Kalau minus, kita amankan tandanya 
	isNegative := false
	if amount < 0 {
		isNegative = true
		amount = -amount
	}

	out := fmt.Sprintf("%d", amount)
	if len(out) <= 3 {
		if isNegative {
			return "-" + out
		}
		return out
	}

	var result []string
	// Loop mundur tiap 3 digit kita selipin titik kasta tertinggi
	for len(out) > 3 {
		result = append([]string{out[len(out)-3:]}, result...)
		out = out[:len(out)-3]
	}
	
	if len(out) > 0 {
		result = append([]string{out}, result...)
	}

	finalStr := strings.Join(result, ".")
	if isNegative {
		return "-" + finalStr
	}
	return finalStr
}