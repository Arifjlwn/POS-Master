<script setup>
defineProps({
	searchQuery: String,
	startDate: String,
	endDate: String,
});
const emit = defineEmits(['update:searchQuery', 'update:startDate', 'update:endDate']);
</script>

<template>
	<div class="bg-white p-5 rounded-[24px] border border-slate-200 shadow-sm mb-6 flex flex-col md:flex-row items-center gap-4">
		<div class="relative flex-1 w-full">
			<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
				</svg>
			</div>
			<input :value="searchQuery" @input="emit('update:searchQuery', $event.target.value)" type="text" placeholder="Cari Supplier / No Faktur..." class="w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs text-slate-800 transition-all placeholder:text-slate-400" />
		</div>

		<div class="flex items-center gap-3 w-full md:w-auto">
			<div class="flex-1 md:w-40">
				<input :value="startDate" @input="emit('update:startDate', $event.target.value)" type="date" class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs text-slate-800 transition-all" />
			</div>
			<span class="text-slate-400 font-black text-[10px] uppercase">S/D</span>
			<div class="flex-1 md:w-40">
				<input :value="endDate" @input="emit('update:endDate', $event.target.value)" type="date" class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs text-slate-800 transition-all" />
			</div>
		</div>
	</div>
</template>
