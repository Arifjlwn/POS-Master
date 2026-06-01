package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// FonntePayload menyesuaikan format API Fonnte
type FonntePayload struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

func SendWhatsAppFonnte(token, targetPhone, message string) {
	// Gunakan Goroutine agar tidak memblokir respon API Kasir (biar loading kasir tetap cepat)
	go func() {
		if token == "" || targetPhone == "" {
			return // Skip kalau toko gak punya token WA atau pelanggan gak ngasih nomor
		}

		payload := FonntePayload{
			Target:  targetPhone,
			Message: message,
		}
		jsonPayload, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "https://api.fonnte.com/send", bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Println("Gagal membuat request WA:", err)
			return
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Gagal mengirim WA:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Respon Fonnte:", string(body))
	}()
}

// 🚀 FUNGSI BARU: Khusus kirim OTP / Notifikasi Sistem pakai Token Sistem milik LU dari .env
func SendSystemWhatsApp(targetPhone, message string) {
	// Ambil token sistem rahasia lu langsung dari .env
	// Pastiin nama key di .env lu sama persis (WA_SYSTEM_TOKEN)
	systemToken := os.Getenv("WA_SYSTEM_TOKEN")
	
	if systemToken == "" {
		log.Println("⚠️ Peringatan: WA_SYSTEM_TOKEN di .env belum diisi! Gagal mengirim pesan sistem.")
		return
	}

	// Panggil fungsi utama di atas, tapi kita suntik pakai token sistem milik lu
	SendWhatsAppFonnte(systemToken, targetPhone, message)
}