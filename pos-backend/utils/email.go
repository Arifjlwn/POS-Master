package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

func SendOTPEmail(targetEmail string, otp string) error {
	// Ambil dari .env/Railway. Jika kosong, return error (Lebih Aman daripada di-hardcode)
	from := os.Getenv("SMTP_FROM")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || username == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("konfigurasi SMTP di Environment Variable belum lengkap")
	}

	// Format Email Header & Body (Disesuaikan jadi ARZURA POS)
	subject := "Subject: [ARZURA POS] Kode Verifikasi Akun\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<div style="font-family: sans-serif; max-width: 400px; margin: auto; border: 1px solid #eee; padding: 20px; border-radius: 20px;"><h2 style="color: #000000; text-align: center;">ARZURA<span style="color: #64748b;">POS</span></h2><p style="text-align: center; color: #64748b;">Gunakan kode OTP di bawah untuk verifikasi:</p><div style="font-size: 40px; font-weight: 900; letter-spacing: 10px; text-align: center; color: #000000; padding: 20px; background: #f8fafc; border-radius: 15px; margin: 20px 0;">%s</div><p style="font-size: 10px; color: #94a3b8; text-align: center; text-transform: uppercase; letter-spacing: 1px;">Kode berlaku 5 menit</p></div>`, otp)

	msg := []byte(subject + mime + body)

	// Setup Autentikasi SMTP
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// JALUR AMAN: Menggunakan TLS Dial untuk jaminan enkripsi koneksi (Wajib untuk Gmail Port 587/465)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // WAJIB false agar server memvalidasi sertifikat SSL/TLS asli milik Google
		ServerName:         smtpHost,
	}

	address := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Hubungi server SMTP
	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		// Jika port 587 menolak TLS langsung di awal, coba jalur standar SMTP biasa dengan STARTTLS
		return sendMailStandard(address, auth, from, targetEmail, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	defer client.Quit()

	// Jalankan proses autentikasi dan pengiriman
	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(targetEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	fmt.Println("OTP Berhasil terkirim !")
	return nil
}

// Fallback method jika server membutuhkan jabat tangan STARTTLS biasa
func sendMailStandard(addr string, auth smtp.Auth, from string, to string, msg []byte) error {
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		fmt.Println("Gagal kirim email SMTP Standard:", err)
		return err
	}
	fmt.Println("OTP Berhasil terkirim via SMTP Standard! 🚀")
	return nil
}
