package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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