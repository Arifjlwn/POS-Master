import Swal from 'sweetalert2';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import api from '../../../../api.js';

export function useKanbanLaundry() {
	const riwayat = ref([]);
	const isLoading = ref(false);
	const searchQuery = ref('');
	let pollingInterval = null;

	const qrisTokoUrl = ref('https://upload.wikimedia.org/wikipedia/commons/d/d0/QR_code_for_mobile_English_Wikipedia.svg');

	const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);

	const formatDate = (dateStr) => {
		if (!dateStr) return '-';
		const d = new Date(dateStr);
		if (isNaN(d.getTime())) return '-';
		return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }).format(d);
	};

	const fetchSettingToko = async () => {
		try {
			// 🚀 FIXED: Ganti URL ke jalur yang bener (punya Global Store)
			const response = await api.get('/store/settings');
			const data = response.data.data || response.data;

			if (data && data.qris_image) {
				const url = data.qris_image.startsWith('http') ? data.qris_image : `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/${data.qris_image}`;
				qrisTokoUrl.value = url;
			}
		} catch (error) {
			console.error('Gagal mengambil QRIS Toko', error);
		}
	};

	const fetchData = async (isBackground = false) => {
		if (!isBackground) isLoading.value = true;
		try {
			const response = await api.get('/laundry/report?period=tahun_ini');
			const semuaTransaksi = response.data.transaksi || [];
			riwayat.value = semuaTransaksi.filter((t) => t.status_pesanan !== 'DIAMBIL');
		} catch (error) {
			console.error('Gagal sinkronisasi data cucian', error);
		} finally {
			if (!isBackground) isLoading.value = false;
		}
	};

	const orderAntri = computed(() => {
		return riwayat.value
			.filter((t) => t.status_pesanan === 'ANTRI' || !t.status_pesanan)
			.filter((t) => t.pelanggan.toLowerCase().includes(searchQuery.value.toLowerCase()) || t.invoice.toLowerCase().includes(searchQuery.value.toLowerCase()))
			.sort((a, b) => new Date(a.estimasi_waktu) - new Date(b.estimasi_waktu));
	});

	const orderProses = computed(() => {
		return riwayat.value
			.filter((t) => t.status_pesanan === 'PROSES')
			.filter((t) => t.pelanggan.toLowerCase().includes(searchQuery.value.toLowerCase()) || t.invoice.toLowerCase().includes(searchQuery.value.toLowerCase()))
			.sort((a, b) => new Date(a.estimasi_waktu) - new Date(b.estimasi_waktu));
	});

	const orderSelesai = computed(() => {
		return riwayat.value
			.filter((t) => t.status_pesanan === 'SELESAI')
			.filter((t) => t.pelanggan.toLowerCase().includes(searchQuery.value.toLowerCase()) || t.invoice.toLowerCase().includes(searchQuery.value.toLowerCase()))
			.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
	});

	const updateStatusKanban = async (trx, statusBaru) => {
		try {
			await api.put(`/laundry/transactions/${trx.id}/status`, { status_pesanan: statusBaru });
			fetchData(true);

			if (statusBaru === 'SELESAI' && trx.whatsapp) {
				Swal.fire({
					title: 'Kirim Notifikasi WA?',
					text: 'Pesanan telah selesai. Kirim pesan otomatis ke pelanggan?',
					icon: 'question',
					showCancelButton: true,
					confirmButtonColor: '#10b981',
					confirmButtonText: 'Buka WhatsApp',
					cancelButtonText: 'Lewati',
					customClass: { popup: 'rounded-[24px]' },
				}).then((result) => {
					if (result.isConfirmed) {
						const textRaw = `Halo Kak *${trx.pelanggan}*,\n\nPesanan laundry dengan No. Resi: *${trx.invoice}* telah SELESAI diproses dan siap untuk diambil.\n\nTerima kasih telah mempercayakan layanan kami!`;
						window.open(`https://wa.me/${trx.whatsapp}?text=${encodeURIComponent(textRaw)}`, '_blank');
					}
				});
			}
		} catch (error) {
			Swal.fire('Error!', 'Gagal mengubah status', 'error');
		}
	};

	const prosesPengambilan = async (trx) => {
		if (trx.status_bayar === 'BELUM_LUNAS') {
			const { value: metode_bayar } = await Swal.fire({
				title: `Pelunasan Tagihan`,
				html: `
                    <div class="text-center mb-4">
                        <p class="text-xs font-black uppercase tracking-widest text-slate-500 mb-1">${trx.pelanggan}</p>
                        <p class="text-rose-500 font-black text-4xl tracking-tighter">${formatRupiah(trx.sub_total || trx.total_harga)}</p>
                    </div>
                `,
				input: 'select',
				inputOptions: { TUNAI: 'Uang Tunai (Cash)', QRIS: 'Scan QRIS / Transfer', DEBIT: 'Mesin EDC / Debit' },
				showCancelButton: true,
				confirmButtonText: 'Lanjutkan',
				cancelButtonText: 'Batal',
				confirmButtonColor: '#4f46e5',
				customClass: { popup: 'rounded-[32px]', input: 'rounded-xl border-slate-200 text-sm font-bold p-3' },
				inputValidator: (value) => (!value ? 'Pilih metode pembayaran terlebih dahulu!' : undefined),
			});

			if (!metode_bayar) return;

			let buktiBase64 = '';

			if (metode_bayar === 'QRIS') {
				const { value: fileData } = await Swal.fire({
					title: 'Verifikasi QRIS',
					html: `
                        <div class="bg-slate-50 p-4 rounded-3xl border border-slate-200 mb-4 inline-block shadow-inner">
                            <img src="${qrisTokoUrl.value}" alt="QRIS Toko" class="w-48 h-48 mx-auto rounded-xl object-contain mix-blend-multiply">
                        </div>
                        
                        <div id="kamera-box" class="hidden flex-col items-center gap-3 mb-2 w-full">
                            <div class="relative w-full h-56 bg-slate-900 rounded-2xl overflow-hidden shadow-inner ring-4 ring-slate-100">
                                <video id="live-video" autoplay playsinline class="w-full h-full object-cover"></video>
                            </div>
                            <button id="btn-jepret" type="button" class="w-full bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3.5 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg flex items-center justify-center gap-2 active:scale-95 transition-all">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><circle cx="12" cy="13" r="3" /></svg>
                                Ambil Gambar Bukti
                            </button>
                        </div>

                        <div id="preview-box" class="hidden flex-col items-center gap-3 mb-2 w-full">
                            <img id="hasil-foto" class="w-full h-56 object-cover rounded-2xl border-4 border-emerald-100 shadow-md">
                            <button id="btn-ulangi" type="button" class="text-rose-500 font-bold text-[10px] uppercase tracking-widest underline active:scale-95 py-2">Ulangi Pengambilan</button>
                        </div>

                        <button id="btn-buka-kamera" type="button" class="w-full bg-slate-800 hover:bg-slate-900 text-white py-3.5 rounded-xl font-black text-[10px] tracking-widest uppercase flex items-center justify-center gap-2 shadow-md transition-all active:scale-95">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                            Buka Kamera Kasir
                        </button>
                        <canvas id="hidden-canvas" class="hidden"></canvas>
                    `,
					showCancelButton: true,
					confirmButtonText: 'Sahkan Pembayaran',
					cancelButtonText: 'Batal',
					confirmButtonColor: '#10b981',
					cancelButtonColor: '#94a3b8',
					customClass: { popup: 'rounded-[32px] p-6', confirmButton: 'rounded-xl uppercase font-black text-[10px] tracking-widest py-3 px-6', cancelButton: 'rounded-xl uppercase font-black text-[10px] tracking-widest py-3 px-6' },
					didOpen: () => {
						const btnBuka = document.getElementById('btn-buka-kamera');
						const btnJepret = document.getElementById('btn-jepret');
						const btnUlangi = document.getElementById('btn-ulangi');
						const kameraBox = document.getElementById('kamera-box');
						const previewBox = document.getElementById('preview-box');
						const video = document.getElementById('live-video');
						const canvas = document.getElementById('hidden-canvas');
						const hasilFoto = document.getElementById('hasil-foto');

						let stream = null;

						const matikanKamera = () => {
							if (stream) stream.getTracks().forEach((track) => track.stop());
						};

						btnBuka.addEventListener('click', async () => {
							try {
								stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
								video.srcObject = stream;
								btnBuka.classList.add('hidden');
								kameraBox.classList.remove('hidden');
								kameraBox.classList.add('flex');
							} catch (err) {
								Swal.showValidationMessage('Akses kamera ditolak / tidak ditemukan!');
							}
						});

						btnJepret.addEventListener('click', () => {
							canvas.width = video.videoWidth;
							canvas.height = video.videoHeight;
							canvas.getContext('2d').drawImage(video, 0, 0);
							const base64Data = canvas.toDataURL('image/jpeg', 0.8);

							hasilFoto.src = base64Data;
							hasilFoto.dataset.base64 = base64Data;

							kameraBox.classList.add('hidden');
							kameraBox.classList.remove('flex');
							previewBox.classList.remove('hidden');
							previewBox.classList.add('flex');

							matikanKamera();
						});

						btnUlangi.addEventListener('click', async () => {
							previewBox.classList.add('hidden');
							previewBox.classList.remove('flex');
							hasilFoto.dataset.base64 = '';
							try {
								stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
								video.srcObject = stream;
								kameraBox.classList.remove('hidden');
								kameraBox.classList.add('flex');
							} catch (e) {}
						});
					},
					willClose: () => {
						const video = document.getElementById('live-video');
						if (video && video.srcObject) video.srcObject.getTracks().forEach((track) => track.stop());
					},
					preConfirm: () => {
						const hasilFoto = document.getElementById('hasil-foto');
						if (!hasilFoto || !hasilFoto.dataset.base64) {
							Swal.showValidationMessage('Bukti transfer wajib diunggah!');
							return false;
						}
						return hasilFoto.dataset.base64;
					},
				});

				if (!fileData) return;
				buktiBase64 = fileData;
			}

			try {
				await api.put(`/laundry/transactions/${trx.id}/lunas`, {
					metode_bayar: metode_bayar,
					bukti_transfer_base64: buktiBase64,
				});

				await updateStatusKanban(trx, 'DIAMBIL');
				Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Lunas & Diserahkan!', showConfirmButton: false, timer: 1500 });
			} catch (error) {
				Swal.fire('Error!', 'Sistem menolak penyimpanan pelunasan.', 'error');
			}
		} else {
			Swal.fire({
				title: 'Serahkan Cucian?',
				text: 'Pesanan ini sudah berstatus LUNAS.',
				icon: 'info',
				showCancelButton: true,
				confirmButtonText: 'Tandai Selesai & Diambil',
				cancelButtonText: 'Batal',
				confirmButtonColor: '#10b981',
				customClass: { popup: 'rounded-[24px]' },
			}).then((result) => {
				if (result.isConfirmed) {
					updateStatusKanban(trx, 'DIAMBIL');
					Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Pesanan ditutup!', showConfirmButton: false, timer: 1500 });
				}
			});
		}
	};

	onMounted(() => {
		fetchSettingToko();
		fetchData();
		pollingInterval = setInterval(() => fetchData(true), 10000);
	});

	onUnmounted(() => {
		if (pollingInterval) clearInterval(pollingInterval);
	});

	return {
		riwayat,
		isLoading,
		searchQuery,
		orderAntri,
		orderProses,
		orderSelesai,
		formatDate,
		updateStatusKanban,
		prosesPengambilan,
	};
}
