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
    // 🚀 REVISI AMAN: Validasi tipe data gambar agar Base64 dan URL Cloud tidak bentrok 404
    const getCleanUrl = (url) => {
        if (!url || url === 'null' || url === 'undefined' || url === '') return null;
        
        // Jika format Base64 dari webcam lokal, langsung loloskan tanpa ditambahin alamat API
        if (url.startsWith('data:image/')) return url;
        
        // Jika sudah berupa link utuh eksternal (Supabase cloud)
        if (url.startsWith('http')) return url;
        
        // Jika berupa relative path sisa sistem lokal backend
        return `${API_BASE}${url.startsWith('/') ? url : `/${url}`}`;
    };
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

            // 🚀 SUNTIKAN SAKTI SINKRONISASI: Paksa currentUser.user_id megang public_id biar sinkron sama tombol di UI
            const profilSendiri = allEmp.find(e => String(e.id) === String(currentUser.value.user_id) || String(e.public_id) === String(currentUser.value.user_id));
            if (profilSendiri && profilSendiri.public_id) {
                currentUser.value.user_id = profilSendiri.public_id;
            }

            const [resAbs, resSch] = await Promise.all([api.get('/retail/attendance', { params: { tanggal: today } }), api.get('/retail/schedules', { params: { start_date: today, end_date: today } })]);
            const absHrIni = resAbs.data.data || [], schHrIni = resSch.data.data || [];

            karyawan.value = karyawan.value.map(emp => {
                const kStr = String(emp.public_id || emp.id).trim();
                const empIdLama = String(emp.id).trim();
                
                // FIX: Gunakan string hybrid matching, buang jauh-jauh fungsi Number()
                const s = schHrIni.find(x => String(x.user_id).trim() === kStr || String(x.user_id).trim() === empIdLama);
                const a = absHrIni.find(x => String(x.user_id).trim() === kStr || String(x.user_id).trim() === empIdLama);
                
                return { 
                    ...emp, 
                    shift_hari_ini: s ? s.shift_type.replace(/ \(Approved\)|\ \(Pending\)/g, '') : 'Belum Set', 
                    sudah_masuk: !!(a && a.jam_masuk), 
                    sudah_pulang: !!(a && a.jam_pulang && a.jam_pulang !== "") 
                };
            });

            const [thn, bln] = bulanDipilih.value.split('-'), tglAwal = `${bulanDipilih.value}-01`, tglAkhir = `${bulanDipilih.value}-${new Date(thn, bln, 0).getDate()}`;
            const params = filterMode.value === 'harian' ? { tanggal: tanggalDipilih.value } : { bulan: bln, tahun: thn };
            const [resRiwayat, resSched] = await Promise.all([api.get('/retail/attendance', { params }), api.get('/retail/schedules', { params: { start_date: filterMode.value === 'harian' ? tanggalDipilih.value : tglAwal, end_date: filterMode.value === 'harian' ? tanggalDipilih.value : tglAkhir } })]);
            const dAbs = resRiwayat.data.data || [], dSch = resSched.data.data || [], gabung = [];

            if (filterMode.value === 'harian') {
                karyawan.value.forEach(emp => {
                    const kStr = String(emp.public_id || emp.id).trim();
                    const empIdLama = String(emp.id).trim();
                    
                    // FIX: Gunakan string hybrid matching
                    const j = dSch.find(s => (String(s.user_id).trim() === kStr || String(s.user_id).trim() === empIdLama) && s.tanggal.substring(0, 10) === tanggalDipilih.value);
                    const a = dAbs.find(x => (String(x.user_id).trim() === kStr || String(x.user_id).trim() === empIdLama) && x.tanggal.substring(0, 10) === tanggalDipilih.value);
                    
                    if (j && j.shift_type === 'OFF') return;
                    gabung.push({ id: a?.id || `temp-h-${kStr}`, user_id: kStr, tanggal: tanggalDipilih.value, User: emp, shift: j ? j.shift_type.replace(/ \(Approved\)|\ \(Pending\)/g, '') : 'Belum Set', foto_masuk: a?.foto_masuk || null, jam_masuk: a?.jam_masuk || null, foto_pulang: a?.foto_pulang || null, jam_pulang: a?.jam_pulang || null, status: a?.status || (tanggalDipilih.value < today ? 'Mangkir' : 'Belum Absen') });
                });
            } else {
                const jmlHr = new Date(thn, bln, 0).getDate();
                karyawan.value.forEach(emp => {
                    const kStr = String(emp.public_id || emp.id).trim();
                    const empIdLama = String(emp.id).trim();
                    
                    for (let h = 1; h <= jmlHr; h++) {
                        const tglStr = `${bulanDipilih.value}-${String(h).padStart(2, '0')}`;
                        
                        // FIX: Gunakan string hybrid matching
                        const a = dAbs.find(x => (String(x.user_id).trim() === kStr || String(x.user_id).trim() === empIdLama) && x.tanggal.substring(0, 10) === tglStr);
                        const j = dSch.find(s => (String(s.user_id).trim() === kStr || String(s.user_id).trim() === empIdLama) && s.tanggal.substring(0, 10) === tglStr);
                        
                        const sh = j ? j.shift_type.replace(/ \(Approved\)|\ \(Pending\)/g, '') : 'Belum Set';
                        let stat = a ? a.status : (sh === 'OFF' ? 'Libur (OFF)' : (emp.role === 'owner' ? 'Owner' : (tglStr < today ? 'Mangkir' : 'Belum Absen')));
                        if (tglStr <= today || a) gabung.push({ id: a?.id || `temp-b-${kStr}-${tglStr}`, user_id: kStr, tanggal: tglStr, User: emp, shift: sh, foto_masuk: a?.foto_masuk || null, jam_masuk: a?.jam_masuk || null, foto_pulang: a?.foto_pulang || null, jam_pulang: a?.jam_pulang || null, status: stat });
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
            // Mengambil profil diri sendiri yang tervalidasi dari server
            const resMe = await api.get('/me');
            const targetEmp = resMe.data;
            
            if (!targetEmp.biometric_url) throw new Error("Data Biometrik Wajah belum terdaftar.");

            const masterUrl = getCleanUrl(targetEmp.biometric_url);
            
            const imgMaster = await faceapi.fetchImage(masterUrl);
            const masterDetections = await faceapi.detectSingleFace(imgMaster, new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptor();
            if (!masterDetections) throw new Error("Wajah di foto sistem tidak terbaca.");

            const faceMatcher = new faceapi.FaceMatcher(masterDetections);
            const queryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks().withFaceDescriptor();
            if (!queryDetections) throw new Error("Wajah tidak terdeteksi di kamera.");

            const bestMatch = faceMatcher.findBestMatch(queryDetections.descriptor);
            if (bestMatch.distance > 0.5) throw new Error("Wajah tidak cocok dengan database.");

            const canvas = document.createElement('canvas'); 
            canvas.width = 320; 
            canvas.height = 400;
            const ctx = canvas.getContext('2d'); 
            
            // Menggambar tangkapan kamera (dibalik agar tidak seperti cermin)
            ctx.translate(canvas.width, 0);
            ctx.scale(-1, 1);
            ctx.drawImage(videoRef.value, 0, 0, canvas.width, canvas.height);
            
            // Reset transformasi sebelum menggambar teks agar tulisan tidak terbalik
            ctx.setTransform(1, 0, 0, 1, 0, 0);

            // MEMBAKAR WATERMARK ANTI-SPOOFING KE DALAM GAMBAR
            ctx.fillStyle = "rgba(0, 0, 0, 0.6)"; // Background hitam transparan untuk teks
            ctx.fillRect(0, canvas.height - 60, canvas.width, 60);

            ctx.font = "bold 14px monospace";
            ctx.fillStyle = "#fbbf24"; // Warna kuning amber
            ctx.fillText(`ID: ${targetEmp.public_id || targetEmp.user_id} - ${absenTarget.value.jenis.toUpperCase()}`, 10, canvas.height - 35);
            
            ctx.font = "12px monospace";
            ctx.fillStyle = "#ffffff";
            const waktuCap = new Date().toLocaleString('id-ID');
            ctx.fillText(`WAKTU: ${waktuCap}`, 10, canvas.height - 15);

            const photoBase64 = canvas.toDataURL('image/jpeg', 0.5);
            
            stopCamera();

            // SUNTIKAN KEAMANAN: Hapus user_id dari payload karena backend harusnya mengambil 
            // user_id dari JWT token yang sudah terautentikasi, bukan dari client-side request.
            await api.post('/retail/attendance', { 
                jenis: absenTarget.value.jenis, 
                foto: photoBase64 
            });

            Swal.fire({ 
                icon: 'success', 
                title: `Absen ${absenTarget.value.jenis} Berhasil!`, 
                timer: 2000, 
                showConfirmButton: false 
            });
            fetchData(); 
        } catch (error) {
            console.error("DEBUG ABSEN ERROR:", error);
            let errorMsg = error.response?.data?.error || error.message || "Kesalahan verifikasi.";
            Swal.fire('Gagal Verifikasi', errorMsg, 'error');
        } finally { 
            isSubmitting.value = false; 
        }
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