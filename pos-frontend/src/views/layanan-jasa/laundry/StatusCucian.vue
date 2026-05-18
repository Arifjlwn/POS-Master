<script setup>
import { ref, onMounted, computed } from 'vue';
import SidebarLaundry from './SidebarLaundry.vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

const cucianList = ref([]);
const isLoading = ref(false);
const searchQuery = ref('');

const formatRupiah = (angka) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
};

const fetchData = async () => {
    isLoading.value = true;
    try {
        const response = await api.get('/laundry/tracking');
        cucianList.value = response.data || [];
    } catch (error) {
        Swal.fire('Gagal!', 'Gagal memuat status cucian.', 'error');
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => fetchData());

const filteredList = computed(() => {
    if (!searchQuery.value) return cucianList.value;
    return cucianList.value.filter(item => 
        item.pelanggan.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        item.invoice.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
});

// 🚀 FUNGSI UPDATE STATUS
const updateStatus = async (item, statusBaru) => {
    try {
        await api.put(`/laundry/tracking/${item.id}/status`, { status_baru: statusBaru });
        
        // Update state lokal biar UI langsung berubah tanpa refresh
        item.status = statusBaru;

        // 🚀🔥 LOGIKA WHATSAPP GRATIS JIKA SELESAI
        if (statusBaru === 'SELESAI' && item.whatsapp) {
            Swal.fire({
                title: 'Kirim Notif WA?',
                text: "Status cucian sudah selesai, mau kabarin pelanggan sekarang?",
                icon: 'question',
                showCancelButton: true,
                confirmButtonColor: '#10b981', // Emerald
                cancelButtonColor: '#94a3b8',
                confirmButtonText: 'Ya, Buka WhatsApp!'
            }).then((result) => {
                if (result.isConfirmed) {
                    // Tulis teksnya normal aja pakai Enter biasa (\n)
                    const textRaw = `Halo Kak *${item.pelanggan}*,

Cucian kakak dengan Nomor Resi:
*${item.invoice}* - (${item.layanan})
*SUDAH SELESAI* dicuci. 🥳✨

Silakan diambil di toko ya kak.
Terima kasih sudah mempercayakan cuciannya di Arzu Laundry ! 💙`;
                    
                    // 🚀 BUNGKUS PAKAI encodeURIComponent BIAR SPASI GAK BIKIN PUTUS!
                    const textEncoded = encodeURIComponent(textRaw);
                    
                    // Buka tab WhatsApp Web / App
                    window.open(`https://wa.me/${item.whatsapp}?text=${textEncoded}`, '_blank');
                }
            });
        } else {
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Status Diperbarui!', showConfirmButton: false, timer: 1500 });
        }
        
        // Kalau diambil, hapus dari list tracking aktif
        if(statusBaru === 'DIAMBIL') {
            cucianList.value = cucianList.value.filter(c => c.id !== item.id);
        }

    } catch (error) {
        Swal.fire('Gagal!', 'Terjadi kesalahan sistem.', 'error');
    }
};
</script>

<template>
    <SidebarLaundry>
        <div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden relative">
            
            <div class="p-5 md:p-8 shrink-0 bg-white border-b border-slate-100 flex flex-col z-10 shadow-sm relative gap-5">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-sky-50 border border-sky-100 rounded-2xl flex items-center justify-center shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-sky-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                    </div>
                    <div>
                        <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-slate-800 leading-tight">Tracking Cucian</h1>
                        <p class="text-[9px] md:text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Pantau Status & Notif Pelanggan</p>
                    </div>
                </div>

                <div class="relative w-full group max-w-2xl">
                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-sky-500 transition-colors" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                    </div>
                    <input v-model="searchQuery" type="text" placeholder="Cari No. Resi atau Nama Pelanggan..." class="w-full pl-12 pr-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-xl focus:bg-white focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 outline-none font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300">
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-8">
                
                <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 text-slate-400">
                    <div class="w-12 h-12 border-4 border-sky-200 border-t-sky-600 rounded-full animate-spin mb-4"></div>
                    <p class="font-black text-xs uppercase tracking-[0.2em] animate-pulse">Memuat Data...</p>
                </div>
                
                <div v-else-if="filteredList.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-400 bg-white rounded-3xl border border-slate-100 border-dashed border-2">
                    <div class="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center mb-4">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22v-5"/><path d="M9 7V2"/><path d="M15 7V2"/><rect width="16" height="15" x="4" y="7" rx="2" ry="2"/></svg>
                    </div>
                    <p class="font-black text-sm uppercase tracking-widest text-slate-500 mb-1">Pekerjaan Beres!</p>
                    <p class="text-xs font-bold">Tidak ada antrian cucian saat ini.</p>
                </div>

                <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">
                    <div v-for="item in filteredList" :key="item.id" class="bg-white rounded-[24px] border border-slate-200 p-5 shadow-sm hover:shadow-xl transition-all flex flex-col relative overflow-hidden">
                        
                        <div class="absolute top-0 right-0 px-4 py-1.5 rounded-bl-xl font-black text-[10px] uppercase tracking-widest"
                             :class="{
                                 'bg-amber-100 text-amber-600': item.status === 'ANTRI',
                                 'bg-sky-100 text-sky-600': item.status === 'DICUCI',
                                 'bg-emerald-100 text-emerald-600': item.status === 'SELESAI'
                             }">
                            {{ item.status }}
                        </div>

                        <div class="mb-4 pr-16">
                            <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">{{ item.invoice }}</span>
                            <h3 class="text-base font-black text-slate-800 uppercase leading-tight mt-1 truncate">{{ item.pelanggan }}</h3>
                            <div class="flex items-center gap-1.5 mt-1 text-slate-500">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                                <span class="text-xs font-bold">{{ item.whatsapp }}</span>
                            </div>
                        </div>

                        <div class="bg-slate-50 p-3 rounded-xl border border-slate-100 mb-5">
                            <h4 class="text-[11px] font-black text-indigo-600 uppercase tracking-wider truncate mb-2">{{ item.layanan }}</h4>
                            <div class="flex justify-between items-end border-t border-slate-200/60 pt-2">
                                <div>
                                    <span class="block text-[9px] font-bold text-slate-400 uppercase">Kuantitas / Qty</span>
                                    <span class="text-sm font-black text-slate-700 flex items-baseline gap-1">
                                        {{ item.berat_kg }} 
                                        <span class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">
                                            {{ item.satuan ? item.satuan : (item.berat_kg % 1 !== 0 ? 'KG' : 'Item') }}
                                        </span>
                                    </span>
                                </div>
                                <div class="text-right">
                                    <span class="block text-[9px] font-bold text-slate-400 uppercase">Sub Total</span>
                                    <span class="text-sm font-black text-slate-700">{{ formatRupiah(item.sub_total) }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="mt-auto">
                            <button v-if="item.status === 'ANTRI'" @click="updateStatus(item, 'DICUCI')" class="w-full bg-sky-50 hover:bg-sky-500 text-sky-600 hover:text-white border border-sky-200 py-3 rounded-xl font-black text-xs uppercase tracking-widest transition-all active:scale-95 flex items-center justify-center gap-2 group">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 transition-transform group-hover:rotate-180" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                Mulai Dicuci
                            </button>

                            <button v-if="item.status === 'DICUCI'" @click="updateStatus(item, 'SELESAI')" class="w-full bg-emerald-500 hover:bg-emerald-600 text-white shadow-lg shadow-emerald-200 py-3 rounded-xl font-black text-xs uppercase tracking-widest transition-all active:scale-95 flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                                Selesai & WA Pelanggan
                            </button>

                            <button v-if="item.status === 'SELESAI'" @click="updateStatus(item, 'DIAMBIL')" class="w-full bg-slate-800 hover:bg-slate-900 text-white shadow-md py-3 rounded-xl font-black text-xs uppercase tracking-widest transition-all active:scale-95 flex items-center justify-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                                Tandai Telah Diambil
                            </button>
                        </div>

                    </div>
                </div>

            </div>
        </div>
    </SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>