import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import * as faceapi from 'face-api.js';

export function useAbsensi() {
    // --- STATE WAKTU REALTIME ---
    const currentTime = ref('');
    let timer;

    const updateTime = () => {
        const now = new Date();
        currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' }).replace(/:/g, '.');
    };

    // --- STATE DATA ---
    const karyawan = ref([]);
    const riwayat = ref([]);
    const urutanTanggalTerbaru = ref(true);
    const isLoading = ref(true);
    const filterMode = ref('harian');
    const bulanDipilih = ref(new Date().toISOString().slice(0, 7));
    const isAiLoading = ref(true);
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    // 🚀 HELPER SAKTI: Menjamin link Cloud Supabase (https://) dan sisa data Lokal aman dikunyah oleh AI & Img tag
    const getCleanUrl = (url) => {
        if (!url) return null;
        if (url.startsWith('http://') || url.startsWith('https://')) return url;
        const cleanPath = url.startsWith('/') ? url : `/${url}`;
        return `${API_BASE_URL}${cleanPath}`;
    };

    const toggleSortTanggal = () => {
        urutanTanggalTerbaru.value = !urutanTanggalTerbaru.value;
        riwayat.value.sort((a, b) => urutanTanggalTerbaru.value ? b.tanggal.localeCompare(a.tanggal) : a.tanggal.localeCompare(b.tanggal));
    };

    const getPayloadFromToken = () => {
        const token = localStorage.getItem('token'), name = localStorage.getItem('name') || 'User', role = localStorage.getItem('role') || 'kasir';
        if (!token) return { user_id: 0, role: '', name: '' };
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return { user_id: payload.user_id, role: role.toLowerCase(), name: name };
        } catch (e) { return { user_id: 0, role: '', name: '' }; }
    };
    const currentUser = ref(getPayloadFromToken());
    const today = new Date().toLocaleDateString('en-CA');
    const tanggalDipilih = ref(today);

    // --- STATE KAMERA ---
    const showCameraModal = ref(false);
    const videoRef = ref(null);
    const stream = ref(null);
    const absenTarget = ref({ id: null, nama: '', jenis: '' });
    const isSubmitting = ref(false);

    const setVideoRef = (el) => { videoRef.value = el; };
    const me = computed(() => karyawan.value.find(k => Number(k.id || k.user_id) === Number(currentUser.value.user_id)) || {});
    
    const fetchData = async () => {
        isLoading.value = true;
        try {
            const resKaryawan = await api.get('/retail/employees');
            const allEmp = resKaryawan.data.data || [], staffSaja = allEmp.filter(e => e.role !== 'owner');
            karyawan.value = currentUser.value.role === 'owner' ? (staffSaja.length === 0 ? allEmp : staffSaja) : staffSaja;

            const [resPanelAbsen, resPanelSched] = await Promise.all([
                api.get('/retail/attendance', { params: { tanggal: today } }),
                api.get('/retail/schedules', { params: { start_date: today, end_date: today } })
            ]);
            
            const absenHariIni = resPanelAbsen.data.data || [], jadwalHariIni = resPanelSched.data.data || [];

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

            const [tahun, bulan] = bulanDipilih.value.split('-'), tanggalAwalBulan = `${bulanDipilih.value}-01`, tanggalAkhirBulan = `${bulanDipilih.value}-${new Date(tahun, bulan, 0).getDate()}`;
            const params = filterMode.value === 'harian' ? { tanggal: tanggalDipilih.value } : { bulan, tahun };

            const [resRiwayat, resSched] = await Promise.all([
                api.get('/retail/attendance', { params }),
                api.get('/retail/schedules', { params: { start_date: filterMode.value === 'harian' ? tanggalDipilih.value : tanggalAwalBulan, end_date: filterMode.value === 'harian' ? tanggalDipilih.value : tanggalAkhirBulan } })
            ]);

            const dataAbsenReal = resRiwayat.data.data || [], dataJadwalReal = resSched.data.data || [], matriksGabungan = [];

            if (filterMode.value === 'harian') {
                karyawan.value.forEach(emp => {
                    const empKey = emp.id || emp.user_id;
                    const jadwalHariIni = dataJadwalReal.find(s => Number(s.user_id) === Number(empKey) && s.tanggal.substring(0, 10) === tanggalDipilih.value);
                    const absenHariIni = dataAbsenReal.find(a => Number(a.user_id) === Number(empKey) && a.tanggal.substring(0, 10) === tanggalDipilih.value);
                    if (jadwalHariIni && jadwalHariIni.shift_type === 'OFF') return;

                    matriksGabungan.push({
                        id: absenHariIni?.id || `temp-harian-${empKey}`, user_id: empKey, tanggal: tanggalDipilih.value, User: emp,
                        shift: jadwalHariIni ? jadwalHariIni.shift_type.replace(' (Approved)','').replace(' (Pending)','') : 'Belum Set',
                        foto_masuk: absenHariIni?.foto_masuk || null, jam_masuk: absenHariIni?.jam_masuk || null,
                        foto_pulang: absenHariIni?.foto_pulang || null, jam_pulang: absenHariIni?.jam_pulang || null,
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
                        
                        if (absenMatch) statusDinamis = absenMatch.status;
                        else if (shiftClean === 'OFF') statusDinamis = 'Libur (OFF)';
                        else if (emp.role === 'owner') statusDinamis = 'Owner';
                        else if (tglLoopStr < today) statusDinamis = 'Mangkir';

                        if (tglLoopStr <= today || absenMatch) {
                            matriksGabungan.push({
                                id: absenMatch?.id || `temp-bulan-${empKey}-${tglLoopStr}`, user_id: empKey, tanggal: tglLoopStr, User: emp, shift: shiftClean,
                                foto_masuk: absenMatch?.foto_masuk || null, jam_masuk: absenMatch?.jam_masuk || null,
                                foto_pulang: absenMatch?.foto_pulang || null, jam_pulang: absenMatch?.jam_pulang || null, status: statusDinamis
                            });
                        }
                    }
                });
            }
            riwayat.value = [...matriksGabungan].sort((a, b) => b.tanggal.localeCompare(a.tanggal));
        } catch (error) { console.error("Gagal sinkronisasi data log:", error); } 
        finally { isLoading.value = false; }
    };

    const lihatFoto = (url, nama, tipe, jam) => {
        if (!url) return;
        Swal.fire({
            title: `Foto ${tipe}`, html: `<div class="text-left font-bold text-sm">👤 Karyawan: ${nama}<br>⏰ Jam: ${jam}</div>`,
            imageUrl: getCleanUrl(url), imageAlt: 'Foto Absensi', confirmButtonText: 'Tutup', confirmButtonColor: '#2563eb',
            customClass: { image: 'rounded-2xl border-4 border-gray-100 shadow-lg' }
        });
    };

    onMounted(async () => {
        updateTime(); timer = setInterval(updateTime, 1000);
        try {
            const MODEL_URL = '/models';
            await faceapi.nets.tinyFaceDetector.loadFromUri(MODEL_URL);
            await faceapi.nets.faceLandmark68Net.loadFromUri(MODEL_URL);
            await faceapi.nets.faceRecognitionNet.loadFromUri(MODEL_URL);
            isAiLoading.value = false;
        } catch (e) { console.error("AI Load Failed. Check network tab for 404s:", e); }
        fetchData();
    });

    onUnmounted(() => { clearInterval(timer); stopCamera(); });
    watch([tanggalDipilih, bulanDipilih, filterMode], () => fetchData());

    const mulaiAbsen = async (id, nama, jenis) => {
        absenTarget.value = { id, nama, jenis }; showCameraModal.value = true;
        try {
            stream.value = await navigator.mediaDevices.getUserMedia({ video: { facingMode: "user", width: 1280, height: 720 } });
            await nextTick();
            if (videoRef.value) videoRef.value.srcObject = stream.value;
        } catch (err) {
            Swal.fire('Kamera Error', 'Nyalakan izin kamera ya!', 'error'); showCameraModal.value = false;
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
            const resMe = await api.get('/me'), targetEmp = resMe.data;
            if (!targetEmp.biometric_url) throw new Error("Data Biometrik Wajah belum terdaftar. Minta HRD untuk mendaftarkan wajah Anda.");

            // 🚀 SUNTIKAN KUNCI: Ambil URL Master dari helper cloud-ready tanpa ngerusak string path
            const masterUrl = getCleanUrl(targetEmp.biometric_url);
            
            const imgMaster = await faceapi.fetchImage(masterUrl);
            const masterDetections = await faceapi.detectSingleFace(imgMaster, new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptor();
            if (!masterDetections) throw new Error("Wajah di foto sistem tidak terbaca oleh AI. Harap ganti foto profil.");

            const faceMatcher = new faceapi.FaceMatcher(masterDetections);
            const queryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptor();
            if (!queryDetections) throw new Error("Wajah tidak terdeteksi di kamera! Pastikan wajah di tengah dan terang.");

            const bestMatch = faceMatcher.findBestMatch(queryDetections.descriptor);
            if (bestMatch.distance > 0.5) throw new Error("Wajah tidak cocok dengan database! (Verifikasi Ditolak)");

            const canvas = document.createElement('canvas'); canvas.width = 320; canvas.height = 400;
            const ctx = canvas.getContext('2d'); ctx.drawImage(videoRef.value, 0, 0, canvas.width, canvas.height);
            const photoBase64 = canvas.toDataURL('image/jpeg', 0.4); 
            
            stopCamera();
            const finalUserId = Number(absenTarget.value.id || targetEmp.user_id);

            await api.post('/retail/attendance', { user_id: finalUserId, jenis: absenTarget.value.jenis, foto: photoBase64 });
            Swal.fire({ icon: 'success', title: `Absen ${absenTarget.value.jenis} Berhasil!`, text: 'Data telah tercatat di server.', timer: 2000, showConfirmButton: false });
            fetchData(); 
        } catch (error) {
            console.error("DEBUG ABSEN ERROR:", error.response || error);
            let errorMsg = error.response?.data?.error || error.message || "Terjadi kesalahan pada sistem verifikasi.";
            if (errorMsg.includes("Failed to fetch")) errorMsg = "Gagal memuat foto server. Cek koneksi atau setting CORS backend.";
            Swal.fire('Gagal Verifikasi', errorMsg, 'error');
        } finally { isSubmitting.value = false; }
    };

    const downloadLaporan = async () => {
        let bulan, tahun;
        filterMode.value === 'harian' ? [tahun, bulan] = tanggalDipilih.value.split('-').slice(0, 2) : [tahun, bulan] = bulanDipilih.value.split('-');
        try {
            const token = localStorage.getItem('token');
            Swal.fire({ title: 'Menyiapkan Laporan...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
            const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/retail/attendance/export?bulan=${bulan}&tahun=${tahun}`, { method: 'GET', headers: { 'Authorization': `Bearer ${token}` } });
            const blob = await response.blob(), url = window.URL.createObjectURL(blob), a = document.createElement('a');
            a.href = url; a.download = `ABSENSI_${bulan}_${tahun}.csv`; a.click(); Swal.close();
        } catch (error) { Swal.fire('Gagal!', 'Gagal download laporan.', 'error'); }
    };

    return {
        currentTime, karyawan, riwayat, urutanTanggalTerbaru, isLoading, filterMode, bulanDipilih, isAiLoading, toggleSortTanggal, currentUser,
        tanggalDipilih, showCameraModal, absenTarget, isSubmitting, me, setVideoRef, fetchData, lihatFoto, mulaiAbsen, stopCamera, jepretDanKirim, downloadLaporan
    };
}