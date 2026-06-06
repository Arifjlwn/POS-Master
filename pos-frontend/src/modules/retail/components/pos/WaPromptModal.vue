<script setup>
import { ref, watch, nextTick } from 'vue';

const props = defineProps({
    show: Boolean,
    totalBelanja: Number
});

const emit = defineEmits(['close', 'skip', 'submit']);
const phone = ref('');
const inputRef = ref(null);

// 🚀 UX SMART: Setiap kali modal WA muncul, paksa kursor langsung berkedip di kotak input otomatis bray!
watch(() => props.show, (newVal) => {
    if (newVal) {
        nextTick(() => {
            if (inputRef.value) inputRef.value.focus();
        });
    }
});

const handleSubmit = () => {
    // Bersihkan spasi atau karakter hantu jika ada bray
    let rawPhone = String(phone.value).trim();

    // 🛡️ SECURITY BLOCK: Kunci total tombol enter agar tidak bisa bypass data sampah bray
    if (rawPhone.length < 10) return;

    // 🛡️ FONNTE API NORMALIZATION: Ubah otomatis format 08xx menjadi 628xx demi kelancaran gateway bray!
    if (rawPhone.startsWith('0')) {
        rawPhone = '62' + rawPhone.slice(1);
    } else if (!rawPhone.startsWith('62') && rawPhone.startsWith('8')) {
        rawPhone = '62' + rawPhone;
    }

    emit('submit', rawPhone);
    phone.value = ''; // reset state
};

const handleSkip = () => {
    phone.value = ''; // reset state
    emit('skip');
};
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[200] p-4 backdrop-blur-sm print:hidden">
        <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-sm overflow-hidden flex flex-col border-t-8 border-emerald-500 animate-slide-in-right">
            
            <div class="p-6 md:p-8 text-center">
                <!-- ICON BRAND DESIGN EMERALD -->
                <div class="w-16 h-16 bg-emerald-50 text-emerald-500 rounded-full flex items-center justify-center mx-auto mb-4 border border-emerald-100 shadow-inner">
                    <svg class="w-8 h-8" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M3 21l1.65-3.8a9 9 0 1 1 3.4 2.9L3 21" /><path stroke-linecap="round" stroke-linejoin="round" d="M9 10a.5.5 0 0 0 1 0V9a.5.5 0 0 0-1 0v1a5 5 0 0 0 5 5h1a.5.5 0 0 0 0-1h-1a.5.5 0 0 0 0 1" /></svg>
                </div>
                
                <h2 class="text-xl font-black text-slate-900 tracking-tight mb-1">Kirim E-Struk WA?</h2>
                <p class="text-xs font-bold text-slate-500 mb-6">Total Tagihan: <span class="text-emerald-600 font-black">Rp {{ totalBelanja.toLocaleString('id-ID') }}</span></p>

                <!-- INPUT FIELD WA PELANGGAN -->
                <div class="relative mb-6 text-left">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-1 block">Nomor WhatsApp Pelanggan</label>
                    <!-- Ditukar ke type text dengan pattern angka agar aman dari bug limit exponensial HTML5 bray -->
                    <input 
                        ref="inputRef"
                        v-model="phone" 
                        type="text" 
                        inputmode="numeric"
                        oninput="this.value = this.value.replace(/\D/g, '')"
                        class="w-full px-5 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-emerald-500 focus:ring-4 focus:ring-emerald-500/10 outline-none font-bold text-slate-800 transition-all placeholder:text-slate-300" 
                        placeholder="Contoh: 08123456789"
                        @keyup.enter="handleSubmit"
                    >
                </div>

                <!-- ACTION BUTTON ACTION -->
                <div class="flex flex-col gap-3">
                    <button @click="handleSubmit" :disabled="phone.length < 10" class="w-full bg-emerald-500 text-white py-4 rounded-xl font-black text-[11px] uppercase tracking-widest shadow-lg shadow-emerald-500/30 hover:bg-emerald-600 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
                        Kirim & Cetak Struk
                    </button>
                    <button @click="handleSkip" class="w-full bg-slate-100 text-slate-500 py-4 rounded-xl font-black text-[11px] uppercase tracking-widest hover:bg-slate-200 transition-all active:scale-95">
                        Lewati / Print Thermal Saja
                    </button>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped>
/* Transisi Smooth Slide-In */
.animate-slide-in-right {
    animation: slideIn 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes slideIn {
    from { transform: translateY(20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
}
</style>