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
	ProcessCheckout(storeID, userID uint, input domain.CheckoutLaundryInput) (string, string, string, error)
	ProcessPelunasan(trxID, storeID uint, input domain.PelunasanInput) error
	GetLaporanRingkasan(storeID uint) (domain.ReportSummaryResponse, error)
	UpdateStatusCucian(trxID uint, status string) error
}

type laundryUseCase struct {
	repo repository.LaundryRepository
}

func NewLaundryUseCase(repo repository.LaundryRepository) LaundryUseCase {
	return &laundryUseCase{repo: repo}
}

// FUNGSI INTERNAL: Algoritma Pencari Rak Otomatis (Load Balancer)
func (u *laundryUseCase) CariRakTerbaik(storeID uint) (*domain.LaundryRack, error) {
	db := u.repo.GetDB()
	var rack domain.LaundryRack

	// 1. PRIORITAS UTAMA: Cari rak 100% KOSONG pakai NOT EXISTS (Anti Meleset!)
	db.Raw(`
        SELECT r.* FROM laundry_racks r
        WHERE r.store_id = ? AND r.status = 'TERSEDIA' 
        AND NOT EXISTS (
            SELECT 1 FROM laundry_transaction_details ltd 
            WHERE ltd.rack_id = r.id AND ltd.status_cucian NOT IN ('DIAMBIL', 'SELESAI')
        ) 
        ORDER BY r.zona ASC, r.baris ASC, r.kolom ASC LIMIT 1`, storeID).Scan(&rack)

	// 2. FALLBACK: Jika semua rak sudah terisi (ID = 0), cari rak dengan muatan paling sedikit
	if rack.ID == 0 {
		db.Raw(`
            SELECT r.* FROM laundry_racks r
            LEFT JOIN laundry_transaction_details ltd ON ltd.rack_id = r.id AND ltd.status_cucian NOT IN ('DIAMBIL', 'SELESAI')
            WHERE r.store_id = ? AND r.status = 'TERSEDIA'
            GROUP BY r.id
            ORDER BY COUNT(ltd.id) ASC, r.zona ASC, r.baris ASC, r.kolom ASC LIMIT 1`, storeID).Scan(&rack)
	}

	if rack.ID == 0 {
		return nil, fmt.Errorf("tidak ada rak yang tersedia")
	}

	return &rack, nil
}

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

	urlResult, errUpload := utils.UploadToSupabase(fileReader, fileName, "image/jpeg", bucketName, remotePath)
	if errUpload != nil {
		return "", fmt.Errorf("gagal upload ke supabase cloud: %v", errUpload)
	}

	return urlResult, nil
}

// 🚀 FUNGSI CHECKOUT KASTA TERTINGGI (ANTI ERROR FOREIGN KEY)
func (u *laundryUseCase) ProcessCheckout(storeID, userID uint, input domain.CheckoutLaundryInput) (string, string, string, error) {
	estimasiTime, err := time.Parse("2006-01-02", input.EstimasiSelesai)
	if err != nil {
		estimasiTime = time.Now().Add(time.Hour * 48)
	}

	invoiceCode := fmt.Sprintf("INV/LD/%s/%s", time.Now().Format("20060102"), time.Now().Format("150405"))
	trxPublicID := utils.GenerateULID()

	var buktiPath, fotoBarangPath string
	if input.PaymentMethod == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ = u.uploadBase64ToSupabase(input.BuktiTransferBase64, trxPublicID, "qris")
	}
	if input.FotoBarangBase64 != "" {
		fotoBarangPath, _ = u.uploadBase64ToSupabase(input.FotoBarangBase64, trxPublicID, "items")
	}

	// EKSEKUSI ALGORITMA RAK
	bestRack, errRack := u.CariRakTerbaik(storeID)
	rackID := uint(0)
	nomorRak := "-"

	if errRack == nil && bestRack != nil {
		rackID = bestRack.ID
		nomorRak = fmt.Sprintf("%s / %s", bestRack.Zona, bestRack.NamaRak)
	}

	db := u.repo.GetDB()
	tx := db.Begin()

	// 🚀 POLYFILL ENGINE: BIKIN SHIFT DUMMY SESUAI STRUCT CORE_SESSION LU BRAY!
	var activeSession models.CashierSession
	errSession := tx.Where("store_id = ? AND user_id = ? AND status = 'OPEN'", storeID, userID).Order("id DESC").First(&activeSession).Error

	if errSession != nil {
		// Bikin shift siluman yang lolos gembok constraint PostgreSQL lu
		activeSession = models.CashierSession{
			PublicID:      utils.GenerateULID(), // Wajib ada buat UUID
			StoreID:       storeID,
			UserID:        userID,
			StationNumber: "AUTO-LD", // Bypass unique index per stasiun
			ModalAwal:     0,
			Status:        "OPEN",
			OpenedAt:      time.Now(), // Nama kolom yang bener sesuai struct lu
			OpenedBy:      userID,     // Saksi forensik pembuka shift
		}
		if err := tx.Create(&activeSession).Error; err != nil {
			tx.Rollback()
			return "", "", "", fmt.Errorf("gagal membuat sesi kasir otomatis: %v", err)
		}
	}

	// 🚀 INJEKSI DATA TRANSAKSI
	newTx := models.Transaction{
		PublicID:      trxPublicID,
		SessionID:     activeSession.ID,
		StoreID:       storeID,
		UserID:        userID,
		NoInvoice:     invoiceCode,
		SubTotal:      input.TotalAmount,
		Pajak:         0,
		Pembulatan:    0,
		TotalHarga:    input.TotalAmount,
		MetodeBayar:   input.PaymentMethod,
		StatusBayar:   input.PaymentStatus,
		TipeBisnis:    "LAUNDRY",
		StatusPesanan: "ANTRI",
		NominalBayar:  input.TotalAmount,
		Kembalian:     0,
		BuktiTransfer: buktiPath,
	}

	if err := u.repo.CreateTransactionTx(tx, &newTx); err != nil {
		tx.Rollback()
		return "", "", "", fmt.Errorf("gagal membuat invoice induk: %v", err)
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
			RackID:        rackID,
			NomorRak:      nomorRak,
		}

		if err := u.repo.CreateLaundryDetailTx(tx, &laundryDetail); err != nil {
			tx.Rollback()
			return "", "", "", fmt.Errorf("gagal menyimpan rincian cucian: %v", err)
		}
	}

	existingCustomer, err := u.repo.FindCustomerByPhone(storeID, input.CustomerPhone)
	if err != nil {
		newCustomer := models.Customer{
			PublicID:   utils.GenerateULID(),
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
	return invoiceCode, fotoBarangPath, nomorRak, nil
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

		// 🚀 FIXED: Tambahin variabel buat nampung nomorRak
		var namaPelanggan, noWhatsapp, nomorRak string
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
			nomorRak = detail.NomorRak // 🎯 TARIK DATA RAK DARI DATABASE BRAY!
		}

		layananName := "Paket Laundry"
		satuanDasar := "KG" // 🚀 Default-nya KG

		if productID > 0 {
			if prod, err := u.repo.GetProductByIDSimple(productID); err == nil {
				layananName = prod.NamaProduk
				if prod.SatuanDasar != "" {
					satuanDasar = prod.SatuanDasar // 🎯 Tarik satuan asli dari Master Layanan!
				}
			}
		}

		reportData = append(reportData, domain.TransactionReportResponse{
			Transaction:   trx,
			Invoice:       trx.NoInvoice,
			Pelanggan:     namaPelanggan,
			Whatsapp:      noWhatsapp,
			Layanan:       layananName,
			BeratKg:       beratKg,
			SatuanDasar:   satuanDasar,
			SubTotal:      subTotalDetail,
			EstimasiWaktu: estimasiWaktu,
			NomorRak:      nomorRak,
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

func (u *laundryUseCase) UpdateStatusCucian(trxID uint, status string) error {
	db := u.repo.GetDB()

	// 1. 🚀 UPDATE TABEL INDUK (Biar Papan Kanban Vue lu bisa pindah kotak)
	err := db.Model(&models.Transaction{}).
		Where("id = ?", trxID).
		Update("status_pesanan", status).Error
	if err != nil {
		return fmt.Errorf("gagal merubah status pesanan induk: %v", err)
	}

	// 2. 🚀 UPDATE TABEL ANAK (Biar rincian cucian di dalamnya ikut sinkron)
	err = db.Model(&domain.TransactionLaundryDetail{}).
		Where("transaction_id = ?", trxID).
		Update("status_cucian", status).Error

	return err
}
