import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import * as faceapi from 'face-api.js';

export function useAbsensi() {
    const currentTime = ref(''), karyawan = ref([]), riwayat = ref([]), urutanTanggalTerbaru = ref(true);
    const isLoading = ref(true), filterMode = ref('harian'), bulanDipilih = ref(new Date().toISOString().slice(0, 7)), isAiLoading = ref(true);
    const showCameraModal = ref(false), videoRef = ref(null), stream = ref(null), absenTarget = ref({ id: null, nama: '', jenis: '' }), isSubmitting = ref(false);
    const today = new Date().toLocaleDateString('en-CA'), tanggalDipilih = ref(today);
    let timer;

    const API_BASE = import.meta.env.VITE_API_BASE_URL || '';
    const getCleanUrl = (url) => url ? (url.startsWith('http') ? url : `${API_BASE}${url.startsWith('/') ? url : `/${url}`}`) : null;
    const updateTime = () => { currentTime.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' }).replace(/:/g, '.'); };
    const toggleSortTanggal = () => { urutanTanggalTerbaru.value = !urutanTanggalTerbaru.value; riwayat.value.sort((a, b) => urutanTanggalTerbaru.value ? b.tanggal.localeCompare(a.tanggal) : a.tanggal.localeCompare(b.tanggal)); };
    
    const getPayloadFromToken = () => {
        const t = localStorage.getItem('token'); if (!t) return { user_id: 0, role: '', name: '' };
        try { return { user_id: JSON.parse(atob(t.split('.')[1])).user_id, role: (localStorage.getItem('role') || 'kasir').toLowerCase(), name: localStorage.getItem('name') || 'User' }; } 
        catch (e) { return { user_id: 0, role: '', name: '' }; }
    };
    
    const currentUser = ref(getPayloadFromToken());
    const setVideoRef = (el) => { videoRef.value = el; };
    const me = computed(() => karyawan.value.find(k => Number(k.id || k.user_id) === Number(currentUser.value.user_id)) || {});

    const fetchData = async () => {
        isLoading.value = true;
        try {
            const resKar = await api.get('/retail/employees'), allEmp = resKar.data.data || [], staff = allEmp.filter(e => e.role !== 'owner');
            karyawan.value = currentUser.value.role === 'owner' ? (staff.length === 0 ? allEmp : staff) : staff;

            const [resAbs, resSch] = await Promise.all([api.get('/retail/attendance', { params: { tanggal: today } }), api.get('/retail/schedules', { params: { start_date: today, end_date: today } })]);
            const absHrIni = resAbs.data.data || [], schHrIni = resSch.data.data || [];

            karyawan.value = karyawan.value.map(emp => {
                const k = emp.id || emp.user_id, s = schHrIni.find(x => Number(x.user_id) === Number(k)), a = absHrIni.find(x => Number(x.user_id) === Number(k));
                return { ...emp, shift_hari_ini: s ? s.shift_type.replace(/ \(Approved\)|\ \(Pending\)/g, '') : 'Belum Set', sudah_masuk: !!(a && a.jam_masuk), sudah_pulang: !!(a && a.jam_pulang && a.jam_pulang !== "") };
            });

            const [thn, bln] = bulanDipilih.value.split('-'), tglAwal = `${bulanDipilih.value}-01`, tglAkhir = `${bulanDipilih.value}-${new Date(thn, bln, 0).getDate()}`;
            const params = filterMode.value === 'harian' ? { tanggal: tanggalDipilih.value } : { bulan: bln, tahun: thn };
            const [resRiwayat, resSched] = await Promise.all([api.get('/retail/attendance', { params }), api.get('/retail/schedules', { params: { start_date: filterMode.value === 'harian' ? tanggalDipilih.value : tglAwal, end_date: filterMode.value === 'harian' ? tanggalDipilih.value : tglAkhir } })]);
            const dAbs = resRiwayat.data.data || [], dSch = resSched.data.data || [], gabung = [];

            if (filterMode.value === 'harian') {
                karyawan.value.forEach(emp => {
                    const k = emp.id || emp.user_id, j = dSch.find(s => Number(s.user_id) === Number(k) && s.tanggal.substring(0, 10) === tanggalDipilih.value), a = dAbs.find(x => Number(x.user_id) === Number(k) && x.tanggal.substring(0, 10) === tanggalDipilih.value);
                    if (j && j.shift_type === 'OFF') return;
                    gabung.push({ id: a?.id || `temp-h-${k}`, user_id: k, tanggal: tanggalDipilih.value, User: emp, shift: j ? j.shift_type.replace(/ \(Approved\)|\ \(Pending\)/g, '') : 'Belum Set', foto_masuk: a?.foto_masuk || null, jam_masuk: a?.jam_masuk || null, foto_pulang: a?.foto_pulang || null, jam_pulang: a?.jam_pulang || null, status: a?.status || (tanggalDipilih.value < today ? 'Mangkir' : 'Belum Absen') });
                });
            } else {
                const jmlHr = new Date(thn, bln, 0).getDate();
                karyawan.value.forEach(emp => {
                    const k = emp.id || emp.user_id;
                    for (let h = 1; h <= jmlHr; h++) {
                        const tglStr = `${bulanDipilih.value}-${String(h).padStart(2, '0')}`, a = dAbs.find(x => Number(x.user_id) === Number(k) && x.tanggal.substring(0, 10) === tglStr), j = dSch.find(s => Number(s.user_id) === Number(k) && s.tanggal.substring(0, 10) === tglStr);
                        const sh = j ? j.shift_type.replace(/ \(Approved\)|\ \(Pending\)/g, '') : 'Belum Set';
                        let stat = a ? a.status : (sh === 'OFF' ? 'Libur (OFF)' : (emp.role === 'owner' ? 'Owner' : (tglStr < today ? 'Mangkir' : 'Belum Absen')));
                        if (tglStr <= today || a) gabung.push({ id: a?.id || `temp-b-${k}-${tglStr}`, user_id: k, tanggal: tglStr, User: emp, shift: sh, foto_masuk: a?.foto_masuk || null, jam_masuk: a?.jam_masuk || null, foto_pulang: a?.foto_pulang || null, jam_pulang: a?.jam_pulang || null, status: stat });
                    }
                });
            }
            riwayat.value = [...gabung].sort((a, b) => b.tanggal.localeCompare(a.tanggal));
        } catch (err) { console.error("Gagal sinkron log:", err); } finally { isLoading.value = false; }
    };

    const lihatFoto = (url, nama, tipe, jam) => {
        if (!url) return;
        Swal.fire({ title: `Foto ${tipe}`, html: `<div class="text-left font-bold text-sm">👤 Karyawan: ${nama}<br>⏰ Jam: ${jam}</div>`, imageUrl: getCleanUrl(url), imageAlt: 'Foto Absensi', confirmButtonText: 'Tutup', confirmButtonColor: '#2563eb', customClass: { image: 'rounded-2xl border-4 border-gray-100 shadow-lg' } });
    };

    const mulaiAbsen = async (id, nama, jenis) => {
        absenTarget.value = { id, nama, jenis }; showCameraModal.value = true;
        try { stream.value = await navigator.mediaDevices.getUserMedia({ video: { facingMode: "user", width: 1280, height: 720 } }); await nextTick(); if (videoRef.value) videoRef.value.srcObject = stream.value; } 
        catch (err) { Swal.fire('Error', 'Nyalakan izin kamera ya!', 'error'); showCameraModal.value = false; }
    };

    const stopCamera = () => { if (stream.value) stream.value.getTracks().forEach(t => t.stop()); showCameraModal.value = false; };

    const jepretDanKirim = async () => {
        if (isAiLoading.value) return Swal.fire('Tunggu', 'Sistem AI sedang memuat model...', 'warning');
        isSubmitting.value = true;
        try {
            const tgtEmp = (await api.get('/me')).data;
            if (!tgtEmp.biometric_url) throw new Error("Data Biometrik Wajah belum terdaftar.");
            const masterDetections = await faceapi.detectSingleFace(await faceapi.fetchImage(getCleanUrl(tgtEmp.biometric_url)), new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptor();
            if (!masterDetections) throw new Error("Wajah di sistem tidak terbaca.");
            
            const qryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptor();
            if (!qryDetections) throw new Error("Wajah tidak terdeteksi di kamera.");
            if (new faceapi.FaceMatcher(masterDetections).findBestMatch(qryDetections.descriptor).distance > 0.5) throw new Error("Wajah tidak cocok.");

            const cvs = document.createElement('canvas'); cvs.width = 320; cvs.height = 400;
            cvs.getContext('2d').drawImage(videoRef.value, 0, 0, cvs.width, cvs.height);
            const photoBase64 = cvs.toDataURL('image/jpeg', 0.4); stopCamera();

            await api.post('/retail/attendance', { jenis: absenTarget.value.jenis, foto: photoBase64 });
            Swal.fire({ icon: 'success', title: `Absen ${absenTarget.value.jenis} Berhasil!`, timer: 2000, showConfirmButton: false });
            fetchData();
        } catch (err) { Swal.fire('Gagal Verifikasi', err.response?.data?.error || err.message, 'error'); } 
        finally { isSubmitting.value = false; }
    };

    // 🚀 FIXED: Ganti fetch mentah jadi api.get Axios (Anti Bypass Token Expired & No Hardcode)
    const downloadLaporan = async () => {
        let bln, thn; filterMode.value === 'harian' ? [thn, bln] = tanggalDipilih.value.split('-').slice(0, 2) : [thn, bln] = bulanDipilih.value.split('-');
        try {
            Swal.fire({ title: 'Menyiapkan Laporan...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });
            const res = await api.get('/retail/attendance/export', { params: { bulan: bln, tahun: thn }, responseType: 'blob' });
            const url = window.URL.createObjectURL(new Blob([res.data])), a = document.createElement('a');
            a.href = url; a.download = `ABSENSI_${bln}_${thn}.csv`; document.body.appendChild(a); a.click(); a.remove(); URL.revokeObjectURL(url); Swal.close();
        } catch (err) { Swal.fire('Gagal!', 'Gagal download laporan.', 'error'); }
    };

    onMounted(async () => {
        updateTime(); timer = setInterval(updateTime, 1000);
        try { await Promise.all([faceapi.nets.tinyFaceDetector.loadFromUri('/models'), faceapi.nets.faceLandmark68Net.loadFromUri('/models'), faceapi.nets.faceRecognitionNet.loadFromUri('/models')]); isAiLoading.value = false; } 
        catch (e) { console.error("AI Load Failed:", e); }
        fetchData();
    });
    onUnmounted(() => { clearInterval(timer); stopCamera(); });
    watch([tanggalDipilih, bulanDipilih, filterMode], () => fetchData());

    return { currentTime, karyawan, riwayat, urutanTanggalTerbaru, isLoading, filterMode, bulanDipilih, isAiLoading, toggleSortTanggal, currentUser, tanggalDipilih, showCameraModal, absenTarget, isSubmitting, me, setVideoRef, fetchData, lihatFoto, mulaiAbsen, stopCamera, jepretDanKirim, downloadLaporan };
}