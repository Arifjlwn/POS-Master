package utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// Menentukan batas maksimal pesan WA yang diproses di background secara bersamaan (misal: 10)
var waJobQueue = make(chan func(), 100)

func init() {
	// Menjalankan 10 worker permanen untuk memproses antrean WhatsApp agar aman dari RAM Leak
	for i := 0; i < 10; i++ {
		go func() {
			for job := range waJobQueue {
				job()
			}
		}()
	}
}

// SendWhatsAppFonnte mengirim notifikasi WA atas nama Token milik Toko Tenant (Owner)
func SendWhatsAppFonnte(token, targetPhone, message string) {
	if token == "" || targetPhone == "" {
		return
	}

	// Memasukkan fungsi ke dalam antrean aman (Worker Pool) menggantikan `go func()` liar
	waJobQueue <- func() {
		payload := new(bytes.Buffer)
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("target", targetPhone)
		_ = writer.WriteField("message", message)
		_ = writer.Close()

		req, err := http.NewRequest("POST", "https://api.fonnte.com/send", payload)
		if err != nil {
			log.Println("Gagal membuat request WA:", err)
			return
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Gagal mengirim WA:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		log.Println("Respon Fonnte:", string(body))
	}
}

// SendSystemWhatsApp mengirim notifikasi sistem (OTP/Billing) menggunakan Token Utama Developer Aplikasi
func SendSystemWhatsApp(targetPhone, message string) {
	systemToken := os.Getenv("WA_SYSTEM_TOKEN")
	if systemToken == "" {
		log.Println("⚠️ Peringatan: WA_SYSTEM_TOKEN di .env belum diisi!")
		return
	}
	SendWhatsAppFonnte(systemToken, targetPhone, message)
}