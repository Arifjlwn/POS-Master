<script setup>
const props = defineProps({ form: Object });
</script>

<template>
    <div class="space-y-6 animate-fade-in-up">
        <div class="p-6 border-2 border-slate-100 rounded-3xl bg-white">
            <div class="flex items-center justify-between mb-6">
                <div>
                    <h4 class="font-black text-slate-800 uppercase tracking-widest text-sm">Pajak Pelanggan (PPN/PB1)</h4>
                    <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Aktifkan untuk membebankan pajak otomatis pada struk</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" v-model="form.is_tax_active" class="sr-only peer">
                    <div class="w-14 h-7 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-emerald-500"></div>
                </label>
            </div>

            <div v-if="form.is_tax_active" class="animate-fade-in-up">
                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Persentase Pajak (%)</label>
                <input v-model="form.pajak_persen" type="number" step="0.1" class="w-full md:w-1/3 p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-emerald-500 outline-none font-black text-lg text-emerald-600">
            </div>
        </div>

        <div class="p-6 border-2 border-slate-100 rounded-3xl bg-white">
            <h4 class="font-black text-slate-800 uppercase tracking-widest text-sm mb-1">Catatan Kaki Struk</h4>
            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-4">Pesan yang tercetak di bagian paling bawah struk thermal</p>
            
            <textarea v-model="form.receipt_footer" rows="3" placeholder="Terima kasih atas kunjungan Anda. Barang yang sudah dibeli tidak dapat ditukar/dikembalikan." class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm text-center resize-none"></textarea>
        </div>
    </div>
</template>