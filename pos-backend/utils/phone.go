package utils

import "strings"

func FormatPhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", ""); phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+", ""); phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", ""); phone = strings.ReplaceAll(phone, ".", "")
	if strings.HasPrefix(phone, "0") { return "62" + phone[1:] }
	if strings.HasPrefix(phone, "8") { return "62" + phone }
	return phone
}