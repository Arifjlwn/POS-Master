<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import SidebarLaundry from '../components/SidebarLaundry.vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

const riwayat = ref([]);
const isLoading = ref(false);
const searchQuery = ref('');
let pollingInterval = null;

// 🚀 STATE BUAT NAMPUNG QRIS TOKO
const qrisTokoUrl = ref('https://upload.wikimedia.org/wikipedia/commons/d/d0/QR_code_for_mobile_English_Wikipedia.svg'); // Gambar default kalau toko belum setting QRIS

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);

const formatDate = (dateStr) => {
    if (!dateStr) return '-';
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return '-';
    return new Intl.DateTimeFormat('id-ID', { day: '2-digit', month: 'long', year: 'numeric' }).format(d);
};

// 🚀 FUNGSI AMBIL SETTING TOKO (BUAT NAMPILIN QRIS ASLI)
const fetchSettingToko = async () => {
    try {
        const response = await api.get('/laundry/setting');
        if (response.data && response.data.qris_image) {
            // Asumsi backend ngirim full URL atau path relative
            const url = response.data.qris_image.startsWith('http') 
                ? response.data.qris_image 
                : `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/${response.data.qris_image}`;
            qrisTokoUrl.value = url;
        }
    } catch (error) {
        console.error("Gagal mengambil QRIS Toko");
    }
};

// 📡 FETCH DATA CUCIAN
const fetchData = async (isBackground = false) => {
    if (!isBackground) isLoading.value = true;
    try {
        const response = await api.get('/laundry/report?period=tahun_ini'); 
        const semuaTransaksi = response.data.transaksi || [];
        riwayat.value = semuaTransaksi.filter(t => t.status_pesanan !== 'DIAMBIL');
    } catch (error) {
        console.error("Gagal sinkronisasi data cucian");
    } finally {
        if (!isBackground) isLoading.value = false;
    }
};

onMounted(() => {
    fetchSettingToko(); // Panggil setting toko dulu
    fetchData();
    pollingInterval = setInterval(() => fetchData(true), 10000); 
});

onUnmounted(() => { if (pollingInterval) clearInterval(pollingInterval); });

// KANBAN LOGIC
const orderAntri = computed(() => {
    return riwayat.value
        .filter(t => t.status_pesanan === 'ANTRI' || !t.status_pesanan)
        .sort((a, b) => new Date(a.estimasi_waktu) - new Date(b.estimasi_waktu)); 
});

const orderProses = computed(() => {
    return riwayat.value
        .filter(t => t.status_pesanan === 'PROSES')
        .sort((a, b) => new Date(a.estimasi_waktu) - new Date(b.estimasi_waktu));
});

const orderSelesai = computed(() => {
    return riwayat.value
        .filter(t => t.status_pesanan === 'SELESAI')
        .sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
});

// UPDATE STATUS KANBAN
const updateStatusKanban = async (trx, statusBaru) => {
    try {
        await api.put(`/laundry/transactions/${trx.id}/status`, { status_pesanan: statusBaru });
        fetchData(true);

        if (statusBaru === 'SELESAI' && trx.whatsapp) {
            Swal.fire({
                title: 'Kirim Notif WA?',
                text: "Cucian sudah selesai, mau kabarin pelanggan?",
                icon: 'question',
                showCancelButton: true,
                confirmButtonColor: '#10b981',
                confirmButtonText: 'Ya, Buka WhatsApp!'
            }).then((result) => {
                if (result.isConfirmed) {
                    const textRaw = `Halo Kak *${trx.pelanggan}*,\n\nCucian kakak dengan Resi: *${trx.invoice}* SUDAH SELESAI dicuci. 🥳✨\n\nSilakan diambil. Terima kasih!`;
                    window.open(`https://wa.me/${trx.whatsapp}?text=${encodeURIComponent(textRaw)}`, '_blank');
                }
            });
        }
    } catch (error) {
        Swal.fire('Error!', 'Gagal mengubah status', 'error');
    }
};

// 🚀 LOGIKA LUNASI DENGAN LIVE KAMERA JEPRET QRIS!
const prosesPengambilan = async (trx) => {
    if (trx.status_bayar === 'BELUM_LUNAS') {
        
        // 1. TANYA METODE BAYAR
        const { value: metode_bayar } = await Swal.fire({
            title: `Pelunasan Piutang`,
            html: `<p class="text-sm text-slate-500 mb-2">Tagihan <b>${trx.pelanggan}</b></p>
                   <p class="text-rose-500 font-black text-3xl mb-6 tracking-tighter">${formatRupiah(trx.sub_total || trx.total_harga)}</p>`,
            input: 'select',
            inputOptions: { 'TUNAI': 'Uang Tunai (Cash)', 'QRIS': 'Scan QRIS / Transfer', 'DEBIT': 'Mesin EDC / Debit' },
            showCancelButton: true, 
            confirmButtonText: 'Lanjut ➔', 
            confirmButtonColor: '#4f46e5',
            inputValidator: (value) => !value ? 'Pilih metode bayar!' : undefined
        });

        if (!metode_bayar) return;

        let buktiBase64 = "";

        // 2. LIVE KAMERA KHUSUS QRIS 📸
        if (metode_bayar === 'QRIS') {
            const { value: fileData } = await Swal.fire({
                title: 'Scan QRIS & Foto Bukti',
                html: `
                    <div class="bg-slate-50 p-3 rounded-2xl border border-slate-200 mb-3 inline-block shadow-inner">
                        <img src="${qrisTokoUrl.value}" alt="QRIS Toko" class="w-48 h-48 mx-auto rounded-xl object-cover mix-blend-multiply">
                    </div>
                    <p class="text-[10px] font-bold text-slate-500 mb-4 uppercase tracking-widest">Arahkan Pelanggan Scan QRIS</p>

                    <div id="kamera-box" class="hidden flex-col items-center gap-3 mb-2">
                        <div class="relative w-full h-48 bg-black rounded-2xl overflow-hidden shadow-inner">
                            <video id="live-video" autoplay playsinline class="w-full h-full object-cover"></video>
                        </div>
                        <button id="btn-jepret" type="button" class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2.5 rounded-full font-black text-[10px] uppercase tracking-widest shadow-lg flex items-center justify-center gap-2 active:scale-95 transition-all">
                            📸 Jepret Bukti Transfer
                        </button>
                    </div>

                    <div id="preview-box" class="hidden flex-col items-center gap-3 mb-2">
                        <img id="hasil-foto" class="w-full h-48 object-cover rounded-2xl border-4 border-emerald-100 shadow-md">
                        <button id="btn-ulangi" type="button" class="text-rose-500 font-bold text-[10px] uppercase underline active:scale-95">Ulangi Foto</button>
                    </div>

                    <button id="btn-buka-kamera" type="button" class="w-full bg-slate-800 hover:bg-slate-900 text-white py-3.5 rounded-xl font-black text-xs uppercase flex items-center justify-center gap-2 shadow-md transition-all active:scale-95">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                        Buka Kamera Kasir
                    </button>
                    <canvas id="hidden-canvas" class="hidden"></canvas>
                `,
                showCancelButton: true,
                confirmButtonText: 'Sahkan Pembayaran',
                confirmButtonColor: '#10b981',
                cancelButtonColor: '#94a3b8',
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

                    const matikanKamera = () => { if (stream) stream.getTracks().forEach(track => track.stop()); };

                    btnBuka.addEventListener('click', async () => {
                        try {
                            // Kamera belakang buat HP, webcam biasa buat Laptop
                            stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
                            video.srcObject = stream;
                            btnBuka.classList.add('hidden');
                            kameraBox.classList.remove('hidden');
                            kameraBox.classList.add('flex');
                        } catch (err) {
                            Swal.showValidationMessage('Akses kamera ditolak browser / tidak ada kamera!');
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
                        } catch(e){}
                    });
                },
                willClose: () => {
                    const video = document.getElementById('live-video');
                    if(video && video.srcObject) video.srcObject.getTracks().forEach(track => track.stop());
                },
                preConfirm: () => {
                    const hasilFoto = document.getElementById('hasil-foto');
                    if (!hasilFoto || !hasilFoto.dataset.base64) {
                        Swal.showValidationMessage('Harap jepret bukti transfer dulu!');
                        return false;
                    }
                    return hasilFoto.dataset.base64;
                }
            });

            if (!fileData) return; 
            buktiBase64 = fileData;
        }

        // 3. TEMBAK API GOLANG BUAT LUNASIN
        try {
            await api.put(`/laundry/transactions/${trx.id}/lunas`, { 
                metode_bayar: metode_bayar,
                bukti_transfer_base64: buktiBase64 
            });
            
            await updateStatusKanban(trx, 'DIAMBIL');
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Lunas & Diambil!', showConfirmButton: false, timer: 1500 });
        } catch (error) {
            Swal.fire('Error!', 'Gagal menyimpan pelunasan.', 'error');
        }

    } else {
        Swal.fire({
            title: 'Serahkan Cucian?', text: 'Cucian ini sudah LUNAS.', icon: 'info', showCancelButton: true, confirmButtonText: 'Tandai Diambil'
        }).then((result) => {
            if (result.isConfirmed) {
                updateStatusKanban(trx, 'DIAMBIL');
                Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Cucian diserahkan!', showConfirmButton: false, timer: 1500 });
            }
        });
    }
};
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-slate-50/50 overflow-hidden relative">
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-200 flex flex-col z-10 shadow-sm relative gap-5">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-sky-50 border border-sky-100 rounded-2xl flex items-center justify-center shrink-0"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-sky-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg></div>
                    <div>
                        <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Operasional Kanban</h1>
                        <p class="text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest flex items-center gap-2">Pantau & Serahkan Cucian <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span></p>
                    </div>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8 pb-20">
                <div v-if="isLoading && riwayat.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-400">
                    <div class="w-12 h-12 border-4 border-sky-200 border-t-sky-600 rounded-full animate-spin mb-4"></div><p class="font-black text-xs uppercase tracking-[0.2em] animate-pulse">Memuat Data...</p>
                </div>
                
                <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div class="bg-slate-100/50 p-4 rounded-[24px] border border-slate-200 flex flex-col h-[70vh]">
                        <div class="flex items-center justify-between mb-4 px-2"><h3 class="font-black text-xs uppercase tracking-widest text-slate-700 flex items-center gap-2"><div class="w-2.5 h-2.5 rounded-full bg-rose-400"></div> Antrean Masuk</h3><span class="bg-slate-200 text-slate-600 text-[10px] font-black px-2.5 py-1 rounded-lg">{{ orderAntri.length }}</span></div>
                        <div class="flex-1 overflow-y-auto hide-scrollbar space-y-3">
                            <div v-for="order in orderAntri" :key="order.id" class="bg-white p-4 rounded-2xl shadow-sm border border-slate-100">
                                <div class="flex justify-between items-start mb-2">
                                    <span class="text-[9px] font-black uppercase tracking-widest text-indigo-500">{{ order.invoice }}</span>
                                </div>
                                <p class="font-black text-sm text-slate-800 uppercase mb-1">{{ order.pelanggan }}</p>
                                
                                <div class="flex items-center gap-1 text-[10px] font-bold text-rose-500 mb-3 bg-rose-50 w-fit px-2 py-0.5 rounded-md border border-rose-100">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                    Target: {{ formatDate(order.estimasi_waktu) }}
                                </div>

                                <div class="flex justify-between items-center mt-3 pt-3 border-t border-slate-50">
                                    <span v-if="order.status_bayar === 'BELUM_LUNAS'" class="text-[8px] font-black bg-rose-50 text-rose-500 px-2 py-0.5 rounded uppercase">Piutang</span>
                                    <span v-else class="text-[8px] font-black bg-emerald-50 text-emerald-600 px-2 py-0.5 rounded uppercase">Lunas</span>
                                    <button @click="updateStatusKanban(order, 'PROSES')" class="bg-indigo-600 hover:bg-indigo-700 text-white text-[9px] font-black px-3 py-1.5 rounded-lg uppercase tracking-widest transition-all shadow-sm">Mulai Cuci ➔</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="bg-slate-100/50 p-4 rounded-[24px] border border-slate-200 flex flex-col h-[70vh]">
                        <div class="flex items-center justify-between mb-4 px-2"><h3 class="font-black text-xs uppercase tracking-widest text-slate-700 flex items-center gap-2"><div class="w-2.5 h-2.5 rounded-full bg-sky-400 animate-pulse"></div> Sedang Dicuci</h3><span class="bg-slate-200 text-slate-600 text-[10px] font-black px-2.5 py-1 rounded-lg">{{ orderProses.length }}</span></div>
                        <div class="flex-1 overflow-y-auto hide-scrollbar space-y-3">
                            <div v-for="order in orderProses" :key="order.id" class="bg-white p-4 rounded-2xl shadow-sm border border-slate-100">
                                <div class="flex justify-between items-start mb-2"><span class="text-[9px] font-black uppercase tracking-widest text-sky-500">{{ order.invoice }}</span></div>
                                <p class="font-black text-sm text-slate-800 uppercase mb-1">{{ order.pelanggan }}</p>
                                
                                <div class="flex items-center gap-1 text-[10px] font-bold text-sky-600 mb-3 bg-sky-50 w-fit px-2 py-0.5 rounded-md border border-sky-100">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                    Target: {{ formatDate(order.estimasi_waktu) }}
                                </div>

                                <div class="flex justify-between items-center mt-3 pt-3 border-t border-slate-50">
                                    <button @click="updateStatusKanban(order, 'ANTRI')" class="text-[9px] font-black text-slate-400 hover:text-slate-600 uppercase">Batal</button>
                                    <button @click="updateStatusKanban(order, 'SELESAI')" class="bg-emerald-500 hover:bg-emerald-600 text-white text-[9px] font-black px-3 py-1.5 rounded-lg uppercase tracking-widest transition-all shadow-sm">Selesai ➔</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="bg-slate-100/50 p-4 rounded-[24px] border border-slate-200 flex flex-col h-[70vh]">
                        <div class="flex items-center justify-between mb-4 px-2"><h3 class="font-black text-xs uppercase tracking-widest text-slate-700 flex items-center gap-2"><div class="w-2.5 h-2.5 rounded-full bg-emerald-400"></div> Siap Diambil</h3><span class="bg-slate-200 text-slate-600 text-[10px] font-black px-2.5 py-1 rounded-lg">{{ orderSelesai.length }}</span></div>
                        <div class="flex-1 overflow-y-auto hide-scrollbar space-y-3">
                            <div v-for="order in orderSelesai" :key="order.id" class="bg-white p-4 rounded-2xl shadow-sm border border-emerald-200">
                                <div class="flex justify-between items-start mb-2"><span class="text-[9px] font-black uppercase tracking-widest text-emerald-600">{{ order.invoice }}</span></div>
                                <p class="font-black text-sm text-slate-800 uppercase mb-3">{{ order.pelanggan }}</p>
                                <div class="flex justify-end items-center mt-3 pt-3 border-t border-slate-50">
                                    <button v-if="order.status_bayar === 'BELUM_LUNAS'" @click="prosesPengambilan(order)" class="w-full bg-rose-500 hover:bg-rose-600 text-white text-[9px] font-black px-3 py-2.5 rounded-lg uppercase tracking-[0.2em] transition-all shadow-lg shadow-rose-200 animate-pulse">
                                        Bayar & Serahkan
                                    </button>
                                    <button v-else @click="prosesPengambilan(order)" class="w-full bg-slate-800 hover:bg-slate-900 text-white text-[9px] font-black px-3 py-2.5 rounded-lg uppercase tracking-widest transition-all shadow-md">
                                        Serahkan Cucian (LUNAS)
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </SidebarLaundry>
</template>
<style scoped>.hide-scrollbar::-webkit-scrollbar { display: none; } @keyframes fadeIn { from { opacity: 0; transform: translateY(15px); } to { opacity: 1; transform: translateY(0); } }</style>