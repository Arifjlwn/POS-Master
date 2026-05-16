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
const urutanTanggalTerbaru = ref(true);
const isLoading = ref(true);
const filterMode = ref('harian');
const bulanDipilih = ref(new Date().toISOString().slice(0, 7));
const isAiLoading = ref(true);

const toggleSortTanggal = () => {
    urutanTanggalTerbaru.value = !urutanTanggalTerbaru.value;
    if (urutanTanggalTerbaru.value) {
        riwayat.value.sort((a, b) => b.tanggal.localeCompare(a.tanggal));
    } else {
        riwayat.value.sort((a, b) => a.tanggal.localeCompare(b.tanggal));
    }
};

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

const me = computed(() => karyawan.value.find(k => Number(k.id || k.user_id) === Number(currentUser.value.user_id)) || {});
const fetchData = async () => {
    isLoading.value = true;
    try {
        const resKaryawan = await api.get('/employees');
        const allEmp = resKaryawan.data.data || [];
        const staffSaja = allEmp.filter(e => e.role !== 'owner');

        if (currentUser.value.role === 'owner') {
            karyawan.value = staffSaja.length === 0 ? allEmp : staffSaja; 
        } else {
            karyawan.value = staffSaja;
        }

        // 🚀 RAHASIA 2: Tarik data Hari Ini SECARA INDEPENDEN khusus buat Panel Absen di atas!
        // Jadi nggak akan terpengaruh biarpun filter Log Riwayat di bawah diganti-ganti ke bulan kemaren.
        const [resPanelAbsen, resPanelSched] = await Promise.all([
            api.get('/attendance', { params: { tanggal: today } }),
            api.get('/schedules', { params: { start_date: today, end_date: today } })
        ]);
        
        const absenHariIni = resPanelAbsen.data.data || [];
        const jadwalHariIni = resPanelSched.data.data || [];

        // Suntikkan status absen & jadwal HARI INI ke masing-masing kartu karyawan
        karyawan.value = karyawan.value.map(emp => {
            const empKey = emp.id || emp.user_id;
            const schedMatch = jadwalHariIni.find(s => Number(s.user_id) === Number(empKey));
            const absenMatch = absenHariIni.find(a => Number(a.user_id) === Number(empKey));
            
            return {
                ...emp,
                shift_hari_ini: schedMatch ? schedMatch.shift_type.replace(' (Approved)', '').replace(' (Pending)', '') : 'Belum Set',
                sudah_masuk: !!(absenMatch && absenMatch.jam_masuk),
                sudah_pulang: !!(absenMatch && absenMatch.jam_pulang && absenMatch.jam_pulang !== "")
            };
        });

        // --- LANJUTAN UNTUK TABEL LOG RIWAYAT DI BAWAH ---
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

onMounted(async () => {
    updateTime();
    timer = setInterval(updateTime, 1000);
    
    try {
        await Promise.all([
            faceapi.nets.tinyFaceDetector.loadFromUri('/models'),
            faceapi.nets.faceLandmark68Net.loadFromUri('/models'),
            faceapi.nets.faceRecognitionNet.loadFromUri('/models')
        ]);
        isAiLoading.value = false;
        console.log("AI Models Loaded!");
    } catch (e) {
        console.error("Gagal load AI models. Pastikan folder /models ada di public folder Vue.", e);
    }
    
    fetchData();
});

onUnmounted(() => {
    clearInterval(timer);
    stopCamera();
});

watch([tanggalDipilih, bulanDipilih, filterMode], () => fetchData());

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
    if (isAiLoading.value) return Swal.fire('Tunggu', 'Sistem AI sedang memuat model...', 'warning');
    isSubmitting.value = true;
    
    try {
        const resMe = await api.get('/me'); 
        const targetEmp = resMe.data;

        if (!targetEmp.biometric_url) {
            throw new Error("Data Biometrik Wajah belum terdaftar. Minta HRD untuk mendaftarkan wajah Anda.");
        }

        const baseUrl = import.meta.env.VITE_API_BASE_URL.replace(/\/$/, ''); 
        const photoPath = targetEmp.biometric_url.startsWith('/') ? targetEmp.biometric_url : `/${targetEmp.biometric_url}`; 
        const masterUrl = `${baseUrl}${photoPath}`;
        
        const imgMaster = await faceapi.fetchImage(masterUrl);
        const masterDetections = await faceapi.detectSingleFace(imgMaster, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!masterDetections) throw new Error("Wajah di foto sistem tidak terbaca oleh AI. Harap ganti foto profil.");

        const faceMatcher = new faceapi.FaceMatcher(masterDetections);
        const queryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!queryDetections) throw new Error("Wajah tidak terdeteksi di kamera! Pastikan wajah di tengah dan terang.");

        const bestMatch = faceMatcher.findBestMatch(queryDetections.descriptor);
        
        if (bestMatch.distance > 0.5) {
            throw new Error("Wajah tidak cocok dengan database! (Verifikasi Ditolak)");
        }

        const canvas = document.createElement('canvas');
        canvas.width = 320; 
        canvas.height = 400;
        const ctx = canvas.getContext('2d');
        ctx.drawImage(videoRef.value, 0, 0, canvas.width, canvas.height);
        const photoBase64 = canvas.toDataURL('image/jpeg', 0.4); 
        
        stopCamera();

        // 🚀 FIX: Pastikan ID berupa angka mutlak!
        const finalUserId = Number(absenTarget.value.id || targetEmp.user_id);

        await api.post('/attendance', {
            user_id: finalUserId, 
            jenis: absenTarget.value.jenis,
            foto: photoBase64
        });

        Swal.fire({
            icon: 'success',
            title: `Absen ${absenTarget.value.jenis} Berhasil!`,
            text: 'Data telah tercatat di server.',
            timer: 2000,
            showConfirmButton: false
        });

        fetchData(); 

    } catch (error) {
        console.error("DEBUG ABSEN ERROR:", error.response || error);
        
        // 🚀 FIX: Tampilkan pesan error ASLI dari Golang biar nggak bingung!
        let errorMsg = "Terjadi kesalahan pada sistem verifikasi.";
        if (error.response && error.response.data && error.response.data.error) {
            errorMsg = error.response.data.error; // Nangkep error dari Golang
        } else if (error.message) {
            errorMsg = error.message;
        }

        if (errorMsg.includes("Failed to fetch")) {
            errorMsg = "Gagal memuat foto server. Cek koneksi atau setting CORS backend.";
        }
        
        Swal.fire('Gagal Verifikasi', errorMsg, 'error');
    } finally {
        isSubmitting.value = false;
    }
};

const downloadLaporan = async () => {
    let bulan, tahun;
    filterMode.value === 'harian' ? [tahun, bulan] = tanggalDipilih.value.split('-').slice(0, 2) : [tahun, bulan] = bulanDipilih.value.split('-');
    try {
        const token = localStorage.getItem('token');
        Swal.fire({ title: 'Menyiapkan Laporan...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/attendance/export?bulan=${bulan}&tahun=${tahun}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url; a.download = `ABSENSI_${bulan}_${tahun}.csv`; a.click();
        Swal.close();
    } catch (error) { Swal.fire('Gagal!', 'Gagal download laporan.', 'error'); }
};
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-indigo-900 via-blue-900 to-slate-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-indigo-800">
                <div class="absolute -right-10 -top-10 opacity-10 text-white pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-72 h-72" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                </div>
                
                <div class="relative z-10 text-center md:text-left mb-6 md:mb-0 w-full md:w-auto">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Presensi <span class="text-blue-400">Tim</span></h1>
                    <p class="text-indigo-200 font-bold text-[10px] md:text-xs uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 9a2 2 0 0 1-2 2H6l-4 4V4c0-1.1.9-2 2-2h8a2 2 0 0 1 2 2v5Z"/><path d="M18 9h2a2 2 0 0 1 2 2v11l-4-4h-6a2 2 0 0 1-2-2v-1"/></svg>
                        Absen Tepat Waktu Ya!
                    </p>
                </div>
                
                <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner w-full md:w-auto">
                    <div class="text-[9px] md:text-[10px] font-black text-indigo-200 uppercase tracking-[0.2em] mb-1">Status Shift Anda</div>
                    <div class="font-mono text-2xl md:text-3xl font-black flex items-center justify-center gap-2 tracking-tighter"
                         :class="me.shift_hari_ini === 'OFF' ? 'text-slate-400' : (me.sudah_masuk ? (me.sudah_pulang ? 'text-slate-300' : 'text-emerald-400') : 'text-amber-400')">
                        <template v-if="me.shift_hari_ini === 'OFF'">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg>
                            LIBUR (OFF)
                        </template>
                        <template v-else>
                            <template v-if="me.sudah_masuk">
                                <template v-if="me.sudah_pulang">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                                    SELESAI
                                </template>
                                <template v-else>
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8 animate-pulse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                    AKTIF
                                </template>
                            </template>
                            <template v-else>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-8 md:h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                BELUM
                            </template>
                        </template>
                    </div>
                </div>
            </div>

            <div class="flex items-center gap-3 mb-6">
                <div class="w-10 h-10 bg-indigo-100 text-indigo-600 rounded-xl flex items-center justify-center shadow-inner">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                </div>
                <h2 class="text-xl font-black text-slate-800 uppercase tracking-tighter italic">
                    Panel Absensi
                </h2>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4 md:gap-6 mb-10">
                <div v-for="user in karyawan" :key="user.id || user.user_id"
                     class="bg-white rounded-[24px] p-5 shadow-sm border-2 transition-all duration-300 relative overflow-hidden group flex flex-col"
                     :class="(user.id || user.user_id) == currentUser.user_id ? 'border-indigo-500 shadow-indigo-100 ring-4 ring-indigo-50' : 'border-slate-100 opacity-70 hover:opacity-100'">
                    
                    <div class="flex items-center gap-4 mb-4 pb-4 border-b border-slate-100 relative z-10">
                        <div class="w-12 h-12 md:w-14 md:h-14 rounded-[16px] flex items-center justify-center shrink-0 border-2 border-white shadow-md bg-gradient-to-br from-slate-800 to-slate-900 text-white font-black text-lg uppercase">{{ user.name.substring(0, 2) }}</div>
                        <div class="flex-1 min-w-0">
                            <h3 class="font-black text-base md:text-lg leading-tight text-slate-800 uppercase truncate">{{ user.name }}</h3>
                            <div class="text-[9px] font-black px-2 py-1 mt-1.5 rounded bg-slate-100 text-slate-500 uppercase tracking-widest inline-flex items-center gap-1 border border-slate-200">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                                {{ user.role }} • {{ user.nik || 'ID:'+(user.id || user.user_id) }}
                            </div>
                        </div>
                    </div>

                    <div class="flex gap-2 md:gap-3 relative z-10 mt-auto">
                        <div v-if="user.shift_hari_ini === 'OFF'" class="w-full py-3 md:py-3.5 text-center bg-slate-100 text-slate-400 font-black rounded-[16px] text-[10px] md:text-[11px] uppercase tracking-widest border-2 border-slate-200 flex items-center justify-center gap-2 shadow-inner">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg>
                            HARI INI LIBUR (OFF)
                        </div>
                        
                        <template v-else>
                            <button v-if="!user.sudah_masuk" @click="mulaiAbsen(user.id || user.user_id, user.name, 'Masuk')"
                                :disabled="(user.id || user.user_id) != currentUser.user_id"
                                class="flex-1 py-3 md:py-3.5 rounded-[16px] font-black text-[10px] md:text-[11px] uppercase tracking-widest transition-all border-2 flex items-center justify-center gap-2 bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-500 hover:text-white hover:border-emerald-500 disabled:opacity-40 disabled:cursor-not-allowed shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                                Absen Masuk
                            </button>
                            
                            <button v-if="user.sudah_masuk && !user.sudah_pulang" @click="mulaiAbsen(user.id || user.user_id, user.name, 'Pulang')"
                                :disabled="(user.id || user.user_id) != currentUser.user_id"
                                class="flex-1 py-3 md:py-3.5 rounded-[16px] font-black text-[10px] md:text-[11px] uppercase tracking-widest transition-all border-2 flex items-center justify-center gap-2 bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-500 hover:text-white hover:border-amber-500 disabled:opacity-40 disabled:cursor-not-allowed shadow-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                                Absen Pulang
                            </button>

                            <div v-if="user.sudah_pulang" class="w-full py-3 md:py-3.5 text-center bg-slate-50 text-slate-400 font-black rounded-[16px] text-[10px] md:text-[11px] uppercase tracking-widest border-2 border-slate-200 flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                Selesai Hari Ini
                            </div>
                        </template>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
                <div class="p-5 md:p-6 border-b border-slate-100 flex flex-col lg:flex-row lg:justify-between lg:items-center bg-slate-50/50 gap-4">
                    <h3 class="font-black text-slate-800 text-lg flex items-center gap-2 uppercase tracking-tight">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                        Log Riwayat
                    </h3>
                    
                    <div class="flex flex-wrap items-center gap-2 md:gap-3 w-full lg:w-auto">
                        <select v-model="filterMode" class="flex-1 lg:flex-none px-4 py-3 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase tracking-widest shadow-sm">
                            <option value="harian">Harian</option><option value="bulanan">Bulanan</option>
                        </select>
                        <input v-if="filterMode === 'harian'" type="date" v-model="tanggalDipilih" class="flex-1 lg:flex-none px-4 py-3 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase shadow-sm cursor-pointer">
                        <input v-else type="month" v-model="bulanDipilih" class="flex-1 lg:flex-none px-4 py-3 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase shadow-sm cursor-pointer">
                        
                        <button v-if="filterMode === 'bulanan'" @click="toggleSortTanggal" class="lg:hidden w-11 h-11 flex items-center justify-center bg-white border border-slate-200 rounded-xl text-slate-500 shadow-sm active:scale-95">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 transition-transform" :class="!urutanTanggalTerbaru ? 'rotate-180 text-indigo-600' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m3 16 4 4 4-4"/><path d="M7 20V4"/><path d="m21 8-4-4-4 4"/><path d="M17 4v16"/></svg>
                        </button>

                        <button v-if="currentUser.role === 'owner'" @click="downloadLaporan" class="w-full lg:w-auto mt-2 lg:mt-0 bg-slate-900 text-white px-5 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-md flex items-center justify-center gap-2 hover:bg-indigo-600 transition-colors active:scale-95">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="3" y2="15"/></svg>
                            Ekspor Laporan
                        </button>
                    </div>
                </div>
                
                <div class="lg:hidden p-4 bg-slate-50/50">
                    <div v-if="isLoading" class="py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse border-2 border-dashed border-slate-200 rounded-[24px]">
                        Sinkronisasi Data...
                    </div>
                    <div v-else-if="riwayat.length === 0" class="py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50 border-2 border-dashed border-slate-200 rounded-[24px]">
                        Data Kosong
                    </div>
                    <div v-else class="flex flex-col gap-4">
                        <div v-for="log in riwayat" :key="log.id" class="bg-white p-4 rounded-[20px] shadow-sm border border-slate-100 flex flex-col gap-3 relative overflow-hidden">
                            <div class="absolute left-0 top-0 bottom-0 w-1.5" 
                                 :class="{
                                    'bg-emerald-400': log.status === 'Hadir',
                                    'bg-amber-400': log.status === 'Lupa Absen Pulang',
                                    'bg-red-500': log.status === 'Mangkir',
                                    'bg-slate-300': log.status === 'Libur (OFF)' || log.status === 'Belum Absen',
                                    'bg-blue-400': log.status === 'Owner' || (!['Hadir','Lupa Absen Pulang','Mangkir','Libur (OFF)','Belum Absen'].includes(log.status))
                                 }">
                            </div>

                            <div class="flex justify-between items-start pl-2">
                                <div>
                                    <div class="font-black text-slate-800 uppercase text-sm tracking-tight">{{ log.User?.name }}</div>
                                    <div class="text-[9px] text-slate-400 font-bold tracking-widest mt-0.5">{{ log.User?.nik || 'ID: '+log.user_id }}</div>
                                </div>
                                <div class="text-right">
                                    <div v-if="filterMode === 'bulanan'" class="text-[10px] font-black text-slate-600 mb-1">{{ log.tanggal ? log.tanggal.substring(0, 10) : '-' }}</div>
                                    <span class="bg-slate-50 border border-slate-100 px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest text-indigo-600">{{ log.shift }}</span>
                                </div>
                            </div>

                            <div class="border-t border-dashed border-slate-100 my-1"></div>

                            <div class="grid grid-cols-2 gap-3 pl-2">
                                <div>
                                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Masuk</p>
                                    <div class="flex items-center gap-2">
                                        <div v-if="log.foto_masuk" @click="lihatFoto(log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-8 h-8 rounded-lg border border-slate-200 overflow-hidden cursor-zoom-in">
                                            <img :src="log.foto_masuk" class="w-full h-full object-cover">
                                        </div>
                                        <div v-else class="w-8 h-8 rounded-lg border border-slate-100 bg-slate-50 flex items-center justify-center">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242"/><path d="M12 12v9"/><path d="m8 17 4-4 4 4"/></svg>
                                        </div>
                                        <span v-if="log.jam_masuk" class="font-black text-xs text-emerald-600">{{ log.jam_masuk }}</span>
                                        <span v-else class="font-black text-xs text-slate-300">-</span>
                                    </div>
                                </div>
                                <div>
                                    <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest mb-1">Pulang</p>
                                    <div class="flex items-center gap-2">
                                        <div v-if="log.foto_pulang" @click="lihatFoto(log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-8 h-8 rounded-lg border border-slate-200 overflow-hidden cursor-zoom-in">
                                            <img :src="log.foto_pulang" class="w-full h-full object-cover">
                                        </div>
                                        <div v-else class="w-8 h-8 rounded-lg border border-slate-100 bg-slate-50 flex items-center justify-center">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242"/><path d="M12 12v9"/><path d="m8 17 4-4 4 4"/></svg>
                                        </div>
                                        <span v-if="log.jam_pulang" class="font-black text-xs text-amber-600">{{ log.jam_pulang }}</span>
                                        <span v-else class="font-black text-xs text-slate-300">-</span>
                                    </div>
                                </div>
                            </div>
                            
                            <div class="pl-2 mt-1">
                                <span v-if="log.status === 'Hadir'" class="inline-flex items-center gap-1 bg-emerald-50 text-emerald-700 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Hadir</span>
                                <span v-else-if="log.status === 'Lupa Absen Pulang'" class="inline-flex items-center gap-1 bg-amber-50 text-amber-700 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Lupa Pulang</span>
                                <span v-else-if="log.status === 'Mangkir'" class="inline-flex items-center gap-1 bg-red-50 text-red-700 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Mangkir</span>
                                <span v-else-if="log.status === 'Libur (OFF)'" class="inline-flex items-center gap-1 bg-slate-100 text-slate-500 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">OFF</span>
                                <span v-else-if="log.status === 'Belum Absen'" class="inline-flex items-center gap-1 bg-blue-50 text-blue-600 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">Belum Absen</span>
                                <span v-else class="inline-flex items-center gap-1 bg-slate-100 text-slate-600 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest">{{ log.status }}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="hidden lg:block overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left whitespace-nowrap">
                        <thead class="bg-white border-b border-slate-100">
                            <tr class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em]">
                                <th class="px-6 py-5">Karyawan</th>
                                <th v-if="filterMode === 'bulanan'" @click="toggleSortTanggal" class="px-6 py-5 text-center cursor-pointer select-none hover:bg-slate-50 transition-colors group">
                                    <div class="flex items-center justify-center gap-2">
                                        Tanggal
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-indigo-600 transition-transform" :class="!urutanTanggalTerbaru ? 'rotate-180' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg>
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
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">
                                    Sinkronisasi Data...
                                </td>
                            </tr>
                            <tr v-else-if="riwayat.length === 0">
                                <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">
                                    Data Riwayat Kosong
                                </td>
                            </tr>
                            <tr v-else v-for="log in riwayat" :key="log.id" class="hover:bg-indigo-50/30 transition-colors group">
                                <td class="px-6 py-5">
                                    <div class="font-black text-slate-800 uppercase text-xs">{{ log.User?.name }}</div>
                                    <div class="text-[9px] text-slate-400 font-bold tracking-widest mt-0.5">{{ log.User?.nik || 'ID: '+log.user_id }}</div>
                                </td>
                                
                                <td v-if="filterMode === 'bulanan'" class="px-6 py-5 text-center font-mono text-[11px] font-black text-slate-600">
                                    {{ log.tanggal ? log.tanggal.substring(0, 10) : '-' }}
                                </td>

                                <td class="px-6 py-5 text-center font-black text-[9px] text-indigo-600 uppercase tracking-widest">
                                    <span class="bg-slate-50 border border-slate-100 px-3 py-1.5 rounded-lg">{{ log.shift }}</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <div v-if="log.foto_masuk" @click="lihatFoto(log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-12 h-12 rounded-[14px] mx-auto border-2 border-slate-200 shadow-sm overflow-hidden cursor-zoom-in hover:border-indigo-500 transition-colors">
                                        <img :src="log.foto_masuk" class="w-full h-full object-cover">
                                    </div>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <span v-if="log.jam_masuk" class="bg-emerald-50 text-emerald-700 border border-emerald-100 font-black px-3 py-1.5 rounded-lg text-[10px]">
                                        {{ log.jam_masuk }}
                                    </span>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <div v-if="log.foto_pulang" @click="lihatFoto(log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-12 h-12 rounded-[14px] mx-auto border-2 border-slate-200 shadow-sm overflow-hidden cursor-zoom-in hover:border-indigo-500 transition-colors">
                                        <img :src="log.foto_pulang" class="w-full h-full object-cover">
                                    </div>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <span v-if="log.jam_pulang" class="bg-amber-50 text-amber-700 border border-amber-100 font-black px-3 py-1.5 rounded-lg text-[10px]">
                                        {{ log.jam_pulang }}
                                    </span>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>
                                
                                <td class="px-6 py-5 text-center">
                                    <span v-if="log.status === 'Hadir'" class="inline-flex items-center gap-1.5 bg-emerald-50 text-emerald-700 border border-emerald-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                        Hadir
                                    </span>
                                    <span v-else-if="log.status === 'Lupa Absen Pulang'" class="inline-flex items-center gap-1.5 bg-amber-50 text-amber-700 border border-amber-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                                        Lupa Pulang
                                    </span>
                                    <span v-else-if="log.status === 'Mangkir'" class="inline-flex items-center gap-1.5 bg-red-50 text-red-700 border border-red-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                                        Mangkir
                                    </span>
                                    <span v-else-if="log.status === 'Libur (OFF)'" class="inline-flex items-center gap-1.5 bg-slate-50 text-slate-500 border border-slate-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg>
                                        OFF
                                    </span>
                                    <span v-else-if="log.status === 'Belum Absen'" class="inline-flex items-center gap-1.5 bg-blue-50 text-blue-700 border border-blue-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                        Belum Absen
                                    </span>
                                    <span v-else class="bg-slate-100 text-slate-600 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm border border-slate-200">
                                        {{ log.status }}
                                    </span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div v-if="showCameraModal" class="fixed inset-0 bg-slate-950/95 z-[100] flex items-center justify-center p-4 backdrop-blur-xl">
            <div class="bg-white rounded-[32px] md:rounded-[40px] overflow-hidden shadow-2xl w-full max-w-sm border-[8px] md:border-[10px] border-slate-900/5">
                <div class="p-5 md:p-6 border-b border-slate-50 flex justify-between items-center bg-white">
                    <div>
                        <h3 class="font-black text-slate-800 uppercase tracking-tighter text-lg md:text-xl italic">Verifikasi Wajah</h3>
                        <p class="text-[9px] md:text-[10px] text-indigo-600 font-black uppercase tracking-widest mt-0.5">{{ absenTarget.jenis }} • {{ absenTarget.nama }}</p>
                    </div>
                    <button @click="stopCamera" class="w-10 h-10 rounded-[14px] bg-slate-100 text-slate-400 hover:text-rose-500 hover:bg-rose-50 transition-all flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                
                <div class="relative bg-slate-900 w-full aspect-[3/4] flex items-center justify-center overflow-hidden">
                    <div v-if="isAiLoading" class="absolute inset-0 bg-slate-950/90 backdrop-blur-md flex flex-col items-center justify-center z-50 text-white">
                        <div class="w-10 h-10 md:w-12 md:h-12 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mb-4"></div>
                        <p class="font-black text-[9px] md:text-[10px] uppercase tracking-[0.2em] animate-pulse">Menyiapkan Engine AI...</p>
                    </div>
                    <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover transform scale-x-[-1]"></video>
                    
                    <div class="absolute inset-0 border-[12px] md:border-[16px] border-black/20 pointer-events-none"></div>
                    <div class="absolute inset-x-8 inset-y-12 md:inset-x-12 md:inset-y-16 border-2 border-white/30 rounded-[32px] md:rounded-[40px] pointer-events-none">
                        <div class="absolute top-0 left-0 w-8 h-8 border-t-4 border-l-4 border-white/80 rounded-tl-[30px] md:rounded-tl-[38px]"></div>
                        <div class="absolute top-0 right-0 w-8 h-8 border-t-4 border-r-4 border-white/80 rounded-tr-[30px] md:rounded-tr-[38px]"></div>
                        <div class="absolute bottom-0 left-0 w-8 h-8 border-b-4 border-l-4 border-white/80 rounded-bl-[30px] md:rounded-bl-[38px]"></div>
                        <div class="absolute bottom-0 right-0 w-8 h-8 border-b-4 border-r-4 border-white/80 rounded-br-[30px] md:rounded-br-[38px]"></div>
                    </div>

                    <div class="absolute bottom-4 left-4 md:bottom-6 md:left-6 text-amber-400 font-mono drop-shadow-md">
                        <div class="font-black tracking-widest text-[10px] md:text-xs">{{ absenTarget.nama.toUpperCase() }}</div>
                        <div class="text-lg md:text-xl font-black tracking-widest mt-0.5 md:mt-1">{{ currentTime }}</div>
                        <div class="font-bold opacity-80 mt-0.5 text-[9px] md:text-[10px]">{{ new Date().toLocaleDateString('id-ID') }}</div>
                    </div>
                </div>
                
                <div class="p-5 md:p-6 bg-white">
                    <button @click="jepretDanKirim" :disabled="isSubmitting" class="w-full bg-indigo-600 hover:bg-slate-900 text-white py-4 md:py-5 rounded-[20px] md:rounded-[24px] font-black text-xs md:text-sm uppercase tracking-[0.2em] shadow-xl shadow-indigo-200 transition-all disabled:opacity-50 flex items-center justify-center gap-2 md:gap-3 active:scale-95">
                        <template v-if="isSubmitting">
                            <div class="w-4 h-4 md:w-5 md:h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            Verifikasi...
                        </template>
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                            Scan & Kirim
                        </template>
                    </button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

/* Styling Calendar Icon pada Input Date biar warnanya nyatu */
input[type="date"]::-webkit-calendar-picker-indicator,
input[type="month"]::-webkit-calendar-picker-indicator {
    cursor: pointer;
    opacity: 0.6;
    transition: 0.2s;
}
input[type="date"]::-webkit-calendar-picker-indicator:hover,
input[type="month"]::-webkit-calendar-picker-indicator:hover {
    opacity: 1;
}

/* Buat input tanggal di header biru */
.color-scheme-dark::-webkit-calendar-picker-indicator {
    filter: invert(1);
    opacity: 0.8;
}
</style>