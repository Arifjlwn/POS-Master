<script setup>
import { ref, computed, onMounted } from 'vue'
import { useMenuFnB } from '../composables/useMenuFnB.js'
import SidebarFnB from '../components/SidebarFnB.vue'
import ProductModal from '../components/master-menu/ProductModal.vue'
import ProductStats from '../components/master-menu/ProductStats.vue'
import ProductMobileCard from '../components/master-menu/ProductMobileCard.vue'

// 1. Tarik semua amunisi dari "otak" composables
const { 
    products, 
    isLoading, 
    isSubmitting, 
    fetchProducts, 
    saveMenu, 
    deleteMenu, 
    toggleStatus 
} = useMenuFnB()

// 2. State Lokal (khusus buat urusan UI di halaman ini aja)
const userRole = ref('')
const search = ref('')
const selectedCategory = ref('all')
const isModalOpen = ref(false)

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

// 3. Computed (Kalkulasi UI)
const filteredProducts = computed(() => {
    return products.value.filter((p) => {
        const matchSearch = p.nama?.toLowerCase().includes(search.value.toLowerCase())
        const matchCategory = selectedCategory.value === 'all' || p.kategori === selectedCategory.value
        return matchSearch && matchCategory
    })
})

const totalMenu = computed(() => products.value.length)
const totalAvailable = computed(() => products.value.filter(p => p.is_available).length)
const totalUnavailable = computed(() => products.value.filter(p => !p.is_available).length)

// 4. Helper & Aksi Modal
const formatRupiah = (angka) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0)

const resetForm = () => {
    form.value = { id: null, nama: '', harga: '', kategori: 'makanan', gambar: '' }
}

const openAddModal = () => {
    resetForm()
    isModalOpen.value = true
}

const openEditModal = (prod) => {
    form.value = { ...prod } // Copy data produk ke form
    isModalOpen.value = true
}

const closeModal = () => {
    isModalOpen.value = false
    resetForm()
}

const handleFormSubmit = async () => {
    // Panggil fungsi saveMenu dari composables, kalau true (sukses) baru tutup modal
    const success = await saveMenu(form.value)
    if (success) closeModal()
}

// 5. Lifecycle
onMounted(() => {
    userRole.value = localStorage.getItem('role') || ''
    fetchProducts()
})
</script>

<template>
    <SidebarFnB>
        <div class="flex flex-col h-screen bg-slate-50 overflow-hidden font-sans">
            
            <div class="px-5 py-6 bg-white border-b border-slate-200 shrink-0 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 class="text-2xl font-black uppercase tracking-tighter text-slate-800">Master Menu</h1>
                    <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">Smart Product Management</p>
                </div>
                
                <div class="flex flex-col md:flex-row w-full md:w-auto gap-3">
                    <input v-model="search" placeholder="Cari menu..." class="flex-1 md:w-64 px-4 py-3 md:py-2.5 bg-slate-100 rounded-xl text-xs font-bold outline-none border-2 border-transparent focus:border-indigo-500">
                    
                    <select v-model="selectedCategory" class="px-4 py-3 md:py-2.5 bg-slate-100 rounded-xl text-xs font-bold outline-none border-2 border-transparent focus:border-indigo-500 cursor-pointer">
                        <option v-for="cat in categoryOptions" :key="cat.value" :value="cat.value">{{ cat.label }}</option>
                    </select>

                    <button v-if="userRole === 'owner'" @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-3 md:py-2.5 rounded-xl text-xs font-black uppercase active:scale-95 transition-all">
                        Tambah Menu
                    </button>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto p-4 md:p-6 custom-scrollbar relative">
                
                <ProductStats :total-menu="totalMenu" :total-available="totalAvailable" :total-unavailable="totalUnavailable" class="mb-6" />

                <div v-if="isLoading" class="flex flex-col items-center justify-center py-20">
                    <div class="w-10 h-10 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"></div>
                    <p class="mt-4 text-xs font-black text-slate-400 uppercase tracking-widest">Memuat Katalog...</p>
                </div>

                <div v-else-if="filteredProducts.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-400">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mb-4 opacity-30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 2v20M17 5H9.5a4.5 4.5 0 0 0 0 9H15a4.5 4.5 0 0 1 0 9H6.5"/></svg>
                    <p class="text-xs font-black uppercase tracking-widest text-slate-500">Menu Tidak Ditemukan</p>
                </div>

                <div v-else>
                    <div class="hidden lg:block bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                        <table class="w-full text-left">
                            <thead class="bg-slate-50 border-b border-slate-200">
                                <tr>
                                    <th class="p-4 text-[10px] font-black uppercase tracking-widest text-slate-400">Menu</th>
                                    <th class="p-4 text-[10px] font-black uppercase tracking-widest text-slate-400">Harga</th>
                                    <th class="p-4 text-[10px] font-black uppercase tracking-widest text-slate-400 text-center">Status Jualan</th>
                                    <th class="p-4 text-[10px] font-black uppercase tracking-widest text-slate-400 text-center">Aksi</th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-slate-100">
                                <tr v-for="prod in filteredProducts" :key="prod.id" class="hover:bg-slate-50/50 transition-colors">
                                    <td class="p-4 flex items-center gap-4">
                                        <img :src="prod.gambar" class="w-12 h-12 rounded-xl object-cover border border-slate-200 shadow-sm" :class="!prod.is_available ? 'opacity-50 grayscale' : ''">
                                        <div>
                                            <span class="block text-sm font-black text-slate-800 uppercase">{{ prod.nama }}</span>
                                            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">{{ prod.kategori }}</span>
                                        </div>
                                    </td>
                                    <td class="p-4 text-sm font-black text-indigo-600">{{ formatRupiah(prod.harga) }}</td>
                                    <td class="p-4 text-center">
                                        <button @click="toggleStatus(prod)" class="px-3 py-1.5 rounded-lg text-[9px] font-black uppercase tracking-widest active:scale-95 transition-all" :class="prod.is_available ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200' : 'bg-rose-100 text-rose-700 hover:bg-rose-200'">
                                            {{ prod.is_available ? 'Tersedia' : 'Habis' }}
                                        </button>
                                    </td>
                                    <td class="p-4">
                                        <div class="flex items-center justify-center gap-2" v-if="userRole === 'owner'">
                                            <button @click="openEditModal(prod)" class="p-2 bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white rounded-lg transition-colors active:scale-95">✏️</button>
                                            <button @click="deleteMenu(prod.id)" class="p-2 bg-rose-50 text-rose-600 hover:bg-rose-600 hover:text-white rounded-lg transition-colors active:scale-95">🗑️</button>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <ProductMobileCard 
                        class="lg:hidden"
                        :products="filteredProducts" 
                        :userRole="userRole" 
                        :formatRupiah="formatRupiah"
                        @toggle="toggleStatus"
                        @edit="openEditModal"
                        @delete="deleteMenu"
                    />
                </div>
            </div>
            
            <button v-if="userRole === 'owner'" @click="openAddModal" class="lg:hidden fixed bottom-6 right-6 w-14 h-14 bg-indigo-600 text-white rounded-full shadow-xl flex items-center justify-center active:scale-90 transition-transform z-40">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
            </button>
        </div>

        <ProductModal 
            :open="isModalOpen" 
            :form="form" 
            :loading="isSubmitting" 
            @close="closeModal" 
            @submit="handleFormSubmit" 
        />
    </SidebarFnB>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }

/* DEDEK HAPUS SEMUA MEDIA QUERY YANG RIBET! */
</style>