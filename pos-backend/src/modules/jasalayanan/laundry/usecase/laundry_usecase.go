package usecase

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
	"pos-backend/utils"
)

type LaundryUseCase interface {
	ProcessCheckout(storeID, userID uint, input domain.CheckoutLaundryInput) (string, string, error)
	ProcessPelunasan(trxID, storeID uint, input domain.PelunasanInput) error
	GetLaporanRingkasan(storeID uint) (domain.ReportSummaryResponse, error)
	UpdateStatusCucian(detailID uint, status string) error
}

type laundryUseCase struct {
	repo repository.LaundryRepository
}

func NewLaundryUseCase(repo repository.LaundryRepository) LaundryUseCase {
	return &laundryUseCase{repo: repo}
}

// 🚀 CLOUD STORAGE CONVERSION: Kirim base64 langsung meluncur ke Supabase Storage bray!
func (u *laundryUseCase) uploadBase64ToSupabase(base64Data, publicID, subFolder string) (string, error) {
	if base64Data == "" {
		return "", nil
	}

	parts := strings.Split(base64Data, ",")
	pureBase64 := parts[0]
	if len(parts) > 1 {
		pureBase64 = parts[1]
	}

	decodedData, err := base64.StdEncoding.DecodeString(pureBase64)
	if err != nil {
		return "", fmt.Errorf("gagal decode base64: %v", err)
	}

	bucketName := os.Getenv("SUPABASE_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "nexa-pos-storage"
	}

	remotePath := fmt.Sprintf("laundry/stores/%s/%s", publicID, subFolder)
	fileReader := bytes.NewReader(decodedData)
	fileName := fmt.Sprintf("%d.jpg", time.Now().UnixNano())

	// Tembak langsung ke awan nexa-pos-storage bray!
	urlResult, errUpload := utils.UploadToSupabase(fileReader, fileName, "image/jpeg", bucketName, remotePath)
	if errUpload != nil {
		return "", fmt.Errorf("gagal upload ke supabase cloud: %v", errUpload)
	}

	return urlResult, nil
}

func (u *laundryUseCase) ProcessCheckout(storeID, userID uint, input domain.CheckoutLaundryInput) (string, string, error) {
	estimasiTime, err := time.Parse("2006-01-02", input.EstimasiSelesai)
	if err != nil {
		estimasiTime = time.Now().Add(time.Hour * 48)
	}

	invoiceCode := fmt.Sprintf("INV/LD/%s/%s", time.Now().Format("20060102"), time.Now().Format("150405"))

	// 🛡️ SECURITY PATCH 1: Generate Public ID Transaksi unik berbasis UUID
	trxPublicID := utils.GenerateULID()

	// 🚀 UPLOAD LANGSUNG KE SUPABASE CLOUD (Folder public/uploads sudah punah!)
	var buktiPath, fotoBarangPath string
	if input.PaymentMethod == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ = u.uploadBase64ToSupabase(input.BuktiTransferBase64, trxPublicID, "qris")
	}
	if input.FotoBarangBase64 != "" {
		fotoBarangPath, _ = u.uploadBase64ToSupabase(input.FotoBarangBase64, trxPublicID, "items")
	}

	db := u.repo.GetDB()
	tx := db.Begin()

	newTx := models.Transaction{
		PublicID:      trxPublicID, // 🛡️ SECURITY PATCH 2: Kunci mati unique constraint DB lu bray!
		SessionID:     1,
		StoreID:       storeID,
		UserID:        userID,
		NoInvoice:     invoiceCode,
		SubTotal:      input.TotalAmount,
		Pajak:         0,
		Pembulatan:    0,
		TotalHarga:    input.TotalAmount,
		MetodeBayar:   input.PaymentMethod,
		StatusBayar:   input.PaymentStatus,
		NominalBayar:  input.TotalAmount,
		Kembalian:     0,
		BuktiTransfer: buktiPath,
	}

	if err := u.repo.CreateTransactionTx(tx, &newTx); err != nil {
		tx.Rollback()
		return "", "", fmt.Errorf("gagal membuat invoice induk: %v", err)
	}

	for _, item := range input.Items {
		laundryDetail := domain.TransactionLaundryDetail{
			TransactionID: newTx.ID,
			StoreID:       storeID,
			ProductID:     item.ProductID,
			NamaPelanggan: input.CustomerName,
			NoWhatsapp:    input.CustomerPhone,
			BeratKg:       item.BeratKg,
			HargaPerKg:    item.HargaPerKg,
			SubTotal:      item.SubTotal,
			StatusCucian:  "ANTRI",
			StatusBayar:   input.PaymentStatus,
			MetodeBayar:   input.PaymentMethod,
			EstimasiWaktu: estimasiTime,
			FotoBarang:    fotoBarangPath,
			NamaParfum:    item.NamaParfum,
			HargaParfum:   item.HargaParfum,
		}

		if err := u.repo.CreateLaundryDetailTx(tx, &laundryDetail); err != nil {
			tx.Rollback()
			return "", "", fmt.Errorf("gagal menyimpan rincian cucian: %v", err)
		}
	}

	existingCustomer, err := u.repo.FindCustomerByPhone(storeID, input.CustomerPhone)
	if err != nil {
		newCustomer := models.Customer{
			StoreID:    storeID,
			Nama:       input.CustomerName,
			NoWhatsapp: input.CustomerPhone,
		}
		tx.Create(&newCustomer)
	} else {
		existingCustomer.Nama = input.CustomerName
		tx.Save(existingCustomer)
	}

	tx.Commit()
	return invoiceCode, fotoBarangPath, nil
}

func (u *laundryUseCase) ProcessPelunasan(trxID, storeID uint, input domain.PelunasanInput) error {
	trx, err := u.repo.GetTransactionByID(trxID, storeID)
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan")
	}
	if trx.StatusBayar == "LUNAS" {
		return fmt.Errorf("transaksi ini sudah lunas sebelumnya")
	}
	if input.MetodeBayar == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ := u.uploadBase64ToSupabase(input.BuktiTransferBase64, trx.PublicID, "qris_lunas")
		trx.BuktiTransfer = buktiPath
	}
	trx.StatusBayar = "LUNAS"
	trx.MetodeBayar = input.MetodeBayar
	return u.repo.UpdateTransaction(trx)
}

func (u *laundryUseCase) GetLaporanRingkasan(storeID uint) (domain.ReportSummaryResponse, error) {
	transactions, err := u.repo.GetAllTransactions(storeID)
	if err != nil {
		return domain.ReportSummaryResponse{}, err
	}

	var reportData []domain.TransactionReportResponse
	var tunai, qris, debit, piutang, omset float64
	totalOrder := len(transactions)

	for _, trx := range transactions {
		if trx.StatusBayar == "BELUM_LUNAS" || trx.StatusBayar == "BELUM_BAYAR" {
			piutang += trx.TotalHarga
		} else {
			omset += trx.TotalHarga
			switch strings.ToUpper(trx.MetodeBayar) {
			case "TUNAI", "CASH":
				tunai += trx.TotalHarga
			case "QRIS":
				qris += trx.TotalHarga
			case "DEBIT":
				debit += trx.TotalHarga
			}
		}

		detail, err := u.repo.GetLaundryDetailByTxID(trx.ID)
		var namaPelanggan, noWhatsapp string
		var productID uint
		var beratKg, subTotalDetail float64
		var estimasiWaktu time.Time

		if err == nil {
			namaPelanggan = detail.NamaPelanggan
			noWhatsapp = detail.NoWhatsapp
			productID = detail.ProductID
			beratKg = detail.BeratKg
			subTotalDetail = detail.SubTotal
			estimasiWaktu = detail.EstimasiWaktu
		}

		layananName := "Paket Laundry"
		if productID > 0 {
			if prod, err := u.repo.GetProductByIDSimple(productID); err == nil {
				layananName = prod.NamaProduk
			}
		}

		reportData = append(reportData, domain.TransactionReportResponse{
			Transaction:   trx,
			Invoice:       trx.NoInvoice,
			Pelanggan:     namaPelanggan,
			Whatsapp:      noWhatsapp,
			Layanan:       layananName,
			BeratKg:       beratKg,
			SubTotal:      subTotalDetail,
			EstimasiWaktu: estimasiWaktu,
		})
	}

	avg := 0.0
	if totalOrder > 0 {
		avg = omset / float64(totalOrder)
	}

	return domain.ReportSummaryResponse{
		Ringkasan: domain.ReportSummary{
			TotalOmset: omset,
			TotalOrder: totalOrder,
			RataRata:   math.Round(avg*100) / 100,
			Tunai:      tunai,
			Qris:       qris,
			Debit:      debit,
			Piutang:    piutang,
		},
		Transaksi: reportData,
	}, nil
}

func (u *laundryUseCase) UpdateStatusCucian(detailID uint, status string) error {
	return u.repo.UpdateStatusDetailCucian(detailID, status)
}
