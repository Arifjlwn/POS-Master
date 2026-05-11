<script setup>
import { Head, useForm } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import { ref } from 'vue';

const props = defineProps({
    karyawan: Array
});

const showModal = ref(false);

const form = useForm({
    name: '',
    email: '',
    password: '',
});

const submit = () => {
    form.post(route('karyawan.store'), {
        onSuccess: () => {
            form.reset();
            showModal.value = false;
        }
    });
};

const deleteKaryawan = (id) => {
    if (confirm('Yakin ingin menghapus akun karyawan ini?')) {
        form.delete(route('karyawan.destroy', id));
    }
};
</script>

<template>
    <Head title="Manajemen Karyawan" />

    <AppLayout>
        <div class="p-6 md:p-8 max-w-7xl mx-auto">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 class="text-3xl font-black text-gray-900 leading-tight">Manajemen Karyawan</h1>
                    <p class="text-gray-500 font-medium">Kelola akses kasir dan tim tokomu di sini.</p>
                </div>
                <button @click="showModal = true" class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-2xl font-black text-sm flex items-center gap-2 shadow-lg shadow-blue-200 transition-all active:scale-95">
                    <span>➕</span> Tambah Karyawan Baru
                </button>
            </div>

            <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
                <table class="w-full text-left">
                    <thead class="bg-gray-50 border-b border-gray-100">
                        <tr>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Nama</th>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Email</th>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Jabatan</th>
                            <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-50">
                        <tr v-if="karyawan.length === 0">
                            <td colspan="4" class="px-6 py-12 text-center text-gray-400 font-medium italic">Belum ada karyawan yang didaftarkan.</td>
                        </tr>
                        <tr v-for="user in karyawan" :key="user.id" class="hover:bg-gray-50 transition-colors">
                            <td class="px-6 py-4">
                                <div class="flex items-center gap-3">
                                    <div class="w-9 h-9 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-black text-xs uppercase">{{ user.name.substring(0,2) }}</div>
                                    <span class="font-bold text-gray-800">{{ user.name }}</span>
                                </div>
                            </td>
                            <td class="px-6 py-4 text-sm text-gray-600 font-medium">{{ user.nik || user.email }}</td>
                            <td class="px-6 py-4">
                                <span class="bg-blue-100 text-blue-700 text-[10px] font-black px-2.5 py-1 rounded-lg uppercase tracking-wider">{{ user.role }}</span>
                            </td>
                            <td class="px-6 py-4 text-right">
                                <button @click="deleteKaryawan(user.id)" class="text-red-400 hover:text-red-600 transition-colors">
                                    <svg class="w-5 h-5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div v-if="showModal" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
            <div class="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden transform transition-all">
                <div class="p-6 border-b border-gray-100 flex justify-between items-center">
                    <h3 class="font-black text-xl text-gray-800">Daftarkan Karyawan</h3>
                    <button @click="showModal = false" class="text-gray-400 hover:text-red-500">✕</button>
                </div>
                <form @submit.prevent="submit" class="p-6 space-y-4">
                    <div>
                        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Nama Lengkap</label>
                        <input v-model="form.name" type="text" required class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-blue-500 focus:border-blue-500 font-bold text-sm">
                    </div>
                    <div>
                        <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Password Sementara</label>
                        <input v-model="form.password" type="password" required class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-blue-500 focus:border-blue-500 font-bold text-sm">
                    </div>
                    <div class="pt-2">
                        <button type="submit" :disabled="form.processing" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-4 rounded-xl font-black text-sm shadow-lg shadow-blue-100 transition-all active:scale-95">
                            {{ form.processing ? 'Menyimpan...' : 'Daftarkan Akun 🚀' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </AppLayout>
</template>
