<script setup>
import { Head, usePage, router } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';
import { ref, watch, onMounted, onUnmounted } from 'vue';

const props = defineProps({
    karyawan: Array,
    riwayat: Array,
    filterTanggal: String,
});

const page = usePage();
const currentUser = page.props.auth.user;

const currentTime = ref('');
let timer;

const updateTime = () => {
    const now = new Date();
    currentTime.value = now.toLocaleTimeString('id-ID', {
        hour: '2-digit', minute: '2-digit', second: '2-digit'
    }).replace(/:/g, '.');
};

onMounted(() => {
    updateTime();
    timer = setInterval(updateTime, 1000);
});

onUnmounted(() => {
    clearInterval(timer);
});

// FITUR FILTER TANGGAL
const tanggalDipilih = ref(props.filterTanggal);

watch(tanggalDipilih, (newDate) => {
    router.get(route('absensi.index'), { tanggal: newDate }, { preserveState: true, preserveScroll: true });
});

// FITUR DOWNLOAD CSV
const downloadLaporan = () => {
    const d = new Date(tanggalDipilih.value);
    const bulan = String(d.getMonth() + 1).padStart(2, '0');
    const tahun = d.getFullYear();
    window.open(`/absensi/export?bulan=${bulan}&tahun=${tahun}`, '_blank');
};

// EKSEKUSI ABSEN
const prosesAbsen = (id, nama, jenis) => {
    if (confirm(`Konfirmasi absen ${jenis} sekarang?`)) {
        router.post(route('absensi.store'), {
            user_id: id,
            jenis: jenis
        }, {
            preserveScroll: true
        });
    }
};
</script>

<template>
    <Head title="Absensi Karyawan" />

    <AppLayout>
        <div class="p-6 md:p-8 max-w-7xl mx-auto">

            <div class="bg-blue-800 rounded-3xl p-6 md:p-8 mb-8 text-white shadow-xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-blue-900">
                <div class="absolute -right-10 -top-20 opacity-10 text-[200px] font-black italic pointer-events-none">⏱️</div>
                <div class="relative z-10 text-center md:text-left mb-6 md:mb-0">
                    <h1 class="text-3xl font-black tracking-tight mb-2">Mesin Absensi</h1>
                    <p class="text-blue-200 font-medium text-sm">Catat jam masuk dan pulang karyawan dengan akurat.</p>
                </div>
                <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner">
                    <div class="text-sm font-bold text-blue-200 uppercase tracking-widest mb-1">Waktu Saat Ini</div>
                    <div class="font-mono text-4xl md:text-5xl font-black text-yellow-300 tracking-wider drop-shadow-md">
                        {{ currentTime || '00.00.00' }}
                    </div>
                </div>
            </div>

            <h2 class="text-lg font-black text-gray-800 mb-4 flex items-center gap-2">
                👥 Daftar Karyawan Toko
            </h2>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-10">
                <div v-for="user in karyawan" :key="user.id"
                     class="bg-white rounded-3xl p-5 shadow-sm border transition-all duration-300"
                     :class="user.id === currentUser.id ? 'border-blue-300 ring-2 ring-blue-100 shadow-md transform hover:-translate-y-1' : 'border-gray-100 opacity-80'">

                    <div class="flex items-center gap-4 mb-5 border-b border-gray-50 pb-4">
                        <div class="w-14 h-14 rounded-full flex items-center justify-center shrink-0 border-2 border-white shadow-md"
                             :class="user.id === currentUser.id ? 'bg-gradient-to-br from-blue-500 to-blue-600 text-white' : 'bg-gray-100 text-gray-400'">
                            <span class="font-black text-lg uppercase">{{ user.name.substring(0, 2) }}</span>
                        </div>
                        <div>
                            <h3 class="font-black text-lg leading-tight" :class="user.id === currentUser.id ? 'text-blue-900' : 'text-gray-600'">
                                {{ user.name }}
                            </h3>
                            <div class="text-[10px] font-bold px-2 py-0.5 mt-1 rounded-md inline-block uppercase tracking-widest"
                                 :class="user.role === 'owner' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'">
                                {{ user.role }} • {{ user.nik || 'OWNER' }}
                            </div>
                        </div>
                    </div>

                    <div class="flex gap-3">
                        <button
                            @click="prosesAbsen(user.id, user.name, 'Masuk')"
                            :disabled="user.id !== currentUser.id"
                            :class="user.id === currentUser.id ? 'bg-green-50 text-green-700 hover:bg-green-600 hover:text-white border-green-200 active:scale-95' : 'bg-gray-50 text-gray-400 border-gray-100 cursor-not-allowed'"
                            class="flex-1 py-3 rounded-xl font-black text-sm transition-all border flex items-center justify-center gap-1.5">
                            <span class="text-lg" :class="{'opacity-50 grayscale': user.id !== currentUser.id}">📥</span> Masuk
                        </button>

                        <button
                            @click="prosesAbsen(user.id, user.name, 'Pulang')"
                            :disabled="user.id !== currentUser.id"
                            :class="user.id === currentUser.id ? 'bg-orange-50 text-orange-700 hover:bg-orange-600 hover:text-white border-orange-200 active:scale-95' : 'bg-gray-50 text-gray-400 border-gray-100 cursor-not-allowed'"
                            class="flex-1 py-3 rounded-xl font-black text-sm transition-all border flex items-center justify-center gap-1.5">
                            <span class="text-lg" :class="{'opacity-50 grayscale': user.id !== currentUser.id}">📤</span> Pulang
                        </button>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 md:p-6 border-b border-gray-100 flex flex-col md:flex-row md:justify-between md:items-center bg-gray-50 gap-4">
                    <h3 class="font-black text-gray-800 text-lg flex items-center gap-2">
                        📋 Log Absensi
                    </h3>

                    <div class="flex items-center gap-3">
                        <div class="relative">
                            <input type="date" v-model="tanggalDipilih" class="px-4 py-2 bg-white border border-gray-200 rounded-xl text-sm font-bold text-gray-700 shadow-sm focus:ring-blue-500 focus:border-blue-500 cursor-pointer">
                        </div>

                        <button v-if="currentUser.role === 'owner'" @click="downloadLaporan" class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-xl font-bold text-sm shadow-md transition-all flex items-center gap-2 active:scale-95">
                            <span>📊</span> Download CSV
                        </button>
                    </div>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-left">
                        <thead class="bg-white border-b border-gray-100">
                            <tr>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Karyawan</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">NIK</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Jam Datang</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Jam Pulang</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-right">Status</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-50">
                            <tr v-if="riwayat.length === 0">
                                <td colspan="5" class="px-6 py-12 text-center text-gray-400 font-medium italic">Belum ada data absensi terekam pada tanggal ini.</td>
                            </tr>
                            <tr v-for="log in riwayat" :key="log.id" class="hover:bg-gray-50 transition-colors">
                                <td class="px-6 py-4 font-bold text-gray-800">{{ log.user?.name }}</td>
                                <td class="px-6 py-4 text-sm font-mono text-gray-500">{{ log.user?.nik || '-' }}</td>
                                <td class="px-6 py-4 text-center">
                                    <span v-if="log.jam_masuk" class="bg-green-100 text-green-700 font-bold px-3 py-1 rounded-lg text-sm">{{ log.jam_masuk }}</span>
                                    <span v-else class="text-gray-300">-</span>
                                </td>
                                <td class="px-6 py-4 text-center">
                                    <span v-if="log.jam_pulang" class="bg-orange-100 text-orange-700 font-bold px-3 py-1 rounded-lg text-sm">{{ log.jam_pulang }}</span>
                                    <span v-else class="text-gray-300">-</span>
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <span class="text-xs font-black px-2.5 py-1 rounded-md uppercase tracking-wider bg-blue-50 text-blue-600 border border-blue-100">
                                        {{ log.status }}
                                    </span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

        </div>
    </AppLayout>
</template>
