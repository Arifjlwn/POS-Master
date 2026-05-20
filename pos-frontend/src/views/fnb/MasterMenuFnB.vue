<script setup>
import { ref, onMounted, computed } from 'vue'
import SidebarFnB from './SidebarFnB.vue'
import api from '../../api.js'
import Swal from 'sweetalert2'

const products = ref([])
const isLoading = ref(false)
const isModalOpen = ref(false)
const isSubmitting = ref(false)

const userRole = ref('')
const search = ref('')
const selectedCategory = ref('all')

const form = ref({
    id: null,
    nama: '',
    harga: '',
    kategori: 'makanan',
    gambar: ''
})

const categoryOptions = [
    { value: 'all', label: 'Semua' },
    { value: 'makanan', label: 'Makanan' },
    { value: 'minuman', label: 'Minuman' },
    { value: 'paket', label: 'Paket' },
    { value: 'snack', label: 'Snack' },
    { value: 'dessert', label: 'Dessert' }
]

const defaultImage =
    'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?auto=format&fit=crop&w=1200&q=80'

const formatRupiah = (angka) =>
    new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        maximumFractionDigits: 0
    }).format(angka || 0)

const filteredProducts = computed(() => {
    return products.value.filter((p) => {
        const matchSearch = p.nama
            ?.toLowerCase()
            .includes(search.value.toLowerCase())

        const matchCategory =
            selectedCategory.value === 'all'
                ? true
                : p.kategori === selectedCategory.value

        return matchSearch && matchCategory
    })
})

const totalMenu = computed(() => products.value.length)

const totalAvailable = computed(() =>
    products.value.filter((p) => p.is_available).length
)

const totalUnavailable = computed(() =>
    products.value.filter((p) => !p.is_available).length
)

const fetchProducts = async () => {
    isLoading.value = true

    try {
        const response = await api.get('/fnb/products')

        products.value = response.data.map((p) => ({
            ...p,
            nama: p.nama_produk,
            harga: p.harga_jual,
            kategori: p.kategori,
            gambar: p.gambar || defaultImage
        }))
    } catch (error) {
        Swal.fire(
            'Gagal!',
            'Tidak dapat memuat data menu.',
            'error'
        )
    } finally {
        isLoading.value = false
    }
}

onMounted(() => {
    userRole.value = localStorage.getItem('role') || ''
    fetchProducts()
})

const resetForm = () => {
    form.value = {
        id: null,
        nama: '',
        harga: '',
        kategori: 'makanan',
        gambar: ''
    }
}

const openAddModal = () => {
    resetForm()
    isModalOpen.value = true
}

const closeModal = () => {
    isModalOpen.value = false
    resetForm()
}

const submitMenu = async () => {
    if (userRole.value !== 'owner') {
        return Swal.fire(
            'Akses Ditolak!',
            'Hanya Owner yang dapat mengubah menu.',
            'error'
        )
    }

    isSubmitting.value = true

    try {
        const payload = {
            nama: form.value.nama,
            harga: Number(form.value.harga),
            kategori: form.value.kategori,
            gambar: form.value.gambar || defaultImage
        }

        if (form.value.id) {
            await api.put(
                `/fnb/products/${form.value.id}`,
                payload
            )

            Swal.fire(
                'Berhasil!',
                'Menu berhasil diperbarui.',
                'success'
            )
        } else {
            await api.post('/fnb/products', payload)

            Swal.fire(
                'Berhasil!',
                'Menu berhasil ditambahkan.',
                'success'
            )
        }

        closeModal()
        fetchProducts()
    } catch (error) {
        Swal.fire(
            'Error!',
            'Gagal menyimpan menu.',
            'error'
        )
    } finally {
        isSubmitting.value = false
    }
}

const editMenu = (prod) => {
    if (userRole.value !== 'owner') {
        return Swal.fire(
            'Ditolak!',
            'Hanya Owner.',
            'error'
        )
    }

    form.value = {
        id: prod.id,
        nama: prod.nama,
        harga: prod.harga,
        kategori: prod.kategori,
        gambar: prod.gambar
    }

    isModalOpen.value = true
}

const deleteMenu = async (id) => {
    if (userRole.value !== 'owner') {
        return Swal.fire(
            'Ditolak!',
            'Hanya Owner.',
            'error'
        )
    }

    const result = await Swal.fire({
        title: 'Hapus Menu?',
        text: 'Data menu akan dihapus permanen.',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Ya, Hapus',
        cancelButtonText: 'Batal',
        confirmButtonColor: '#ef4444',
        cancelButtonColor: '#64748b'
    })

    if (!result.isConfirmed) return

    try {
        await api.delete(`/fnb/products/${id}`)

        Swal.fire(
            'Berhasil!',
            'Menu berhasil dihapus.',
            'success'
        )

        fetchProducts()
    } catch (error) {
        Swal.fire(
            'Error!',
            'Gagal menghapus menu.',
            'error'
        )
    }
}

const toggleStatus = async (product) => {
    try {
        const res = await api.put(
            `/fnb/products/${product.id}/toggle`
        )

        Swal.fire({
            toast: true,
            position: 'top-end',
            icon: res.data.is_available
                ? 'success'
                : 'warning',
            title: res.data.message,
            showConfirmButton: false,
            timer: 1500
        })

        fetchProducts()
    } catch (error) {
        Swal.fire(
            'Gagal!',
            'Tidak dapat mengubah status.',
            'error'
        )
    }
}
</script>

<template>
    <SidebarFnB>
        <div
            class="relative flex flex-col h-screen overflow-hidden bg-gradient-to-br from-slate-50 via-white to-indigo-50"
        >
            <!-- BG -->
            <div
                class="absolute inset-0 overflow-hidden pointer-events-none"
            >
                <div
                    class="absolute -top-40 -right-40 w-[400px] h-[400px] rounded-full bg-indigo-300/20 blur-3xl"
                ></div>

                <div
                    class="absolute bottom-0 left-0 w-[350px] h-[350px] rounded-full bg-cyan-300/20 blur-3xl"
                ></div>
            </div>

            <!-- HEADER -->
            <div
                class="sticky top-0 z-30 border-b backdrop-blur-xl bg-white/70 border-white/40"
            >
                <div
                    class="flex flex-col gap-5 px-5 py-5 lg:px-8 lg:flex-row lg:items-center lg:justify-between"
                >
                    <!-- LEFT -->
                    <div class="flex items-center gap-4">
                        <div
                            class="flex items-center justify-center w-14 h-14 text-white rounded-3xl bg-gradient-to-br from-indigo-600 to-indigo-800 shadow-[0_15px_40px_rgba(79,70,229,0.35)]"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="w-7 h-7"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                stroke-width="2.5"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M3 7h18M5 7l1 12h12l1-12"
                                />
                            </svg>
                        </div>

                        <div>
                            <h1
                                class="text-3xl font-black tracking-tight text-slate-800"
                            >
                                Master Menu
                            </h1>

                            <p
                                class="mt-1 text-[11px] uppercase tracking-[0.25em] font-black text-slate-400"
                            >
                                Smart Product Management
                            </p>
                        </div>
                    </div>

                    <!-- RIGHT -->
                    <div
                        class="flex flex-col w-full gap-3 lg:w-auto lg:flex-row"
                    >
                        <!-- SEARCH -->
                        <div class="relative w-full lg:w-80">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="absolute w-5 h-5 -translate-y-1/2 left-4 top-1/2 text-slate-400"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="m21 21-4.35-4.35"
                                />

                                <circle
                                    cx="11"
                                    cy="11"
                                    r="6"
                                />
                            </svg>

                            <input
                                v-model="search"
                                type="text"
                                placeholder="Cari menu..."
                                class="w-full h-12 pl-12 pr-4 text-sm font-bold transition-all border outline-none rounded-3xl bg-white/80 border-white/40 text-slate-700 focus:border-indigo-500"
                            />
                        </div>

                        <!-- FILTER -->
                        <select
                            v-model="selectedCategory"
                            class="h-12 px-4 text-sm font-black transition-all border outline-none rounded-2xl bg-white/80 border-white/40 text-slate-700 focus:border-indigo-500"
                        >
                            <option
                                v-for="cat in categoryOptions"
                                :key="cat.value"
                                :value="cat.value"
                            >
                                {{ cat.label }}
                            </option>
                        </select>

                        <!-- BTN -->
                        <button
                            v-if="userRole === 'owner'"
                            @click="openAddModal"
                            class="hidden lg:flex items-center justify-center gap-3 px-6 h-14 text-xs font-black tracking-widest text-white uppercase transition-all rounded-2xl bg-gradient-to-r from-indigo-600 to-indigo-800 shadow-[0_15px_40px_rgba(79,70,229,0.35)] hover:scale-[1.02] active:scale-95"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="w-5 h-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                stroke-width="3"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M12 4v16m8-8H4"
                                />
                            </svg>

                            Tambah Menu
                        </button>
                    </div>
                </div>
            </div>

            <!-- CONTENT -->
            <div
                class="relative flex-1 p-5 overflow-y-auto custom-scrollbar lg:p-8"
            >
                <!-- STATS -->
<div
    class="grid grid-cols-3 gap-2 mb-5 md:gap-4 lg:gap-5"
>
    <!-- TOTAL -->
    <div
        class="mobile-stat p-3 md:p-5 border rounded-[20px] md:rounded-[28px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_10px_40px_rgba(15,23,42,0.06)]"
    >
        <p
            class="text-[8px] md:text-[11px] font-black uppercase tracking-[0.18em] md:tracking-[0.25em] text-slate-400"
        >
            Total Menu
        </p>

        <h2
            class="mt-1 md:mt-2 text-[22px] md:text-4xl font-black text-slate-800 leading-none"
        >
            {{ totalMenu }}
        </h2>
    </div>

    <!-- READY -->
    <div
        class="mobile-stat p-3 md:p-5 border rounded-[20px] md:rounded-[28px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_10px_40px_rgba(15,23,42,0.06)]"
    >
        <p
            class="text-[8px] md:text-[11px] font-black uppercase tracking-[0.18em] md:tracking-[0.25em] text-slate-400"
        >
            Menu Ready
        </p>

        <h2
            class="mt-1 md:mt-2 text-[22px] md:text-4xl font-black text-emerald-600 leading-none"
        >
            {{ totalAvailable }}
        </h2>
    </div>

    <!-- HABIS -->
    <div
        class="mobile-stat p-3 md:p-5 border rounded-[20px] md:rounded-[28px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_10px_40px_rgba(15,23,42,0.06)]"
    >
        <p
            class="text-[8px] md:text-[11px] font-black uppercase tracking-[0.18em] md:tracking-[0.25em] text-slate-400"
        >
            Menu Habis
        </p>

        <h2
            class="mt-1 md:mt-2 text-[22px] md:text-4xl font-black text-rose-600 leading-none"
        >
            {{ totalUnavailable }}
        </h2>
    </div>
</div>

                <!-- LOADING -->
                <div
                    v-if="isLoading"
                    class="flex flex-col items-center justify-center h-[60vh]"
                >
                    <div
                        class="w-12 h-12 border-4 rounded-full border-slate-200 border-t-indigo-600 animate-spin"
                    ></div>

                    <p
                        class="mt-5 text-xs font-black tracking-[0.25em] uppercase text-slate-400"
                    >
                        Memuat Data...
                    </p>
                </div>

                <!-- EMPTY -->
                <div
                    v-else-if="filteredProducts.length === 0"
                    class="flex flex-col items-center justify-center h-[60vh] p-10 border rounded-[40px] bg-white/70 backdrop-blur-xl border-white/40"
                >
                    <div
                        class="flex items-center justify-center w-20 h-20 rounded-full bg-slate-100"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-12 h-12 text-slate-300"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="1.5"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M20 13V7a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v6"
                            />
                        </svg>
                    </div>

                    <h2
                        class="mt-6 text-2xl font-black text-slate-800"
                    >
                        Menu Tidak Ditemukan
                    </h2>

                    <p
                        class="mt-2 text-sm font-semibold text-center text-slate-400"
                    >
                        Coba gunakan pencarian lain.
                    </p>
                </div>

                <!-- DESKTOP -->
                <div
                    v-if="filteredProducts.length > 0"
                    class="hidden overflow-hidden border lg:block rounded-[32px] bg-white/70 backdrop-blur-xl border-white/40 shadow-[0_20px_60px_rgba(15,23,42,0.08)]"
                >
                    <div class="overflow-x-auto">
                        <table class="w-full">
                            <thead>
                                <tr
                                    class="border-b bg-slate-50/80 border-slate-100"
                                >
                                    <th
                                        class="p-5 text-[11px] text-left font-black uppercase tracking-[0.2em] text-slate-400"
                                    >
                                        Menu
                                    </th>

                                    <th
                                        class="p-5 text-[11px] text-left font-black uppercase tracking-[0.2em] text-slate-400"
                                    >
                                        Kategori
                                    </th>

                                    <th
                                        class="p-5 text-[11px] text-left font-black uppercase tracking-[0.2em] text-slate-400"
                                    >
                                        Harga
                                    </th>

                                    <th
                                        class="p-5 text-[11px] text-center font-black uppercase tracking-[0.2em] text-slate-400"
                                    >
                                        Status
                                    </th>

                                    <th
                                        class="p-5 text-[11px] text-center font-black uppercase tracking-[0.2em] text-slate-400"
                                    >
                                        Action
                                    </th>
                                </tr>
                            </thead>

                            <tbody
                                class="divide-y divide-slate-100"
                            >
                                <tr
                                    v-for="prod in filteredProducts"
                                    :key="prod.id"
                                    class="transition-all hover:bg-white/60"
                                >
                                    <td class="p-5">
                                        <div
                                            class="flex items-center gap-4"
                                        >
                                            <img
                                                :src="prod.gambar"
                                                class="object-cover w-14 h-14 rounded-2xl"
                                            />

                                            <div>
                                                <h3
                                                    class="text-sm font-black text-slate-800"
                                                >
                                                    {{ prod.nama }}
                                                </h3>

                                                <p
                                                    class="mt-1 text-[10px] uppercase tracking-[0.2em] font-bold text-slate-400"
                                                >
                                                    ID #{{
                                                        prod.id
                                                    }}
                                                </p>
                                            </div>
                                        </div>
                                    </td>

                                    <td class="p-5">
                                        <span
                                            class="px-3 py-2 text-[10px] uppercase rounded-xl font-black tracking-[0.2em] bg-slate-100 text-slate-600"
                                        >
                                            {{ prod.kategori }}
                                        </span>
                                    </td>

                                    <td
                                        class="p-5 text-sm font-black text-indigo-600"
                                    >
                                        {{
                                            formatRupiah(
                                                prod.harga
                                            )
                                        }}
                                    </td>

                                    <td class="p-5 text-center">
                                        <button
                                            @click="
                                                toggleStatus(
                                                    prod
                                                )
                                            "
                                            class="relative inline-flex items-center h-8 transition-all w-14 rounded-full"
                                            :class="
                                                prod.is_available
                                                    ? 'bg-emerald-500'
                                                    : 'bg-slate-300'
                                            "
                                        >
                                            <span
                                                class="inline-block w-6 h-6 transition-all bg-white rounded-full"
                                                :class="
                                                    prod.is_available
                                                        ? 'translate-x-7'
                                                        : 'translate-x-1'
                                                "
                                            ></span>
                                        </button>
                                    </td>

                                    <td class="p-5">
                                        <div
                                            class="flex items-center justify-center gap-2"
                                        >
                                            <button
                                                v-if="
                                                    userRole ===
                                                    'owner'
                                                "
                                                @click="
                                                    editMenu(
                                                        prod
                                                    )
                                                "
                                                class="flex items-center justify-center w-11 h-11 transition-all rounded-2xl bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white active:scale-95"
                                            >
                                                ✏️
                                            </button>

                                            <button
                                                v-if="
                                                    userRole ===
                                                    'owner'
                                                "
                                                @click="
                                                    deleteMenu(
                                                        prod.id
                                                    )
                                                "
                                                class="flex items-center justify-center w-11 h-11 transition-all rounded-2xl bg-rose-50 text-rose-600 hover:bg-rose-600 hover:text-white active:scale-95"
                                            >
                                                🗑️
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>

                <!-- MOBILE CARD -->
                <div
                    v-if="filteredProducts.length > 0"
                    class="grid grid-cols-1 gap-3 lg:hidden"
                >

                    <div
                        v-for="prod in filteredProducts"
                        :key="prod.id"
                        class="overflow-hidden border border-white/50 bg-white/80 backdrop-blur-xl rounded-[26px] shadow-[0_10px_30px_rgba(15,23,42,0.06)]"
                    >

                        <div class="flex gap-3 p-3">

                            <!-- IMAGE -->
                            <div class="relative shrink-0">

                                <img
                                    :src="prod.gambar"
                                    class="object-cover w-20 h-20 shadow-md rounded-2xl"
                                >

                                <div
                                    class="absolute -bottom-1 left-1/2 -translate-x-1/2 px-2 py-[3px] rounded-full text-[8px] font-black tracking-widest text-white"
                                    :class="prod.is_available
                                        ? 'bg-emerald-500'
                                        : 'bg-rose-500'"
                                >
                                    {{ prod.is_available ? 'READY' : 'HABIS' }}
                                </div>

                            </div>

                            <!-- CONTENT -->
                            <div class="flex flex-col justify-between flex-1 min-w-0">

                                <div>

                                    <div class="flex items-start justify-between gap-2">

                                        <div class="min-w-0">

                                            <h2 class="truncate text-[15px] font-black text-slate-800">
                                                {{ prod.nama }}
                                            </h2>

                                            <p class="mt-1 text-[10px] uppercase tracking-[0.2em] text-slate-400 font-bold">
                                                {{ prod.kategori }}
                                            </p>

                                        </div>

                                        <!-- TOGGLE -->
                                        <button
                                            @click="toggleStatus(prod)"
                                            class="relative inline-flex items-center w-11 h-6 transition-all rounded-full shrink-0"
                                            :class="prod.is_available
                                                ? 'bg-emerald-500'
                                                : 'bg-slate-300'"
                                        >
                                            <span
                                                class="inline-block w-4 h-4 transition-all bg-white rounded-full shadow-md"
                                                :class="prod.is_available
                                                    ? 'translate-x-6'
                                                    : 'translate-x-1'"
                                            />
                                        </button>

                                    </div>

                                    <p class="mt-3 text-2xl font-black text-indigo-600">
                                        {{ formatRupiah(prod.harga) }}
                                    </p>

                                </div>

                                <!-- ACTION -->
                                <div
                                    v-if="userRole === 'owner'"
                                    class="flex gap-2 mt-3"
                                >

                                    <button
                                        @click="editMenu(prod)"
                                        class="flex items-center justify-center flex-1 h-10 gap-2 text-[10px] font-black tracking-widest uppercase transition-all bg-blue-50 text-blue-600 rounded-2xl active:scale-95"
                                    >
                                        Edit
                                    </button>

                                    <button
                                        @click="deleteMenu(prod.id)"
                                        class="flex items-center justify-center flex-1 h-10 gap-2 text-[10px] font-black tracking-widest uppercase transition-all rounded-2xl bg-rose-50 text-rose-600 active:scale-95"
                                    >
                                        Hapus
                                    </button>

                                </div>

                            </div>

                        </div>

                    </div>

                </div>
            </div>

            <!-- FAB -->
            <button
                v-if="userRole === 'owner'"
                @click="openAddModal"
                class="fixed z-50 flex items-center justify-center text-white lg:hidden bottom-6 right-6 w-14 h-14 rounded-3xl bg-gradient-to-br from-indigo-600 to-indigo-800 shadow-[0_15px_40px_rgba(79,70,229,0.4)] active:scale-90"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="w-7 h-7"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M12 4v16m8-8H4"
                    />
                </svg>
            </button>

            <!-- MODAL -->
            <div
                v-if="isModalOpen"
                class="fixed inset-0 z-[100] flex items-center justify-center px-4"
            >
                <div
                    @click="closeModal"
                    class="absolute inset-0 bg-slate-900/60 backdrop-blur-md"
                ></div>

                <div
                    class="relative z-10 w-full max-w-md overflow-hidden bg-white border shadow-2xl rounded-[36px] border-white/40"
                >
                    <div
                        class="flex items-center justify-between p-6 border-b border-slate-100"
                    >
                        <div>
                            <h2
                                class="text-xl font-black text-slate-800"
                            >
                                {{
                                    form.id
                                        ? 'Edit Menu'
                                        : 'Tambah Menu'
                                }}
                            </h2>

                            <p
                                class="mt-1 text-[11px] uppercase tracking-[0.2em] font-bold text-slate-400"
                            >
                                Product Information
                            </p>
                        </div>

                        <button
                            @click="closeModal"
                            class="flex items-center justify-center w-10 h-10 rounded-2xl bg-slate-100"
                        >
                            ✕
                        </button>
                    </div>

                    <div class="p-6">
                        <form
                            @submit.prevent="submitMenu"
                            class="space-y-5"
                        >
                            <div>
                                <label
                                    class="block mb-2 ml-1 text-[11px] font-black uppercase tracking-[0.2em] text-slate-400"
                                >
                                    Nama Menu
                                </label>

                                <input
                                    v-model="form.nama"
                                    type="text"
                                    required
                                    class="w-full px-5 py-4 text-sm font-bold border outline-none rounded-2xl border-slate-200 focus:border-indigo-500"
                                />
                            </div>

                            <div>
                                <label
                                    class="block mb-2 ml-1 text-[11px] font-black uppercase tracking-[0.2em] text-slate-400"
                                >
                                    Kategori
                                </label>

                                <select
                                    v-model="form.kategori"
                                    class="w-full px-5 py-4 text-sm font-bold border outline-none rounded-2xl border-slate-200 focus:border-indigo-500"
                                >
                                    <option
                                        value="makanan"
                                    >
                                        Makanan
                                    </option>

                                    <option
                                        value="minuman"
                                    >
                                        Minuman
                                    </option>

                                    <option value="paket">
                                        Paket
                                    </option>

                                    <option value="snack">
                                        Snack
                                    </option>

                                    <option
                                        value="dessert"
                                    >
                                        Dessert
                                    </option>
                                </select>
                            </div>

                            <div>
                                <label
                                    class="block mb-2 ml-1 text-[11px] font-black uppercase tracking-[0.2em] text-slate-400"
                                >
                                    Harga
                                </label>

                                <input
                                    v-model="form.harga"
                                    type="number"
                                    required
                                    class="w-full px-5 py-4 text-sm font-black border outline-none rounded-2xl border-slate-200 text-indigo-600 focus:border-indigo-500"
                                />
                            </div>

                            <div>
                                <label
                                    class="block mb-2 ml-1 text-[11px] font-black uppercase tracking-[0.2em] text-slate-400"
                                >
                                    Link Gambar
                                </label>

                                <input
                                    v-model="form.gambar"
                                    type="url"
                                    class="w-full px-5 py-4 text-sm font-bold border outline-none rounded-2xl border-slate-200 focus:border-indigo-500"
                                />
                            </div>

                            <button
                                type="submit"
                                :disabled="isSubmitting"
                                class="w-full py-5 text-xs font-black tracking-widest text-white uppercase transition-all rounded-2xl bg-gradient-to-r from-indigo-600 to-indigo-800 active:scale-95 disabled:opacity-60"
                            >
                                {{
                                    isSubmitting
                                        ? 'Menyimpan...'
                                        : form.id
                                          ? 'Update Menu'
                                          : 'Simpan Menu'
                                }}
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </SidebarFnB>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
    background: linear-gradient(
        to bottom,
        #6366f1,
        #4f46e5
    );
    border-radius: 999px;
}

.custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
}

* {
    -webkit-tap-highlight-color: transparent;
    box-sizing: border-box;
}

button {
    user-select: none;
}

input,
select,
button {
    transition:
        0.25s ease,
        transform 0.2s ease;
}

/* =========================================
   RESPONSIVE GLOBAL
========================================= */

@media (max-width: 1280px) {
    .desktop-table-wrap {
        overflow-x: auto;
    }
}

@media (max-width: 1024px) {
    h1 {
        line-height: 1.1;
    }
}

/* =========================================
   TABLET
========================================= */

@media (max-width: 768px) {
    input,
    select,
    button,
    textarea {
        font-size: 16px !important;
    }

    .mobile-padding {
        padding-left: 14px !important;
        padding-right: 14px !important;
    }

    .mobile-card {
        border-radius: 22px !important;
    }

    .mobile-title {
        font-size: 24px !important;
    }

    .mobile-subtitle {
        font-size: 9px !important;
        letter-spacing: 0.22em !important;
    }

    .mobile-stat {
        padding: 16px !important;
        border-radius: 22px !important;
    }

    .mobile-stat h2 {
        font-size: 28px !important;
        margin-top: 8px !important;
    }

    .mobile-stat p {
        font-size: 9px !important;
        letter-spacing: 0.2em !important;
    }

    .mobile-search,
    .mobile-select {
        height: 48px !important;
        border-radius: 18px !important;
    }

    .mobile-product-card {
        border-radius: 22px !important;
    }

    .mobile-product-image {
        width: 72px !important;
        height: 72px !important;
        border-radius: 18px !important;
    }

    .mobile-product-name {
        font-size: 14px !important;
        line-height: 1.2 !important;
    }

    .mobile-product-category {
        font-size: 9px !important;
    }

    .mobile-product-price {
        font-size: 24px !important;
        margin-top: 10px !important;
    }

    .mobile-action-btn {
        height: 38px !important;
        border-radius: 14px !important;
        font-size: 10px !important;
    }

    .mobile-modal {
        border-radius: 26px !important;
    }

    .mobile-modal-padding {
        padding: 18px !important;
    }

    .mobile-modal input,
    .mobile-modal select {
        padding: 14px 16px !important;
        border-radius: 16px !important;
    }

    .mobile-modal button[type='submit'] {
        padding: 16px !important;
        border-radius: 16px !important;
    }

    .fab-mobile {
        width: 56px !important;
        height: 56px !important;
        right: 18px !important;
        bottom: 18px !important;
        border-radius: 20px !important;
    }
}

/* =========================================
   HP KECIL
========================================= */

@media (max-width: 480px) {
    .mobile-padding {
        padding-left: 12px !important;
        padding-right: 12px !important;
    }

    .mobile-header {
        gap: 14px !important;
        padding-top: 14px !important;
        padding-bottom: 14px !important;
    }

    .mobile-icon-box {
        width: 50px !important;
        height: 50px !important;
        border-radius: 18px !important;
    }

    .mobile-title {
        font-size: 20px !important;
    }

    .mobile-subtitle {
        font-size: 8px !important;
        letter-spacing: 0.18em !important;
    }

    .mobile-search,
    .mobile-select {
        height: 44px !important;
        border-radius: 16px !important;
        font-size: 13px !important;
    }

    .mobile-search-icon {
        width: 16px !important;
        height: 16px !important;
        left: 14px !important;
    }

    .mobile-search-input {
        padding-left: 40px !important;
        padding-right: 14px !important;
        font-size: 13px !important;
    }

    .mobile-stat {
        padding: 14px !important;
        border-radius: 18px !important;
    }

    .mobile-stat h2 {
        font-size: 24px !important;
    }

    .mobile-stat p {
        font-size: 8px !important;
    }

    .mobile-product-card {
        border-radius: 20px !important;
    }

    .mobile-product-image {
        width: 64px !important;
        height: 64px !important;
        border-radius: 16px !important;
    }

    .mobile-product-name {
        font-size: 13px !important;
    }

    .mobile-product-category {
        font-size: 8px !important;
        margin-top: 3px !important;
    }

    .mobile-product-price {
        font-size: 20px !important;
        line-height: 1 !important;
    }

    .mobile-status-badge {
        font-size: 7px !important;
        padding: 3px 7px !important;
    }

    .mobile-toggle {
        width: 42px !important;
        height: 22px !important;
    }

    .mobile-toggle-dot {
        width: 14px !important;
        height: 14px !important;
    }

    .mobile-action-btn {
        height: 34px !important;
        border-radius: 12px !important;
        font-size: 9px !important;
        letter-spacing: 0.15em !important;
    }

    .mobile-modal {
        border-radius: 22px !important;
    }

    .mobile-modal-padding {
        padding: 16px !important;
    }

    .fab-mobile {
        width: 52px !important;
        height: 52px !important;
        border-radius: 18px !important;
    }
}

/* =========================================
   EXTRA SMALL DEVICE
========================================= */

@media (max-width: 360px) {
    .mobile-title {
        font-size: 18px !important;
    }

    .mobile-product-price {
        font-size: 18px !important;
    }

    .mobile-product-name {
        font-size: 12px !important;
    }

    .mobile-stat h2 {
        font-size: 22px !important;
    }

    .mobile-search,
    .mobile-select {
        height: 42px !important;
    }
}
</style>