<script setup>
import { ref, computed } from 'vue';
import SidebarFnB from './SidebarFnB.vue';
import Swal from 'sweetalert2';

// Dummy Data Orderan Masuk
const activeOrders = ref([
    {
        id: 101,
        invoice: 'KOT-001',
        tipe: 'DINE_IN',
        meja: 'Meja 05',
        waktu_pesan: '13:45',
        status: 'PROSES',
        items: [
            { nama: 'Nasi Goreng Arzu Spesial', qty: 2, notes: 'Pedas Mampus, telurnya dadar ya' },
            { nama: 'Es Teh Manis Jumbo', qty: 2, notes: 'Es nya dikit aja' }
        ]
    },
    {
        id: 102,
        invoice: 'KOT-002',
        tipe: 'TAKE_AWAY',
        meja: 'Bungkus (A/n Budi)',
        waktu_pesan: '13:50',
        status: 'PROSES',
        items: [
            { nama: 'Mie Ayam Pangsit Badas', qty: 1, notes: 'Gak pake seledri, kuah pisah' },
            { nama: 'Kopi Susu Gula Aren Dewa', qty: 1, notes: '' }
        ]
    }
]);

const selesaikanOrder = (order) => {
    Swal.fire({
        title: 'Pesanan Selesai?',
        text: `Makanan untuk ${order.meja} sudah siap dihidangkan?`,
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#10b981',
        confirmButtonText: 'Ya, Siap Dihidangkan!'
    }).then((result) => {
        if (result.isConfirmed) {
            activeOrders.value = activeOrders.value.filter(o => o.id !== order.id);
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Notif dikirim ke Kasir / Pelayan!', showConfirmButton: false, timer: 1500 });
        }
    });
};
</script>

<template>
    <SidebarFnB>
        <div class="flex-1 flex flex-col h-full bg-slate-900 overflow-hidden relative font-sans">
            <div class="p-5 md:p-6 bg-slate-900 border-b border-slate-800 flex justify-between items-center z-10 shadow-md shrink-0">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 bg-orange-500/20 border border-orange-500/50 rounded-2xl flex items-center justify-center shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-orange-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 13.87A4 4 0 0 1 7.41 6a5.11 5.11 0 0 1 1.05-1.54 5 5 0 0 1 7.08 0A5.11 5.11 0 0 1 16.59 6 4 4 0 0 1 18 13.87V21H6Z"/><line x1="6" y1="17" x2="18" y2="17"/></svg>
                    </div>
                    <div>
                        <h1 class="text-xl md:text-2xl font-black tracking-tighter uppercase text-white leading-tight">Kitchen Display</h1>
                        <p class="text-[10px] font-black text-slate-400 mt-1 uppercase tracking-widest">Monitor Koki & Barista</p>
                    </div>
                </div>
                <div class="bg-slate-800 px-5 py-2.5 rounded-xl border border-slate-700">
                    <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Antrean Masak</p>
                    <p class="text-xl font-black text-white leading-none">{{ activeOrders.length }} <span class="text-xs text-orange-400">Order</span></p>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-6">
                <div v-if="activeOrders.length === 0" class="h-full flex flex-col items-center justify-center text-slate-500">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 mb-4 opacity-20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" y1="9" x2="9.01" y2="9"/><line x1="15" y1="9" x2="15.01" y2="9"/></svg>
                    <p class="font-black text-sm uppercase tracking-widest text-slate-400">Dapur Kosong, Waktunya Ngopi!</p>
                </div>

                <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 items-start">
                    <div v-for="order in activeOrders" :key="order.id" class="bg-white rounded-3xl overflow-hidden shadow-xl flex flex-col animate-[fadeIn_0.3s_ease-out]">
                        <div class="p-4 flex justify-between items-center" :class="order.tipe === 'DINE_IN' ? 'bg-indigo-600' : 'bg-rose-500'">
                            <div>
                                <span class="text-[10px] font-black text-white/70 uppercase tracking-widest">{{ order.invoice }} • {{ order.waktu_pesan }}</span>
                                <h3 class="font-black text-lg text-white uppercase tracking-tight mt-0.5">{{ order.meja }}</h3>
                            </div>
                            <span class="bg-white/20 text-white border border-white/30 text-[10px] font-black px-3 py-1.5 rounded-lg uppercase tracking-widest">
                                {{ order.tipe === 'DINE_IN' ? 'Makan Sini' : 'Bungkus' }}
                            </span>
                        </div>

                        <div class="p-5 flex-1 bg-slate-50 space-y-3">
                            <div v-for="(item, index) in order.items" :key="index" class="flex gap-3 pb-3 border-b border-slate-200 border-dashed last:border-0 last:pb-0">
                                <div class="w-8 h-8 rounded-lg bg-slate-200 flex items-center justify-center shrink-0 font-black text-sm text-slate-700">{{ item.qty }}x</div>
                                <div>
                                    <p class="font-black text-sm text-slate-800 uppercase leading-tight">{{ item.nama }}</p>
                                    <p v-if="item.notes" class="mt-1.5 text-xs font-bold text-amber-900 bg-amber-100/80 border border-amber-200 px-2.5 py-1.5 rounded-md flex items-start gap-1.5">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 mt-0.5 shrink-0 text-amber-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M12 2v20M17 5H9.5a4.5 4.5 0 0 0 0 9H15a4.5 4.5 0 0 1 0 9H6.5"/></svg>
                                        {{ item.notes }}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div class="p-4 bg-white border-t border-slate-100">
                            <button @click="selesaikanOrder(order)" class="w-full bg-emerald-500 hover:bg-emerald-600 text-white py-3.5 rounded-xl font-black text-xs uppercase tracking-widest shadow-lg shadow-emerald-200 transition-all active:scale-95 flex justify-center items-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                                Siap Dihidangkan
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </SidebarFnB>
</template>
<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 8px; width: 8px; }
.custom-scrollbar::-webkit-scrollbar-track { background: rgba(255,255,255,0.05); }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #475569; border-radius: 10px; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(15px); } to { opacity: 1; transform: translateY(0); } }
</style>