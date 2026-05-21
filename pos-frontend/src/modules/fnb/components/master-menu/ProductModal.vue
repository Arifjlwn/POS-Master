<script setup>
defineProps({
    open: Boolean,
    form: Object,
    loading: Boolean
})

const emit = defineEmits([
    'close',
    'submit'
])
</script>

<template>
    <div
        v-if="open"
        class="fixed inset-0 z-[100] flex items-center justify-center px-4"
    >
        <div
            @click="emit('close')"
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
                    @click="emit('close')"
                    class="flex items-center justify-center w-10 h-10 rounded-2xl bg-slate-100"
                >
                    ✕
                </button>
            </div>

            <div class="p-6">
                <form
                    @submit.prevent="emit('submit')"
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
                            <option value="makanan">
                                Makanan
                            </option>

                            <option value="minuman">
                                Minuman
                            </option>

                            <option value="paket">
                                Paket
                            </option>

                            <option value="snack">
                                Snack
                            </option>

                            <option value="dessert">
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
                        :disabled="loading"
                        class="w-full py-5 text-xs font-black tracking-widest text-white uppercase transition-all rounded-2xl bg-gradient-to-r from-indigo-600 to-indigo-800 active:scale-95 disabled:opacity-60"
                    >
                        {{
                            loading
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
</template>