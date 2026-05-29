package utils

import "strings"

// FormatPhoneNumber membersihkan format nomor HP menjadi standar 628xxx
// Contoh: "0812 3456" -> "628123456", "+62-812" -> "62812"
func FormatPhoneNumber(phone string) string {
	// 1. Hapus semua spasi, strip, dan tanda plus
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// 2. Kalau depannya "0", potong "0"-nya dan ganti jadi "62"
	if strings.HasPrefix(phone, "0") {
		return "62" + phone[1:]
	}

	// 3. Kalau udah "62" atau format lain, biarkan apa adanya
	return phone
}