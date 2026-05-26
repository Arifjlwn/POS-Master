<script setup>
defineProps({ supplierName: String, noFaktur: String });
defineEmits(['update:supplierName', 'update:noFaktur']);
</script>

<template>
    <div class="bg-white p-8 rounded-[32px] shadow-xl shadow-slate-200/50 border border-slate-100">
        <h3 class="font-black text-slate-800 text-[10px] uppercase tracking-[0.2em] mb-6 flex items-center gap-2">Data Pengiriman</h3>
        <div class="space-y-4">
            <input :value="supplierName" @input="$emit('update:supplierName', $event.target.value)" type="text" placeholder="Nama Supplier..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm focus:border-blue-500 transition-all">
            <input :value="noFaktur" @input="$emit('update:noFaktur', $event.target.value)" type="text" placeholder="No Faktur..." class="w-full p-4 bg-slate-50 border-2 border-slate-100 rounded-2xl outline-none font-bold text-sm focus:border-blue-500 uppercase transition-all">
        </div>
    </div>
</template>