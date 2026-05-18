<script setup>
import { ref, onMounted } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js'; 
import Swal from 'sweetalert2';

const isLoading = ref(false);
const isSaving = ref(false);

// 🚀 FORM BERSIH TANPA EMAIL/PASSWORD KASIR
const form = ref({
    nama_toko: '',
    telepon: '',
    alamat: '',
    payment_type: 'PRIBADI',
    receipt_footer: '',
    qris_base64: ''
});

const qrisPreviewUrl = ref('');

const formatNoHpToko = () => {
    let val = String(form.value.telepon);
    if (val.startsWith('0')) val = val.substring(1);
    if (val.startsWith('62')) val = val.substring(2);
    form.value.telepon = val;
};

const fetchSetting = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/setting');
        form.value.nama_toko = response.data.nama_toko;
        form.value.alamat = response.data.alamat;
        form.value.payment_type = response.data.payment_type || 'PRIBADI';
        form.value.receipt_footer = response.data.receipt_footer || 'Cucian tidak diambil dalam 30 hari bukan tanggung jawab kami.';
        
        let phone = String(response.data.telepon || '');
        if (phone.startsWith('62')) phone = phone.substring(2);
        else if (phone.startsWith('0')) phone = phone.substring(1);
        form.value.telepon = phone;
        
        if (response.data.qris_image) {
            qrisPreviewUrl.value = response.data.qris_image;
        }
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal memuat pengaturan toko.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchSetting());

const handleQrisUpload = (event) => {
    const file = event.target.files[0];
    if (!file) return;
    if (file.size > 2 * 1024 * 1024) return Swal.fire('Oops!', 'Ukuran gambar maksimal 2 MB beb!', 'warning');

    const reader = new FileReader();
    reader.onload = (e) => {
        form.value.qris_base64 = e.target.result; 
        qrisPreviewUrl.value = e.target.result; 
    };
    reader.readAsDataURL(file);
};

const saveSetting = async () => {
    isSaving.value = true;
    try {
        const payload = {
            ...form.value,
            telepon: form.value.telepon ? `62${form.value.telepon}` : ''
        };

        await api.put('/laundry/setting', payload);
        Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Pengaturan Berhasil Disimpan!', showConfirmButton: false, timer: 1500 });
        
        localStorage.setItem('store_name', form.value.nama_toko);
        
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal menyimpan konfigurasi.', 'error');
    } finally {
        isSaving.value = false;
    }
};
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden relative">
            
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 z-10 shadow-sm">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V3a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
                    </div>
                    <div>
                        <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Pengaturan Toko</h1>
                        <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Konfigurasi Profil Bisnis & Sistem Pembayaran</p>
                    </div>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8">
                
                <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 text-slate-400">
                    <div class="w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                    <p class="font-black text-xs uppercase tracking-widest">Memuat Konfigurasi...</p>
                </div>

                <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8 items-start">
                    
                    <div class="lg:col-span-2 space-y-6 md:space-y-8">
                        
                        <div class="bg-white border border-slate-200 p-6 md:p-8 rounded-[24px] shadow-sm space-y-6">
                            <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest border-b border-slate-100 pb-3">Identitas Bisnis & Cetakan Struk</h3>
                            
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                                <div>
                                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Toko Laundry</label>
                                    <input v-model="form.nama_toko" type="text" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 font-bold text-sm text-slate-800 transition-all">
                                </div>
                                <div>
                                    <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nomor Telepon WhatsApp</label>
                                    <div class="flex items-center bg-slate-50 border-2 border-slate-200 rounded-xl focus-within:bg-white focus-within:border-indigo-500 transition-all overflow-hidden">
                                        <div class="pl-3 pr-2 py-3 bg-slate-100 border-r border-slate-200 flex items-center justify-center select-none"><span class="text-slate-500 font-black text-xs">+62</span></div>
                                        <input v-model="form.telepon" @input="formatNoHpToko" type="number" class="w-full px-3 py-3 bg-transparent outline-none font-bold text-slate-800 text-sm" placeholder="8123...">
                                    </div>
                                </div>
                            </div>

                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Alamat Lengkap (Tampil di Nota)</label>
                                <textarea v-model="form.alamat" rows="2" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 font-bold text-sm text-slate-800 transition-all resize-none"></textarea>
                            </div>

                            <div>
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Pesan Footer Struk (Syarat & Ketentuan)</label>
                                <textarea v-model="form.receipt_footer" rows="2" class="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 font-bold text-sm text-slate-800 transition-all resize-none"></textarea>
                            </div>
                        </div>

                    </div>

                    <div class="bg-white border border-slate-200 p-6 rounded-[24px] shadow-sm space-y-6 flex flex-col justify-between">
                        <div>
                            <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest border-b border-slate-100 pb-3 mb-5">Metode Gerbang QRIS</h3>
                            
                            <div class="space-y-3">
                                <label @click="form.payment_type = 'PRIBADI'" :class="form.payment_type === 'PRIBADI' ? 'border-indigo-500 bg-indigo-50/40' : 'border-slate-200 bg-white'" class="flex items-center justify-between p-4 border-2 rounded-2xl cursor-pointer transition-all hover:border-indigo-300">
                                    <div class="flex items-center gap-3">
                                        <div class="w-4 h-4 rounded-full border-2 border-slate-300 flex items-center justify-center" :class="form.payment_type === 'PRIBADI' ? 'border-indigo-500' : ''">
                                            <div v-if="form.payment_type === 'PRIBADI'" class="w-2 h-2 rounded-full bg-indigo-600"></div>
                                        </div>
                                        <div>
                                            <span class="block text-xs font-black text-slate-800 uppercase tracking-wide">QRIS Manual / Pribadi</span>
                                            <span class="text-[10px] font-bold text-slate-400">Scan Barcode Akrilik Toko Anda</span>
                                        </div>
                                    </div>
                                </label>

                                <label @click="form.payment_type = 'GATEWAY'" :class="form.payment_type === 'GATEWAY' ? 'border-indigo-500 bg-indigo-50/40' : 'border-slate-200 bg-white'" class="flex items-center justify-between p-4 border-2 rounded-2xl cursor-pointer transition-all hover:border-indigo-300">
                                    <div class="flex items-center gap-3">
                                        <div class="w-4 h-4 rounded-full border-2 border-slate-300 flex items-center justify-center" :class="form.payment_type === 'GATEWAY' ? 'border-indigo-500' : ''">
                                            <div v-if="form.payment_type === 'GATEWAY'" class="w-2 h-2 rounded-full bg-indigo-600"></div>
                                        </div>
                                        <div>
                                            <span class="block text-xs font-black text-slate-800 uppercase tracking-wide">Payment Gateway</span>
                                            <span class="text-[10px] font-bold text-slate-400">Otomatis Terintegrasi (Coming Soon)</span>
                                        </div>
                                    </div>
                                </label>
                            </div>

                            <div v-if="form.payment_type === 'PRIBADI'" class="mt-6 border-t border-slate-100 pt-5 animate-[fadeIn_0.3s_ease-out] flex flex-col items-center">
                                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3 self-start">Upload File Gambar QRIS</label>
                                
                                <div class="w-40 h-40 bg-slate-50 border-2 border-dashed border-slate-200 rounded-2xl flex flex-col items-center justify-center relative overflow-hidden group shadow-inner">
                                    <img v-if="qrisPreviewUrl" :src="qrisPreviewUrl" class="w-full h-full object-contain p-2 mix-blend-multiply" />
                                    <div v-else class="flex flex-col items-center text-slate-400">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 mb-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
                                        <span class="text-[9px] font-black uppercase tracking-wider">Cari Foto</span>
                                    </div>
                                    <input type="file" @change="handleQrisUpload" accept="image/*" class="absolute inset-0 opacity-0 cursor-pointer" />
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="mt-8 flex justify-end">
                    <button @click="saveSetting" :disabled="isSaving" class="w-full lg:w-auto bg-slate-900 hover:bg-slate-800 text-white px-10 py-4 rounded-xl font-black text-xs uppercase tracking-[0.2em] transition-all active:scale-95 shadow-xl shadow-slate-900/20 flex items-center justify-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/></svg>
                        {{ isSaving ? 'MENYIMPAN...' : 'SIMPAN SEMUA PERUBAHAN' }}
                    </button>
                </div>

            </div>
        </div>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
input[type=number]::-webkit-inner-spin-button, input[type=number]::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number] { -moz-appearance: textfield; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(-10px); } to { opacity: 1; transform: translateY(0); } }
</style>