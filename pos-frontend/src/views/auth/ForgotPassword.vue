<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const router = useRouter();
const identifier = ref('');
const isLoading = ref(false);

const goBack = () => {
    router.push('/login');
};

const handleCheckAccount = async () => {
    if (!identifier.value) return;
    
    isLoading.value = true;
    try {
        // 🚀 Tembak API Golang untuk check account asli
        const res = await api.post('/auth/check-account', { 
            identifier: identifier.value 
        });
        
        // Ambil data real dari database Golang lu
        const { email, phone } = res.data;

        // Lempar ke halaman verifikasi dengan data asli hasil dari database
        router.push({ 
            path: '/select-verify', 
            query: { 
                email: email, 
                phone: phone, 
                intent: 'reset-password' 
            } 
        });

    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Akun Tidak Ditemukan',
            text: error.response?.data?.error || 'Email atau nomor WA tidak terdaftar di sistem.',
            confirmButtonColor: '#ef4444'
        });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] font-sans p-4 relative overflow-hidden">
    
    <div class="absolute -top-24 -left-24 w-96 h-96 bg-blue-100/50 rounded-full blur-3xl"></div>
    <div class="absolute -bottom-24 -right-24 w-96 h-96 bg-indigo-100/50 rounded-full blur-3xl"></div>

    <div class="w-full max-w-md relative z-10">
      <div class="bg-white rounded-[40px] p-8 md:p-10 shadow-2xl relative border border-white flex flex-col text-center">
        
        <div class="w-20 h-20 bg-blue-50 text-blue-600 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-sm border border-blue-100/30">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
        </div>

        <h2 class="text-3xl font-black text-slate-900 tracking-tighter uppercase mb-2">Lupa Password?</h2>
        <p class="text-slate-400 font-bold text-[10px] uppercase tracking-widest leading-relaxed mb-8">
            Masukkan Email atau Nomor WhatsApp <br/>yang terdaftar pada akun Anda.
        </p>

        <form @submit.prevent="handleCheckAccount" class="space-y-6 text-left">
          <div class="space-y-2">
            <div class="relative group">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-600 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              </div>
              <input 
                v-model="identifier" 
                type="text" 
                required 
                class="w-full pl-12 pr-4 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-blue-600 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300"
                placeholder="Email atau No. WA Anda" 
              > 
            </div>
          </div>

          <button 
            type="submit" 
            :disabled="isLoading || !identifier"
            class="w-full py-4 rounded-2xl bg-blue-600 text-white font-black text-sm uppercase tracking-widest shadow-xl shadow-blue-200 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50"
          >
             {{ isLoading ? 'MENCARI AKUN...' : 'LANJUTKAN' }}
          </button>
        </form>

        <button 
            @click="goBack"
            class="mt-6 w-full py-3 rounded-xl bg-transparent text-slate-400 font-bold text-[10px] uppercase tracking-widest hover:text-slate-700 transition-all flex items-center justify-center gap-2"
        >
            Batal
        </button>

      </div>
    </div>
  </div>
</template>