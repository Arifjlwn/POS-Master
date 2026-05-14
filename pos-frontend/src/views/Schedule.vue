<script setup>
import { ref, onMounted, computed } from 'vue';
import Sidebar from '../components/Sidebar.vue';
import api from '../api.js';
import Swal from 'sweetalert2';

// --- STATE DATA ---
const listKaryawan = ref([]);
const listJadwal = ref([]);
const isLoading = ref(false);

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
    const hariIni = rollers();
    
    function rollers() {
        const day = sekarang.getDay();
        return day;
    }
    
    const selisihKeSenin = hariIni === 0 ? -6 : 1 - hariIni;
    const seninMingguIni = new Date(sekarang);
    seninMingguIni.setDate(sekarang.getDate() + selisihKeSenin);

    for (let i = 0; i < 7; i++) {
        const d = new Date(seninMingguIni);
        d.setDate(seninMingguIni.getDate() + i);
        
        const tglString = d.toLocaleDateString('en-CA');
        const namaHari = d.toLocaleDateString('id-ID', { weekday: 'long' });
        
        dates.push({ tanggal: tglString, hari: namaHari });
    }
    return dates;
};

const mingguJadwal = ref(getJadwalDates());
const startDate = computed(() => mingguJadwal.value[0]?.tanggal);
const endDate = computed(() => mingguJadwal.value[6]?.tanggal);

// Form State untuk Karyawan Input
const formJadwal = ref({});

// --- TARIK DATA DARI BACKEND (DIAMANKAN AGAR REAKTIF) ---
const fetchData = async () => {
    isLoading.value = true;
    try {
        const resEmp = await api.get('/employees');
        const allEmployees = resEmp.data.data || [];
        
        listKaryawan.value = allEmployees.filter(emp => emp.role !== 'owner');

        const resSched = await api.get(`/schedules?start_date=${startDate.value}&end_date=${endDate.value}`);
        listJadwal.value = resSched.data.data || [];

        listKaryawan.value.forEach(emp => {
            const empKey = emp.id || emp.user_id;
            formJadwal.value[empKey] = {};
            
            mingguJadwal.value.forEach(d => {
                const match = listJadwal.value.find(s => {
                    // 🚀 AMANKAN STRING TANGGAL: Ambil 10 karakter pertama saja (YYYY-MM-DD)
                    const sTanggalClean = s.tanggal ? s.tanggal.substring(0, 10) : "";
                    
                    return Number(s.user_id) === Number(empKey) && sTanggalClean === d.tanggal;
                });
                
                // Masukkan data asli shift_type dari database jika match ditemukan
                formJadwal.value[empKey][d.tanggal] = match ? match.shift_type : 'OFF';
            });
        });
    } catch (error) {
        console.error("Gagal sinkronisasi data TSM:", error);
    } finally {
        isLoading.value = false;
    }
};

const canEditSchedule = computed(() => {
    const role = currentUser.value.role;
    return role === 'owner' || role === 'manager' || role === 'supervisor';
});

// --- SIMPAN JADWAL TOKO (FIXED REFERENCE ERROR) ---
const handleSaveJadwalBulk = async () => {
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

        // 🚀 CETAK CONSOLE LOG DI DALAM BLOK AMAN (Sebelum Tembak API)
        console.log("PAYLOAD YANG DIKIRIM KE GO:", payloadSchedules);

        await api.post('/schedules/bulk', { schedules: payloadSchedules });
        
        Swal.fire({
            icon: 'success',
            title: 'Berhasil Disimpan!',
            text: currentUser.value.role === 'owner' 
                ? 'Seluruh jadwal toko resmi aktif! 🟢' 
                : 'Jadwal toko berhasil diajukan, menunggu ACC Owner! ⏳',
            timer: 2000,
            showConfirmButton: false
        });
        
        fetchData(); 
    } catch (e) {
        console.error("ERROR TSM SAVE:", e);
        Swal.fire('Gagal', 'Terjadi kesalahan sistem saat menyimpan master jadwal', 'error');
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

        await api.post('/schedules/bulk', { schedules: singlePayload });
        Swal.fire('Updated!', `Jadwal berhasil di-${action}!`, 'success');
        fetchData();
    } catch (e) {
        Swal.fire('Gagal', 'Gagal memproses approval owner.', 'error');
    }
};

onMounted(() => fetchData());
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-8 max-w-7xl mx-auto font-sans">
            <div class="bg-gradient-to-r from-slate-800 to-indigo-900 rounded-3xl p-6 md:p-8 mb-8 text-white shadow-xl flex flex-col md:flex-row items-center justify-between relative border border-slate-700">
                <div class="absolute -right-10 -top-20 opacity-10 text-[180px] font-black italic pointer-events-none">📅</div>
                <div class="z-10 text-center md:text-left mb-4 md:mb-0">
                    <h1 class="text-2xl font-black tracking-tight mb-2 uppercase">Toko Shift Management (TSM)</h1>
                    <p class="text-indigo-200 font-medium text-xs italic">Periode Mingguan: {{ startDate }} s/d {{ endDate }}</p>
                </div>
                <div v-if="canEditSchedule" class="z-10">
                    <button @click="handleSaveJadwalBulk" class="bg-yellow-400 hover:bg-yellow-500 text-slate-900 font-black px-6 py-3.5 rounded-2xl text-sm shadow-md transition-all uppercase tracking-tight">
                        💾 Simpan Jadwal Toko
                    </button>
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
                <div class="p-5 md:p-6 border-b border-gray-100 bg-gray-50/50">
                    <h3 class="font-black text-gray-800 text-base uppercase tracking-tight">📋 Matriks Jadwal Kerja</h3>
                </div>

                <div class="overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left whitespace-nowrap border-collapse">
                        <thead>
                            <tr class="bg-slate-100 border-b border-gray-200 text-[10px] font-black text-slate-500 uppercase tracking-wider">
                                <th class="px-6 py-4 sticky left-0 bg-slate-100 z-10 border-r border-gray-200 w-52">Nama Karyawan</th>
                                <th v-for="d in mingguJadwal" :key="d.tanggal" class="px-4 py-4 text-center min-w-[150px]">
                                    <div class="text-indigo-600 font-black">{{ d.hari }}</div>
                                    <div class="text-slate-400 text-[9px] mt-0.5">{{ d.tanggal }}</div>
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-100">
                            <tr v-if="isLoading">
                                <td :colspan="8" class="px-6 py-12 text-center text-gray-400 font-bold uppercase animate-pulse">Memuat Matriks Jadwal...</td>
                            </tr>
                            <tr v-else v-for="emp in listKaryawan" :key="emp.id" class="hover:bg-slate-50/50 transition-colors">
                                <td class="px-6 py-4 sticky left-0 bg-white group-hover:bg-slate-50 z-10 border-r border-gray-100 font-black text-gray-700 text-sm uppercase shadow-[4px_0_8px_rgba(0,0,0,0.02)]">
                                    {{ emp.name }}
                                </td>

                                <td v-for="d in mingguJadwal" :key="d.tanggal" class="p-3 text-center">
                                    <div class="flex flex-col items-center gap-1.5">
                                        <select 
                                        v-model="formJadwal[emp.id || emp.user_id][d.tanggal]"
                                        :disabled="!canEditSchedule || formJadwal[emp.id][d.tanggal].includes('(Approved)')"
                                        class="w-full px-2 py-2 text-xs font-black rounded-xl border-2 text-center transition-all outline-none"
                                        :class="{
                                            'bg-green-50 text-green-700 border-green-200': formJadwal[emp.id][d.tanggal].includes('Shift 1'),'bg-blue-50 text-blue-700 border-blue-200': formJadwal[emp.id][d.tanggal].includes('Shift 2'),'bg-purple-50 text-purple-700 border-purple-200': formJadwal[emp.id][d.tanggal].includes('Middle'),'bg-slate-100 text-slate-400 border-slate-200': formJadwal[emp.id][d.tanggal] === 'OFF'
                                            }">
                                            <option value="Shift 1">SHIFT 1</option>
                                            <option value="Middle">MIDDLE</option>
                                            <option value="Shift 2">SHIFT 2</option>
                                            <option value="Shift 3">SHIFT 3</option>
                                            <option value="OFF">❌ LIBUR (OFF)</option>
                                            <slot v-if="formJadwal[emp.id || emp.user_id][d.tanggal]">
                                                <option hidden :value="formJadwal[emp.id || emp.user_id][d.tanggal]">
                                                    {{ formJadwal[emp.id || emp.user_id][d.tanggal] }}
                                                </option>
                                            </slot>
                                        </select>

                                        <div class="text-[8px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md"
                                             v-if="formJadwal[emp.id][d.tanggal] !== 'OFF'">
                                            <span v-if="formJadwal[emp.id][d.tanggal].includes('(Pending)')" class="bg-amber-100 text-amber-700">⏳ PENDING</span>
                                            <span v-else-if="formJadwal[emp.id][d.tanggal].includes('(Approved)')" class="bg-green-100 text-green-700">✅ APPROVED</span>
                                            <span v-else-if="formJadwal[emp.id][d.tanggal].includes('(Rejected)')" class="bg-red-100 text-red-700">❌ REJECTED</span>
                                            <span v-else class="bg-gray-100 text-gray-500">✍️ DRAFT</span>
                                        </div>

                                        <div v-if="currentUser.role === 'owner' && formJadwal[emp.id][d.tanggal].includes('(Pending)')" class="flex gap-1 mt-1 justify-center">
                                            <button @click="handleApproval(emp.id, d.tanggal, 'approve')" class="bg-green-600 hover:bg-green-700 text-white text-[9px] font-bold px-2 py-1 rounded shadow-sm">ACC</button>
                                            <button @click="handleApproval(emp.id, d.tanggal, 'reject')" class="bg-red-600 hover:bg-red-700 text-white text-[9px] font-bold px-2 py-1 rounded shadow-sm">TOLAK</button>
                                        </div>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>