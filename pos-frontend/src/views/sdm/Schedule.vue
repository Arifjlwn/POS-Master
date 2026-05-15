<script setup>
import { ref, onMounted, computed } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
import Swal from 'sweetalert2';

// --- STATE DATA ---
const listKaryawan = ref([]);
const listJadwal = ref([]);
const isSaving = ref(false);
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

        // 🚀 CETAK CONSOLE LOG DI DALAM BLOK AMAN (Sebelum Tembak API)
        console.log("PAYLOAD YANG DIKIRIM KE GO:", payloadSchedules);

        await api.post('/schedules/bulk', { schedules: payloadSchedules });
        
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
        console.error("ERROR TSM SAVE:", e);
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
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans">
            <div class="bg-gradient-to-br from-slate-900 to-blue-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-slate-800">
                <svg xmlns="http://www.w3.org/2000/svg" class="absolute -right-10 -bottom-10 w-64 h-64 text-blue-500 opacity-10 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="M8 14h.01"/><path d="M12 14h.01"/><path d="M16 14h.01"/><path d="M8 18h.01"/><path d="M12 18h.01"/><path d="M16 18h.01"/></svg>
                
                <div class="z-10 text-center md:text-left mb-6 md:mb-0">
                    <h1 class="text-3xl font-black tracking-tighter mb-2 uppercase italic">Shift <span class="text-blue-500">Management</span></h1>
                    <p class="text-blue-200 font-bold text-[10px] uppercase tracking-widest flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                        Periode: {{ startDate }} s/d {{ endDate }}
                    </p>
                </div>
                
                <div v-if="canEditSchedule" class="z-10">
                    <button
                    @click="handleSaveJadwalBulk" 
                    :disabled="isSaving"
                    class="bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-[24px] font-black text-[10px] uppercase tracking-[0.2em] shadow-xl shadow-blue-900/50 flex items-center gap-3 transition-all active:scale-95 border border-blue-400/30 disabled:opacity-70 disabled:cursor-not-allowed"
                    >
                        <template v-if="isSaving">
                            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            Menyimpan...
                        </template>
                        
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
                                <polyline points="17 21 17 13 7 13 7 21"/>
                                <polyline points="7 3 7 8 15 8"/>
                            </svg>
                            Simpan Master Jadwal
                        </template>
                    </button>
                </div>
            </div>

            <div class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
                <div class="p-6 border-b border-slate-50 bg-slate-50/50">
                    <h3 class="font-black text-slate-800 text-sm flex items-center gap-2 uppercase tracking-widest">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                        Matriks Penjadwalan Toko
                    </h3>
                </div>

                <div class="overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left whitespace-nowrap border-collapse">
                        <thead>
                            <tr class="bg-white border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                                <th class="px-6 py-5 sticky left-0 bg-white z-20 border-r border-slate-100 w-52 shadow-[4px_0_12px_rgba(0,0,0,0.03)]">Nama Karyawan</th>
                                <th v-for="d in mingguJadwal" :key="d.tanggal" class="px-4 py-5 text-center min-w-[160px] border-l border-slate-50">
                                    <div class="text-blue-600 font-black text-xs">{{ d.hari }}</div>
                                    <div class="text-slate-400 text-[9px] font-bold mt-1 tracking-widest">{{ d.tanggal }}</div>
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td :colspan="8" class="px-6 py-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">
                                    Memuat Data Matriks...
                                </td>
                            </tr>
                            <tr v-else v-for="emp in listKaryawan" :key="emp.id" class="hover:bg-slate-50/50 transition-colors group">
                                <td class="px-6 py-5 sticky left-0 bg-white group-hover:bg-slate-50 z-10 border-r border-slate-100 shadow-[4px_0_12px_rgba(0,0,0,0.03)] transition-colors">
                                    <div class="font-black text-slate-800 text-sm uppercase tracking-tighter">{{ emp.name }}</div>
                                    <div class="text-[9px] font-bold text-slate-400 mt-1 uppercase tracking-widest">{{ emp.role }}</div>
                                </td>

                                <td v-for="d in mingguJadwal" :key="d.tanggal" class="p-4 text-center border-l border-slate-50">
                                    <div class="flex flex-col items-center gap-2">
                                        <select 
                                        v-model="formJadwal[emp.id || emp.user_id][d.tanggal]"
                                        :disabled="!canEditSchedule || formJadwal[emp.id][d.tanggal].includes('(Approved)')"
                                        class="w-full px-3 py-2.5 text-[11px] font-black uppercase tracking-widest rounded-xl border-2 text-center transition-all outline-none cursor-pointer appearance-none"
                                        :class="{
                                            'bg-emerald-50 text-emerald-700 border-emerald-200 focus:border-emerald-500': formJadwal[emp.id][d.tanggal].includes('Shift 1'),
                                            'bg-blue-50 text-blue-700 border-blue-200 focus:border-blue-500': formJadwal[emp.id][d.tanggal].includes('Shift 2'),
                                            'bg-purple-50 text-purple-700 border-purple-200 focus:border-purple-500': formJadwal[emp.id][d.tanggal].includes('Middle'),
                                            'bg-slate-50 text-slate-400 border-slate-200 focus:border-slate-400': formJadwal[emp.id][d.tanggal] === 'OFF',
                                            'bg-slate-900 text-white border-slate-800': formJadwal[emp.id][d.tanggal].includes('Shift 3')
                                            }">
                                            <option value="Shift 1">SHIFT 1</option>
                                            <option value="Middle">MIDDLE</option>
                                            <option value="Shift 2">SHIFT 2</option>
                                            <option value="Shift 3">SHIFT 3</option>
                                            <option value="OFF">LIBUR (OFF)</option>
                                            <slot v-if="formJadwal[emp.id || emp.user_id][d.tanggal]">
                                                <option hidden :value="formJadwal[emp.id || emp.user_id][d.tanggal]">
                                                    {{ formJadwal[emp.id || emp.user_id][d.tanggal] }}
                                                </option>
                                            </slot>
                                        </select>

                                        <div class="flex justify-center w-full" v-if="formJadwal[emp.id][d.tanggal] !== 'OFF'">
                                            <span v-if="formJadwal[emp.id][d.tanggal].includes('(Pending)')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-amber-50 text-amber-600 border border-amber-100 text-[9px] font-black uppercase tracking-widest">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                                Pending
                                            </span>
                                            <span v-else-if="formJadwal[emp.id][d.tanggal].includes('(Approved)')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100 text-[9px] font-black uppercase tracking-widest">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                                Approved
                                            </span>
                                            <span v-else-if="formJadwal[emp.id][d.tanggal].includes('(Rejected)')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-50 text-red-600 border border-red-100 text-[9px] font-black uppercase tracking-widest">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
                                                Rejected
                                            </span>
                                            <span v-else class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-100 text-slate-500 text-[9px] font-black uppercase tracking-widest">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                                                Draft
                                            </span>
                                        </div>

                                        <div v-if="currentUser.role === 'owner' && formJadwal[emp.id][d.tanggal].includes('(Pending)')" class="flex gap-2 mt-1 justify-center">
                                            <button @click="handleApproval(emp.id, d.tanggal, 'approve')" class="p-2 rounded-xl bg-emerald-100 text-emerald-600 hover:bg-emerald-500 hover:text-white transition-all shadow-sm border border-emerald-200" title="Setujui Jadwal">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                                            </button>
                                            <button @click="handleApproval(emp.id, d.tanggal, 'reject')" class="p-2 rounded-xl bg-red-100 text-red-600 hover:bg-red-500 hover:text-white transition-all shadow-sm border border-red-200" title="Tolak Jadwal">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                                            </button>
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
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>