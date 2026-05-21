<script setup>
defineProps({
    products: Array,
    userRole: String,
    formatRupiah: Function
})

const emit = defineEmits([
    'toggle',
    'edit',
    'delete'
])
</script>

<template>
    <div
        class="grid grid-cols-1 gap-3 lg:hidden"
    >
        <div
            v-for="prod in products"
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
                                @click="emit('toggle', prod)"
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
                            @click="emit('edit', prod)"
                            class="flex items-center justify-center flex-1 h-10 gap-2 text-[10px] font-black tracking-widest uppercase transition-all bg-blue-50 text-blue-600 rounded-2xl active:scale-95"
                        >
                            Edit
                        </button>

                        <button
                            @click="emit('delete', prod.id)"
                            class="flex items-center justify-center flex-1 h-10 gap-2 text-[10px] font-black tracking-widest uppercase transition-all rounded-2xl bg-rose-50 text-rose-600 active:scale-95"
                        >
                            Hapus
                        </button>

                    </div>

                </div>

            </div>
        </div>
    </div>
</template>