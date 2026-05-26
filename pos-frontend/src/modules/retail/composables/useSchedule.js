import { ref, computed, onMounted } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useSchedule() {
    // --- STATE DATA ---
    const listKaryawan = ref([]);
    const listJadwal = ref([]);
    const isSaving = ref(false);
    const isLoading = ref(true);

    // --- STATE PENCARIAN KARYAWAN ---
    const searchQuery = ref('');

    // --- STATE USER & ROLE ---
    const getPayloadFromToken = () => {
        const token = localStorage.getItem('token');
        const role = localStorage.getItem('role') || 'staff';
        const name = localStorage.getItem('name') || 'User';
        if (!token) return { user_id: 0, role: '', name: '' };
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            return { user_id: payload.user_id, role: role.toLowerCase(), name: name };
        } catch (e) { return { user_id: 0, role: '', name: '' }; }
    };
    const currentUser = ref(getPayloadFromToken());

    // --- LOGIKA MINGGUAN / SPLIT-WEEK TSM ---
    const getJadwalDates = () => {
        const dates = [];
        const sekarang = new Date();
        const hariIni = sekarang.getDay();
        
        const selisihKeSenin = hariIni === 0 ? -6 : 1 - hariIni;
        const seninMingguIni = new Date(sekarang);
        seninMingguIni.setDate(sekarang.getDate() + selisihKeSenin);

        for (let i = 0; i < 7; i++) {
            const d = new Date(seninMingguIni);
            d.setDate(seninMingguIni.getDate() + i);
            
            const tglString = d.toLocaleDateString('en-CA');
            const namaHari = d.toLocaleDateString('id-ID', { weekday: 'short' }); 
            
            dates.push({ tanggal: tglString, hari: namaHari, tglAngka: d.getDate() });
        }
        return dates;
    };

    const mingguJadwal = ref(getJadwalDates());
    const startDate = computed(() => mingguJadwal.value[0]?.tanggal);
    const endDate = computed(() => mingguJadwal.value[6]?.tanggal);

    // Form State untuk Karyawan Input
    const formJadwal = ref({});

    // --- TARIK DATA DARI BACKEND ---
    const fetchData = async () => {
        isLoading.value = true;
        try {
            const resEmp = await api.get('/retail/employees');
            const allEmployees = resEmp.data.data || [];
            
            listKaryawan.value = allEmployees.filter(emp => emp.role !== 'owner');

            const resSched = await api.get(`/retail/schedules?start_date=${startDate.value}&end_date=${endDate.value}`);
            listJadwal.value = resSched.data.data || [];

            listKaryawan.value.forEach(emp => {
                const empKey = emp.id || emp.user_id;
                formJadwal.value[empKey] = {};
                
                mingguJadwal.value.forEach(d => {
                    const match = listJadwal.value.find(s => {
                        const sTanggalClean = s.tanggal ? s.tanggal.substring(0, 10) : "";
                        return Number(s.user_id) === Number(empKey) && sTanggalClean === d.tanggal;
                    });
                    
                    formJadwal.value[empKey][d.tanggal] = match ? match.shift_type : 'OFF';
                });
            });
        } catch (error) {
            console.error("Gagal sinkronisasi data TSM:", error);
        } finally {
            isLoading.value = false;
        }
    };

    // --- FITUR PENCARIAN REALTIME ---
    const filteredKaryawan = computed(() => {
        if (!searchQuery.value) return listKaryawan.value;
        const query = searchQuery.value.toLowerCase();
        return listKaryawan.value.filter(emp => 
            emp.name.toLowerCase().includes(query) || 
            (emp.nik && String(emp.nik).toLowerCase().includes(query))
        );
    });

    const canEditSchedule = computed(() => {
        const role = currentUser.value.role;
        return role === 'owner' || role === 'manager' || role === 'supervisor';
    });

    // --- SIMPAN JADWAL TOKO ---
    const handleSaveJadwalBulk = async () => {
        isSaving.value = true;
        try {
            const payloadSchedules = [];

            listKaryawan.value.forEach(emp => {
                const empKey = emp.id || emp.user_id;

                mingguJadwal.value.forEach(d => {
                    let shift = formJadwal.value[empKey]?.[d.tanggal] || 'OFF';
                    
                    if (currentUser.value.role !== 'owner' && shift !== 'OFF' && !shift.includes('(Approved)')) {
                        shift = `${shift} (Pending)`;
                    } 
                    else if (currentUser.value.role === 'owner' && shift !== 'OFF' && !shift.includes('(Approved)')) {
                        shift = `${shift} (Approved)`;
                    }

                    payloadSchedules.push({
                        user_id: Number(empKey),
                        tanggal: d.tanggal,
                        shift_type: shift
                    });
                });
            });

            await api.post('/retail/schedules/bulk', { schedules: payloadSchedules });
            
            Swal.fire({
                icon: 'success',
                title: 'Berhasil Disimpan!',
                text: currentUser.value.role === 'owner' 
                    ? 'Seluruh jadwal toko resmi aktif!' 
                    : 'Jadwal toko berhasil diajukan, menunggu ACC Owner!',
                timer: 2000,
                showConfirmButton: false
            });
            
            fetchData(); 
        } catch (e) {
            Swal.fire('Gagal', 'Terjadi kesalahan sistem saat menyimpan master jadwal', 'error');
        } finally {
            isSaving.value = false
        }
    };

    // --- APPROVAL WORKFLOW (KHUSUS OWNER) ---
    const handleApproval = async (empId, tanggal, action) => {
        try {
            let currentShift = formJadwal.value[empId][tanggal];
            const rawShift = currentShift.replace(' (Pending)', '').replace(' (Approved)', '').replace(' (Rejected)', '');
            
            const updatedShift = action === 'approve' ? `${rawShift} (Approved)` : `${rawShift} (Rejected)`;

            const singlePayload = [{
                user_id: Number(empId),
                tanggal: tanggal,
                shift_type: updatedShift
            }];

            await api.post('/retail/schedules/bulk', { schedules: singlePayload });
            
            Swal.fire({
                toast: true, position: 'top-end', icon: 'success', 
                title: `Jadwal di-${action}!`, showConfirmButton: false, timer: 1500
            });
            
            fetchData();
        } catch (e) {
            Swal.fire('Gagal', 'Gagal memproses approval owner.', 'error');
        }
    };

    onMounted(() => fetchData());

    return {
        listKaryawan, listJadwal, isSaving, isLoading, searchQuery, currentUser,
        mingguJadwal, startDate, endDate, formJadwal, filteredKaryawan,
        canEditSchedule, handleSaveJadwalBulk, handleApproval
    };
}