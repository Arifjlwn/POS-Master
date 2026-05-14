<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';
import * as faceapi from 'face-api.js';

// --- STATE WAKTU REALTIME ---
const currentTime = ref('');
let timer;

const updateTime = () => {
    const now = new Date();
    currentTime.value = now.toLocaleTimeString('id-ID', {
        hour: '2-digit', minute: '2-digit', second: '2-digit'
    }).replace(/:/g, '.');
};

// --- STATE DATA ---
const karyawan = ref([]);

const riwayat = ref([]);
const urutanTanggalTerbaru = ref(true); // true = terbaru ke tertua, false = tertua ke terbaru

const toggleSortTanggal = () => {
    urutanTanggalTerbaru.value = !urutanTanggalTerbaru.value;
    
    // Langsung urutkan ulang array riwayat yang sudah ada di memori Vue
    if (urutanTanggalTerbaru.value) {
        riwayat.value.sort((a, b) => b.tanggal.localeCompare(a.tanggal));
    } else {
        riwayat.value.sort((a, b) => a.tanggal.localeCompare(b.tanggal));
    }
};

const isLoading = ref(true);
const filterMode = ref('harian');
const bulanDipilih = ref(new Date().toISOString().slice(0, 7)); // YYYY-MM
const isAiLoading = ref(true);

const getPayloadFromToken = () => {
    const token = localStorage.getItem('token');
    const name = localStorage.getItem('name') || 'User';
    const role = localStorage.getItem('role') || 'kasir';
    
    if (!token) return { user_id: 0, role: '', name: '' };
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return { user_id: payload.user_id, role: role.toLowerCase(), name: name };
    } catch (e) {
        return { user_id: 0, role: '', name: '' };
    }
};
const currentUser = ref(getPayloadFromToken());

const today = new Date().toLocaleDateString('en-CA');
const tanggalDipilih = ref(today);

// --- 🚀 FIX PERMANEN: SAMAKAN FORMAT STRING TANGGAL SECARA MURNI ---
const cekAbsenMasukKaryawan = (userId) => {
    return riwayat.value.some(log => {
        if (!log.tanggal) return false;
        const tglClean = log.tanggal.substring(0, 10);
        return Number(log.user_id) === Number(userId) && 
               tglClean === today && 
               log.jam_masuk != null;
    });
};

const cekAbsenPulangKaryawan = (userId) => {
    return riwayat.value.some(log => {
        if (!log.tanggal) return false;
        const tglClean = log.tanggal.substring(0, 10);
        return Number(log.user_id) === Number(userId) && 
               tglClean === today && 
               log.jam_pulang != null && log.jam_pulang !== "";
    });
};

// --- STATE KAMERA ---
const showCameraModal = ref(false);
const videoRef = ref(null);
const stream = ref(null);
const absenTarget = ref({ id: null, nama: '', jenis: '' });
const isSubmitting = ref(false);

// --- SINKRONISASI LOG RIWAYAT HARIAN & BULANAN PENUH ALA RETAIL ---
const fetchData = async () => {
    isLoading.value = true;
    try {
        // Ambil semua daftar karyawan murni dari backend Go (Diizinkan untuk semua role)
        const resKaryawan = await api.get('/employees');
        const allEmp = resKaryawan.data.data || [];
        
        // Pisahkan pasukan staff biasa (yang rolenya bukan owner)
        const staffSaja = allEmp.filter(e => e.role !== 'owner');

        // 🚀 LOGIKA DINAMIS SINKRONISASI PANEL ATAS
        if (currentUser.value.role === 'owner') {
            // Jika login Owner dan belum punya karyawan, tampilkan diri sendiri agar bisa absen solo
            if (staffSaja.length === 0) {
                karyawan.value = allEmp; 
            } else {
                karyawan.value = staffSaja; // Jika ada staff, owner ngalah keluar dari panel biar rapi
            }
        } else {
            // Jika login sebagai Staff biasa, dia tetap melihat SEMUA staff rekan kerjanya agar transparan!
            karyawan.value = staffSaja;
        }

        // Tentukan range tanggal awal dan akhir bulan
        const [tahun, bulan] = bulanDipilih.value.split('-');
        const tanggalAwalBulan = `${bulanDipilih.value}-01`;
        const tanggalAkhirBulan = `${bulanDipilih.value}-${new Date(tahun, bulan, 0).getDate()}`;

        const params = {};
        if (filterMode.value === 'harian') {
            params.tanggal = tanggalDipilih.value;
        } else {
            params.bulan = bulan;
            params.tahun = tahun;
        }

        // Ambil data absensi aktual dan master jadwal dari backend
        const [resRiwayat, resSched] = await Promise.all([
            api.get('/attendance', { params }),
            api.get('/schedules', { params: { 
                start_date: filterMode.value === 'harian' ? tanggalDipilih.value : tanggalAwalBulan, 
                end_date: filterMode.value === 'harian' ? tanggalDipilih.value : tanggalAkhirBulan 
            } })
        ]);

        const dataAbsenReal = resRiwayat.data.data || [];
        const dataJadwalReal = resSched.data.data || [];
        const matriksGabungan = [];

        // --- GENERATOR MATRIKS LOG RIWAYAT TABEL ---
        if (filterMode.value === 'harian') {
            karyawan.value.forEach(emp => {
                const empKey = emp.id || emp.user_id;
                const jadwalHariIni = dataJadwalReal.find(s => Number(s.user_id) === Number(empKey) && s.tanggal.substring(0, 10) === tanggalDipilih.value);
                const absenHariIni = dataAbsenReal.find(a => Number(a.user_id) === Number(empKey) && a.tanggal.substring(0, 10) === tanggalDipilih.value);

                if (jadwalHariIni && jadwalHariIni.shift_type === 'OFF') return;

                matriksGabungan.push({
                    id: absenHariIni?.id || `temp-harian-${empKey}`,
                    user_id: empKey,
                    tanggal: tanggalDipilih.value,
                    User: emp,
                    shift: jadwalHariIni ? jadwalHariIni.shift_type.replace(' (Approved)','').replace(' (Pending)','') : 'Belum Set',
                    foto_masuk: absenHariIni?.foto_masuk || null,
                    jam_masuk: absenHariIni?.jam_masuk || null,
                    foto_pulang: absenHariIni?.foto_pulang || null,
                    jam_pulang: absenHariIni?.jam_pulang || null,
                    status: absenHariIni?.status || (tanggalDipilih.value < today ? 'Mangkir' : 'Belum Absen')
                });
            });
        } else {
            const jumlahHari = new Date(tahun, bulan, 0).getDate();

            karyawan.value.forEach(emp => {
                const empKey = emp.id || emp.user_id;

                for (let hari = 1; hari <= jumlahHari; hari++) {
                    const tglLoopStr = `${bulanDipilih.value}-${String(hari).padStart(2, '0')}`;
                    
                    const absenMatch = dataAbsenReal.find(a => Number(a.user_id) === Number(empKey) && a.tanggal.substring(0, 10) === tglLoopStr);
                    const jadwalMatch = dataJadwalReal.find(s => Number(s.user_id) === Number(empKey) && s.tanggal.substring(0, 10) === tglLoopStr);

                    const shiftClean = jadwalMatch ? jadwalMatch.shift_type.replace(' (Approved)','').replace(' (Pending)','') : 'Belum Set';

                    let statusDinamis = 'Belum Absen';
                    if (absenMatch) {
                        statusDinamis = absenMatch.status;
                    } else if (shiftClean === 'OFF') {
                        statusDinamis = 'Libur (OFF)';
                    } else if (emp.role === 'owner') {
                        statusDinamis = 'Owner';
                    } else if (tglLoopStr < today) {
                        statusDinamis = 'Mangkir';
                    }

                    if (tglLoopStr <= today || absenMatch) {
                        matriksGabungan.push({
                            id: absenMatch?.id || `temp-bulan-${empKey}-${tglLoopStr}`,
                            user_id: empKey,
                            tanggal: tglLoopStr,
                            User: emp,
                            shift: shiftClean,
                            foto_masuk: absenMatch?.foto_masuk || null,
                            jam_masuk: absenMatch?.jam_masuk || null,
                            foto_pulang: absenMatch?.foto_pulang || null,
                            jam_pulang: absenMatch?.jam_pulang || null,
                            status: statusDinamis
                        });
                    }
                }
            });
        }

        riwayat.value = [...matriksGabungan].sort((a, b) => b.tanggal.localeCompare(a.tanggal));
    } catch (error) {
        console.error("Gagal sinkronisasi data log:", error);
    } finally {
        isLoading.value = false;
    }
};

const lihatFoto = (url, nama, tipe, jam) => {
    if (!url) return;
    Swal.fire({
        title: `Foto ${tipe}`,
        html: `<div class="text-left font-bold text-sm">👤 Karyawan: ${nama}<br>⏰ Jam: ${jam}</div>`,
        imageUrl: url,
        imageAlt: 'Foto Absensi',
        confirmButtonText: 'Tutup',
        confirmButtonColor: '#2563eb',
        customClass: { image: 'rounded-2xl border-4 border-gray-100 shadow-lg' }
    });
};

// LOGIC AI ABSENSI
onMounted(async () => {
    updateTime();
    timer = setInterval(updateTime, 1000);
    
    // 🚀 LOAD MODELS AI SAAT HALAMAN DIBUKA
    try {
        await Promise.all([
            faceapi.nets.tinyFaceDetector.loadFromUri('/models'),
            faceapi.nets.faceLandmark68Net.loadFromUri('/models'),
            faceapi.nets.faceRecognitionNet.loadFromUri('/models')
        ]);
        isAiLoading.value = false;
        console.log("AI Models Loaded! 🤖");
    } catch (e) {
        console.error("Gagal load AI models", e);
    }
    
    fetchData();
});

onUnmounted(() => {
    clearInterval(timer);
    stopCamera();
});

watch([tanggalDipilih, bulanDipilih, filterMode], () => fetchData());

// Fungsi untuk membatasi tanggal pengisian jadwal di Vue Mas Arif
const getJadwalRange = () => {
    const sekarang = new Date();
    const tahun = sekarang.getFullYear();
    const bulan = sekarang.getMonth(); // Bulan saat ini (0-11)
    
    // Cari tanggal terakhir di bulan ini
    const tanggalTerakhir = new Date(tahun, bulan + 1, 0).getDate(); 
    const hariTanggalTerakhir = new Date(tahun, bulan, tanggalTerakhir).getDay(); // 0 = Minggu, 4 = Kamis, dll
    
    console.log(`Tanggal terakhir bulan ini: ${tanggalTerakhir}, Jatuh pada hari ke-${hariTanggalTerakhir}`);
    
    // Logika pembatasan ala TSM Indomaret Mas Arif:
    // Jika akhir bulan bukan hari Minggu (0), berarti ada potongan minggu (Split-Week)
    if (hariTanggalTerakhir !== 0) {
        const sisaHari = 7 - hariTanggalTerakhir;
        console.log(`Karyawan harus buat jadwal parsial sisa bulan sebanyak ${sisaHari} hari ke depan.`);
        // Jalur ini yang mengunci form biar cuma bisa isi dari tanggal 1 sampai hari Minggu pertama di bulan baru
    }
}

const mulaiAbsen = async (id, nama, jenis) => {
    absenTarget.value = { id, nama, jenis };
    showCameraModal.value = true;
    try {
        stream.value = await navigator.mediaDevices.getUserMedia({ 
            video: { facingMode: "user", width: 1280, height: 720 } 
        });
        await nextTick();
        if (videoRef.value) videoRef.value.srcObject = stream.value;
    } catch (err) {
        Swal.fire('Kamera Error', 'Nyalakan izin kamera ya!', 'error');
        showCameraModal.value = false;
    }
};

const stopCamera = () => {
    if (stream.value) stream.value.getTracks().forEach(track => track.stop());
    showCameraModal.value = false;
};

const jepretDanKirim = async () => {
    if (isAiLoading.value) return Swal.fire('Tunggu', 'AI sedang pemanasan...', 'warning');
    
    isSubmitting.value = true;
    
    try {
        // 1. Ambil data profil diri sendiri (Pasti diizinkan Backend)
        const resMe = await api.get('/me'); 
        const targetEmp = resMe.data; // Mengambil data user yang sedang login

        if (!targetEmp.foto_url) {
            throw new Error("Foto Master Anda belum ada di sistem. Minta Owner untuk upload foto Anda.");
        }

        // 🚀 FIX URL DINAMIS: Gunakan .env untuk menembak foto master
        // Biar lancar di laptop dan HP, pastikan VITE_API_BASE_URL di .env sudah pakai IP (http://192.168.xx.xx:8080)
        const masterUrl = `${import.meta.env.VITE_API_BASE_URL}${targetEmp.foto_url}`;
        const imgMaster = await faceapi.fetchImage(masterUrl);
        
        const masterDetections = await faceapi.detectSingleFace(imgMaster, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!masterDetections) throw new Error("AI tidak menemukan wajah di Foto Master. Coba upload foto baru yang lebih jelas.");

        const faceMatcher = new faceapi.FaceMatcher(masterDetections);

        // 2. Deteksi Wajah di Kamera (Video)
        const queryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!queryDetections) throw new Error("Wajah tidak terdeteksi di kamera! Pastikan wajah terlihat jelas dan cukup cahaya.");

        // Bandingkan wajah kamera dengan wajah master
        const bestMatch = faceMatcher.findBestMatch(queryDetections.descriptor);
        
        // Tolerance: 0.5 (Makin kecil makin ketat)
        if (bestMatch.distance > 0.5) {
            throw new Error("Wajah tidak cocok! Verifikasi gagal.");
        }

        // 3. Lolos Verifikasi, Kecilkan Foto untuk dikirim ke Backend
        const canvas = document.createElement('canvas');
        canvas.width = 320; 
        canvas.height = 400;
        const ctx = canvas.getContext('2d');
        ctx.drawImage(videoRef.value, 0, 0, canvas.width, canvas.height);

        // Quality 0.4 biar gak kena 413 Payload Too Large lagi
        const photoBase64 = canvas.toDataURL('image/jpeg', 0.4); 
        
        stopCamera(); // Matikan kamera setelah sukses

        // Kirim ke database lewat Axios instance kita
        await api.post('/attendance', {
            user_id: targetEmp.user_id, 
            jenis: absenTarget.value.jenis,
            foto: photoBase64
        });

        Swal.fire({
            icon: 'success',
            title: `Absen ${absenTarget.value.jenis} Berhasil!`,
            text: 'Selamat bekerja!',
            timer: 2000,
            showConfirmButton: false
        });

        fetchData(); // Refresh tabel riwayat absen

    } catch (error) {
        console.error("DEBUG ABSEN:", error);
        Swal.fire('Gagal', error.message || "Terjadi kesalahan", 'error');
    } finally {
        isSubmitting.value = false;
    }
};

const downloadLaporan = async () => {
    let bulan, tahun;
    filterMode.value === 'harian' ? [tahun, bulan] = tanggalDipilih.value.split('-').slice(0, 2) : [tahun, bulan] = bulanDipilih.value.split('-');
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/attendance/export?bulan=${bulan}&tahun=${tahun}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url; a.download = `ABSENSI_${bulan}_${tahun}.csv`; a.click();
    } catch (error) { Swal.fire('Gagal!', 'Gagal download laporan.', 'error'); }
};
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-8 max-w-7xl mx-auto font-sans">
            <div class="bg-blue-800 rounded-3xl p-6 md:p-8 mb-8 text-white shadow-xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-blue-900">
                <div class="absolute -right-10 -top-20 opacity-10 text-[200px] font-black italic pointer-events-none">⏱️</div>
                <div class="relative z-10 text-center md:text-left mb-6 md:mb-0">
                    <h1 class="text-3xl font-black tracking-tight mb-2 uppercase">Presensi Karyawan</h1>
                    <p class="text-blue-200 font-medium text-sm italic">Absen Tepat Waktu ya teman-teman !</p>
                </div>
                <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner">
    <div class="text-[10px] font-black text-blue-200 uppercase tracking-[0.2em] mb-1">Status Shift Anda</div>
    <div class="font-mono text-2xl font-black text-yellow-300">
        {{ cekAbsenMasukKaryawan(currentUser.user_id) ? (cekAbsenPulangKaryawan(currentUser.user_id) ? '🏁 SELESAI' : '🟢 AKTIF') : '⚪ BELUM' }}
    </div>
</div>
            </div>

            <h2 class="text-lg font-black text-gray-800 mb-4 flex items-center gap-2 uppercase tracking-tight">👤 Panel Absensi Karyawan</h2>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-10">
                <div v-for="user in karyawan" :key="user.id"
                     class="bg-white rounded-3xl p-5 shadow-sm border-2 transition-all duration-300"
                     :class="user.id == currentUser.user_id ? 'border-blue-500 shadow-blue-100 ring-4 ring-blue-50' : 'border-gray-100 opacity-60'">

                    <div class="flex items-center gap-4 mb-5 pb-4 border-b border-gray-50">
                        <div class="w-14 h-14 rounded-full flex items-center justify-center shrink-0 border-2 border-white shadow-md bg-blue-600 text-white font-black text-lg uppercase">{{ user.name.substring(0, 2) }}</div>
                        <div>
                            <h3 class="font-black text-lg leading-tight text-gray-800 uppercase">{{ user.name }}</h3>
                            <div class="text-[10px] font-black px-2 py-0.5 mt-1 rounded bg-blue-100 text-blue-700 uppercase tracking-widest">{{ user.role }} • {{ user.nik || 'ID:'+user.id }}</div>
                        </div>
                    </div>

                    <div class="flex gap-3">
    <button v-if="!cekAbsenMasukKaryawan(user.id)" @click="mulaiAbsen(user.id, user.name, 'Masuk')"
        :disabled="user.id != currentUser.user_id"
        class="flex-1 py-3.5 rounded-2xl font-black text-sm transition-all border-2 flex items-center justify-center gap-2 bg-green-50 text-green-700 border-green-200 hover:bg-green-600 hover:text-white disabled:opacity-30">
        📸 MASUK
    </button>
    
    <button v-if="cekAbsenMasukKaryawan(user.id) && !cekAbsenPulangKaryawan(user.id)" @click="mulaiAbsen(user.id, user.name, 'Pulang')"
        :disabled="user.id != currentUser.user_id"
        class="flex-1 py-3.5 rounded-2xl font-black text-sm transition-all border-2 flex items-center justify-center gap-2 bg-orange-50 text-orange-700 border-orange-200 hover:bg-orange-600 hover:text-white disabled:opacity-30">
        🏁 PULANG
    </button>

    <div v-if="cekAbsenPulangKaryawan(user.id)" class="w-full py-3.5 text-center bg-gray-100 text-gray-400 font-black rounded-2xl text-xs uppercase tracking-widest border-2 border-dashed border-gray-200">
        TUGAS SELESAI HARI INI
    </div>
</div>
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
                <div class="p-5 md:p-6 border-b border-gray-100 flex flex-col md:flex-row md:justify-between md:items-center bg-gray-50 gap-4">
                    <h3 class="font-black text-gray-800 text-lg flex items-center gap-2 uppercase tracking-tight">📋 Log Riwayat</h3>
                    <div class="flex flex-wrap items-center gap-3">
                        <select v-model="filterMode" class="px-4 py-2 bg-white border-2 border-gray-200 rounded-xl text-sm font-black text-gray-700 outline-none">
                            <option value="harian">HARIAN</option><option value="bulanan">BULANAN</option>
                        </select>
                        <input v-if="filterMode === 'harian'" type="date" v-model="tanggalDipilih" class="px-4 py-2 border-2 border-gray-200 rounded-xl text-sm font-black text-gray-700">
                        <input v-else type="month" v-model="bulanDipilih" class="px-4 py-2 border-2 border-gray-200 rounded-xl text-sm font-black text-gray-700">
                        <button v-if="currentUser.role === 'owner'" @click="downloadLaporan" class="bg-emerald-600 text-white px-5 py-2 rounded-xl font-black text-xs uppercase shadow-md">Ekspor CSV</button>
                    </div>
                </div>
                <div class="overflow-x-auto">
                    <table class="w-full text-left whitespace-nowrap">
    <thead class="bg-white border-b border-gray-100">
        <tr class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">
            <th class="px-6 py-5">Karyawan</th>
            <th v-if="filterMode === 'bulanan'" 
    @click="toggleSortTanggal" 
    class="px-6 py-5 text-center cursor-pointer select-none hover:bg-slate-50 transition-colors group">
    <div class="flex items-center justify-center gap-1">
        <span>Tanggal</span>
        <span class="text-[11px] text-blue-600 transition-transform">
            {{ urutanTanggalTerbaru ? '🔽' : '🔼' }}
        </span>
    </div>
</th>
            <th class="px-6 py-5 text-center">Shift</th>
            <th class="px-6 py-5 text-center">Foto Masuk</th>
            <th class="px-6 py-5 text-center">Jam Masuk</th>
            <th class="px-6 py-5 text-center">Foto Pulang</th>
            <th class="px-6 py-5 text-center">Jam Pulang</th>
            <th class="px-6 py-5 text-center">Status</th>
        </tr>
    </thead>
    <tbody class="divide-y divide-gray-50">
        <tr v-if="isLoading">
            <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-12 text-center text-gray-400 font-bold uppercase animate-pulse">
                Sinkronisasi...
            </td>
        </tr>
        <tr v-else-if="riwayat.length === 0">
            <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-12 text-center text-gray-400 font-bold italic">
                Data Kosong.
            </td>
        </tr>
        <tr v-else v-for="log in riwayat" :key="log.id" class="hover:bg-blue-50/30 transition-colors group">
            <td class="px-6 py-5">
                <div class="font-black text-gray-800 uppercase text-sm">{{ log.User?.name }}</div>
                <div class="text-[9px] text-gray-400 font-bold tracking-widest">
                    {{ log.User?.nik || 'ID: '+log.user_id }}
                </div>
            </td>
            
            <td v-if="filterMode === 'bulanan'" class="px-6 py-5 text-center font-mono text-xs font-black text-slate-600">
                {{ log.tanggal ? log.tanggal.substring(0, 10) : '-' }}
            </td>

            <td class="px-6 py-5 text-center font-black text-xs text-indigo-600 uppercase">
                {{ log.shift }}
            </td>

            <td class="px-6 py-5 text-center">
                <div v-if="log.foto_masuk" @click="lihatFoto(log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-14 h-14 rounded-xl mx-auto border-4 border-white shadow-sm overflow-hidden cursor-zoom-in">
                    <img :src="log.foto_masuk" class="w-full h-full object-cover">
                </div>
                <span v-else class="text-gray-200">❌</span>
            </td>

            <td class="px-6 py-5 text-center">
                <span v-if="log.jam_masuk" class="bg-green-100 text-green-700 font-black px-3 py-1.5 rounded-lg text-xs">
                    {{ log.jam_masuk }}
                </span>
                <span v-else class="text-gray-300">-</span>
            </td>

            <td class="px-6 py-5 text-center">
                <div v-if="log.foto_pulang" @click="lihatFoto(log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-14 h-14 rounded-xl mx-auto border-4 border-white shadow-sm overflow-hidden cursor-zoom-in">
                    <img :src="log.foto_pulang" class="w-full h-full object-cover">
                </div>
                <span v-else class="text-gray-200">❌</span>
            </td>

            <td class="px-6 py-5 text-center">
                <span v-if="log.jam_pulang" class="bg-orange-100 text-orange-700 font-black px-3 py-1.5 rounded-lg text-xs">
                    {{ log.jam_pulang }}
                </span>
                <span v-else class="text-gray-300">-</span>
            </td>
            
            <td class="px-6 py-5 text-center">
    <span v-if="log.status === 'Hadir'" class="bg-green-100 text-green-800 font-black px-3 py-1.5 rounded-full text-[10px] uppercase tracking-wider shadow-sm">
        🟢 HADIR
    </span>
    <span v-else-if="log.status === 'Lupa Absen Pulang'" class="bg-amber-100 text-amber-800 font-black px-3 py-1.5 rounded-full text-[10px] uppercase tracking-wider shadow-sm">
        ⚠️ LUPA PULANG
    </span>
    <span v-else-if="log.status === 'Mangkir'" class="bg-red-100 text-red-800 font-black px-3 py-1.5 rounded-full text-[10px] uppercase tracking-wider shadow-sm">
        🚨 MANGKIR
    </span>
    <span v-else-if="log.status === 'Libur (OFF)'" class="bg-slate-100 text-slate-500 font-black px-3 py-1.5 rounded-full text-[10px] uppercase tracking-wider shadow-sm border border-dashed border-slate-200">
        ⚪ LIBUR (OFF)
    </span>
    <span v-else-if="log.status === 'Belum Absen'" class="bg-blue-100 text-blue-800 font-black px-3 py-1.5 rounded-full text-[10px] uppercase tracking-wider shadow-sm animate-pulse">
        🔵 BELUM ABSEN
    </span>
    <span v-else class="bg-gray-100 text-gray-600 font-black px-3 py-1.5 rounded-full text-[10px] uppercase tracking-wider shadow-sm">
        {{ log.status }}
    </span>
</td>
        </tr>
    </tbody>
</table>
                </div>
            </div>
        </div>

        <div v-if="showCameraModal" class="fixed inset-0 bg-gray-900/95 z-[100] flex items-center justify-center p-4 backdrop-blur-md">
            <div class="bg-white rounded-[40px] overflow-hidden shadow-2xl w-full max-w-md border-8 border-white">
                <div class="p-6 border-b border-gray-50 flex justify-between items-center bg-white">
                    <div>
                        <h3 class="font-black text-gray-800 uppercase tracking-tighter text-xl">Selfie Absen</h3>
                        <p class="text-xs text-blue-600 font-black uppercase mt-1">{{ absenTarget.jenis }} • {{ absenTarget.nama }}</p>
                    </div>
                    <button @click="stopCamera" class="w-10 h-10 rounded-full bg-gray-100 text-gray-400 hover:text-red-600 transition-all font-black">✕</button>
                </div>
                <div class="relative bg-black w-full aspect-[3/4] flex items-center justify-center overflow-hidden">
                    <div v-if="isAiLoading" class="absolute inset-0 bg-slate-900/90 backdrop-blur-md flex flex-col items-center justify-center z-50 text-white">
        <div class="w-10 h-10 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="font-black text-[10px] uppercase tracking-[0.2em]">AI Engine Pemanasan...</p>
    </div>
                    <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover transform scale-x-[-1]"></video>
                    <div class="absolute bottom-6 left-6 text-yellow-400 font-mono text-xs drop-shadow-[0_2px_2px_rgba(0,0,0,0.8)]">
                        <div class="font-black">{{ absenTarget.nama.toUpperCase() }}</div>
                        <div class="text-lg font-black tracking-widest">{{ currentTime }}</div>
                        <div class="font-bold opacity-80">{{ new Date().toLocaleDateString('id-ID') }}</div>
                    </div>
                </div>
                <div class="p-6 bg-white">
                    <button @click="jepretDanKirim" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-5 rounded-[24px] font-black text-xl shadow-2xl transition-all disabled:opacity-50">
                        {{ isSubmitting ? 'MENGIRIM...' : '◎ JEPRET & ABSEN' }}
                    </button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }
</style>