import Swal from 'sweetalert2';
import { computed, ref } from 'vue';
import api from '../../../api.js';

export function useInboundReport() {
	const rawData = ref([]);
	const isLoading = ref(false);

	// State Filter
	const searchQuery = ref('');
	const startDate = ref('');
	const endDate = ref('');
	const expandedRows = ref([]); // Untuk state buka-tutup accordion tabel

	const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);

	const formatDate = (dateStr) => {
		if (!dateStr) return '-';
		const d = new Date(dateStr);
		return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(d);
	};

	const fetchReports = async () => {
		isLoading.value = true;
		try {
			// Nanti backend Golang kita sesuaikan url-nya
			const response = await api.get('/retail/inbounds');
			rawData.value = response.data.data || response.data || [];
		} catch (error) {
			Swal.fire('Gagal!', 'Sirkuit laporan terputus dari server.', 'error');
		} finally {
			isLoading.value = false;
		}
	};

	const filteredReports = computed(() => {
		let result = rawData.value;

		if (searchQuery.value) {
			const q = searchQuery.value.toLowerCase();
			result = result.filter((r) => (r.nama_supplier && r.nama_supplier.toLowerCase().includes(q)) || (r.no_faktur && r.no_faktur.toLowerCase().includes(q)));
		}

		if (startDate.value && endDate.value) {
			const start = new Date(startDate.value).setHours(0, 0, 0, 0);
			const end = new Date(endDate.value).setHours(23, 59, 59, 999);
			result = result.filter((r) => {
				const itemDate = new Date(r.created_at).getTime();
				return itemDate >= start && itemDate <= end;
			});
		}

		return result;
	});

	const kpiStats = computed(() => {
		const totalModal = filteredReports.value.reduce((acc, curr) => acc + (curr.total_modal || 0), 0);
		const totalBarang = filteredReports.value.reduce((acc, curr) => acc + (curr.total_item || 0), 0);
		const totalFaktur = filteredReports.value.length;

		return { totalModal, totalBarang, totalFaktur };
	});

	const toggleRow = (id) => {
		const index = expandedRows.value.indexOf(id);
		if (index > -1) expandedRows.value.splice(index, 1);
		else expandedRows.value.push(id);
	};

	// 🚀 PRINT ENGINE: Auto generate halaman HTML khusus Faktur!
	const printFaktur = (lpb) => {
		const printWindow = window.open('', '_blank', 'width=800,height=600');
		if (!printWindow) return Swal.fire('Gagal', 'Popup diblokir oleh browser!', 'error');

		let itemsHtml = '';
		if (lpb.items && lpb.items.length > 0) {
			lpb.items.forEach((item, index) => {
				itemsHtml += `
                    <tr>
                        <td style="padding: 8px; border-bottom: 1px solid #e2e8f0; font-size: 12px; text-align: center;">${index + 1}</td>
                        <td style="padding: 8px; border-bottom: 1px solid #e2e8f0; font-size: 12px;">${item.nama_produk || '-'}</td>
                        <td style="padding: 8px; border-bottom: 1px solid #e2e8f0; font-size: 12px; text-align: center;">${item.qty || 0}</td>
                        <td style="padding: 8px; border-bottom: 1px solid #e2e8f0; font-size: 12px; text-align: right;">${formatRupiah(item.harga_modal)}</td>
                        <td style="padding: 8px; border-bottom: 1px solid #e2e8f0; font-size: 12px; text-align: right; font-weight: bold;">${formatRupiah(item.sub_total)}</td>
                    </tr>
                `;
			});
		}

		const htmlLayout = `
            <html>
            <head>
                <title>Cetak LPB - ${lpb.no_faktur}</title>
                <style>
                    body { font-family: 'Courier New', Courier, monospace; color: #1e293b; padding: 20px; }
                    .header { text-align: center; margin-bottom: 30px; border-bottom: 2px dashed #cbd5e1; padding-bottom: 20px; }
                    .title { font-size: 20px; font-weight: bold; margin: 0; text-transform: uppercase; letter-spacing: 2px; }
                    .info-box { display: flex; justify-content: space-between; margin-bottom: 20px; font-size: 12px; }
                    .info-box div { margin-bottom: 5px; }
                    table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
                    th { background-color: #f8fafc; padding: 10px; font-size: 12px; text-transform: uppercase; border-bottom: 2px solid #cbd5e1; text-align: left; }
                    .total-box { text-align: right; font-size: 14px; font-weight: bold; border-top: 2px dashed #cbd5e1; padding-top: 15px; }
                    .signature-box { display: flex; justify-content: space-between; margin-top: 50px; text-align: center; font-size: 12px; }
                    .sign-line { margin-top: 60px; border-top: 1px solid #1e293b; width: 150px; display: inline-block; }
                    @media print { body { padding: 0; } }
                </style>
            </head>
            <body>
                <div class="header">
                    <h1 class="title">Bukti Penerimaan Barang</h1>
                    <p style="margin: 5px 0 0 0; font-size: 12px;">Sistem Inventori Retail</p>
                </div>
                
                <div class="info-box">
                    <div>
                        <div><strong>No. Faktur:</strong> ${lpb.no_faktur}</div>
                        <div><strong>Tanggal Terima:</strong> ${formatDate(lpb.created_at)}</div>
                    </div>
                    <div style="text-align: right;">
                        <div><strong>Supplier:</strong> ${lpb.nama_supplier}</div>
                        <div><strong>Penerima:</strong> Administrator</div>
                    </div>
                </div>

                <table>
                    <thead>
                        <tr>
                            <th style="text-align: center; width: 5%;">No</th>
                            <th>Nama Barang</th>
                            <th style="text-align: center; width: 10%;">Qty</th>
                            <th style="text-align: right; width: 20%;">Modal/Pcs</th>
                            <th style="text-align: right; width: 20%;">Subtotal</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${itemsHtml}
                    </tbody>
                </table>

                <div class="total-box">
                    <div>TOTAL ITEM: ${lpb.total_item || 0}</div>
                    <div style="margin-top: 5px; font-size: 18px;">TOTAL NILAI: ${formatRupiah(lpb.total_modal)}</div>
                </div>

                <div class="signature-box">
                    <div>
                        <p>Supplier / Pengirim</p>
                        <span class="sign-line"></span>
                    </div>
                    <div>
                        <p>Penerima / Checker</p>
                        <span class="sign-line"></span>
                    </div>
                </div>
                
                <script>
                    window.onload = function() { window.print(); window.onafterprint = function(){ window.close(); } };
                </script>
            </body>
            </html>
        `;

		printWindow.document.write(htmlLayout);
		printWindow.document.close();
	};

	return {
		rawData,
		isLoading,
		searchQuery,
		startDate,
		endDate,
		expandedRows,
		filteredReports,
		kpiStats,
		formatRupiah,
		formatDate,
		fetchReports,
		toggleRow,
		printFaktur,
	};
}
