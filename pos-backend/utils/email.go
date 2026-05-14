package utils

import (
	"fmt"
	"net/smtp"
)

func SendOTPEmail(targetEmail string, otp string) error {
	// 🚀 DATA DARI MAILTRAP (Ganti Username & Password sesuai di web Mailtrap Mas)
	from := "no-reply@posumkm.com" // Bebas mau tulis apa aja kalau di Mailtrap
	username := "da1bed098977f5"
	password := "3f5706f134b6eb"
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := "2525" // Port standar Mailtrap

	// Format Email (HTML agar keren)
	subject := "Subject: [POSUMKM] Kode Verifikasi Akun\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 400px; margin: auto; border: 1px solid #eee; padding: 20px; border-radius: 20px;">
			<h2 style="color: #2563eb; text-align: center;">POS<span style="color: #1e293b;">UMKM</span></h2>
			<p style="text-align: center; color: #64748b;">Gunakan kode OTP di bawah untuk verifikasi:</p>
			<div style="font-size: 40px; font-weight: 900; letter-spacing: 10px; text-align: center; color: #2563eb; padding: 20px; background: #f8fafc; border-radius: 15px; margin: 20px 0;">
				%s
			</div>
			<p style="font-size: 10px; color: #94a3b8; text-align: center; text-transform: uppercase; letter-spacing: 1px;">Kode berlaku 5 menit</p>
		</div>
	`, otp)

	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{targetEmail}, []byte(subject+mime+body))
	
	if err != nil {
		fmt.Println("Gagal kirim email Mailtrap:", err)
		return err
	}
	
	fmt.Println("OTP Berhasil terkirim ke Sandbox Mailtrap! 🚀")
	return nil
}