<script setup>
import Sidebar from '../../components/Sidebar.vue';
import AccountProfile from '../../components/settings/account/AccountProfile.vue';
import AccountSecurity from '../../components/settings/account/AccountSecurity.vue';
import { useAccount } from '../../composables/useAccount.js';

const { 
    isLoading, isSaving, activeTab, role, profileForm, passwordForm,
    fotoPreview, handleFileChange, saveProfile, updatePassword 
} = useAccount();

const tabs = [
    { id: 'profile', label: 'Profil Saya', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' },
    { id: 'security', label: 'Keamanan', icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z' }
];
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-5xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 class="text-2xl md:text-3xl font-black text-slate-800 tracking-tighter uppercase">Akun Saya</h1>
                    <p class="text-[10px] md:text-xs font-black text-slate-400 uppercase tracking-widest mt-1">Kelola Informasi Pribadi dan Keamanan Akun</p>
                </div>
            </div>

            <div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100 shadow-sm">
                <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-black text-[10px] uppercase tracking-widest animate-pulse">Menyiapkan Profil...</p>
            </div>

            <div v-else class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col md:flex-row">
                
                <div class="w-full md:w-64 bg-slate-50/50 border-r border-slate-100 p-6 flex flex-row md:flex-col gap-2 overflow-x-auto custom-scrollbar">
                    <button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id" 
                        :class="[
                            'flex items-center gap-3 p-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all whitespace-nowrap',
                            activeTab === tab.id ? 'bg-white text-indigo-600 shadow-sm border border-slate-200' : 'text-slate-500 hover:bg-slate-100'
                        ]">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                            <path stroke-linecap="round" stroke-linejoin="round" :d="tab.icon" />
                        </svg>
                        {{ tab.label }}
                    </button>
                    
                    <button v-if="role === 'owner'" @click="activeTab = 'billing'" 
                        :class="[
                            'flex items-center gap-3 p-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all whitespace-nowrap mt-4 border-t border-slate-200 pt-4',
                            activeTab === 'billing' ? 'bg-white text-amber-600 shadow-sm border border-slate-200' : 'text-amber-500 hover:bg-amber-50'
                        ]">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" /></svg>
                        Paket Berlangganan
                    </button>
                </div>

                <div class="flex-1 p-6 md:p-8 overflow-hidden relative">
                    <div v-if="isSaving" class="absolute inset-0 bg-white/60 backdrop-blur-sm z-10 flex items-center justify-center">
                        <div class="w-8 h-8 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
                    </div>

                    <AccountProfile v-show="activeTab === 'profile'" :form="profileForm" :fotoPreview="fotoPreview" @update-file="handleFileChange" @submit="saveProfile" />
                    
                    <AccountSecurity v-show="activeTab === 'security'" :form="passwordForm" @submit="updatePassword" />
                    
                    <div v-show="activeTab === 'billing'" class="animate-fade-in-up text-center py-10">
                        <div class="w-16 h-16 bg-amber-100 text-amber-500 rounded-2xl flex items-center justify-center mx-auto mb-4">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                        </div>
                        <h3 class="text-lg font-black text-slate-800 uppercase tracking-widest">Paket Enterprise Active</h3>
                        <p class="text-xs font-bold text-slate-400 mt-2">Masa aktif berlangganan Anda sampai dengan 31 Desember 2026.</p>
                    </div>
                </div>

            </div>
        </div>
    </Sidebar>
</template>