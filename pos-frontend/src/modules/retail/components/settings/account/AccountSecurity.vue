<script setup>
const props = defineProps({ form: Object });
const emit = defineEmits(['submit']);
</script>

<template>
    <div class="space-y-6 animate-fade-in-up">
        <div class="bg-amber-50 p-4 rounded-2xl border border-amber-100 flex items-start gap-4 mb-6">
            <div class="p-2 bg-amber-100 text-amber-600 rounded-xl shrink-0"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg></div>
            <div>
                <h4 class="font-black text-amber-800 text-xs uppercase tracking-widest mb-1">Pembaruan Keamanan</h4>
                <p class="text-xs font-medium text-amber-700 leading-relaxed">Ganti password secara berkala untuk menjaga keamanan akun. Jika Anda Kasir dan lupa password lama, silakan hubungi Owner.</p>
            </div>
        </div>

        <div class="max-w-md space-y-4">
            <div>
                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Password Lama</label>
                <input v-model="form.old_password" type="password" placeholder="••••••••" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
            </div>
            <div>
                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Password Baru</label>
                <input v-model="form.new_password" type="password" placeholder="••••••••" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
            </div>
            <div>
                <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Konfirmasi Password Baru</label>
                <input v-model="form.confirm_password" type="password" placeholder="••••••••" class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm">
            </div>
        </div>

        <div class="pt-4 border-t border-slate-100">
            <button @click="emit('submit')" :disabled="!form.old_password || !form.new_password || !form.confirm_password" class="w-full md:w-auto px-8 py-4 bg-slate-900 hover:bg-blue-600 text-white rounded-2xl font-black text-[10px] uppercase tracking-widest shadow-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                Ubah Password
            </button>
        </div>
    </div>
</template>