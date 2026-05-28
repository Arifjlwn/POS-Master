<script setup>

defineProps({
    filteredKaryawan: Array,
    mingguJadwal: Array,
    formJadwal: Object,
    canEditSchedule: Boolean,
    currentUser: Object
});

const emit = defineEmits(['handle-approval']);
</script>

<template>
    <div class="lg:hidden flex flex-col gap-4">
        <div v-for="emp in filteredKaryawan" :key="emp.id" class="bg-white p-5 rounded-[24px] shadow-sm border border-slate-100 flex flex-col gap-4">
            
            <div class="flex items-center gap-4 pb-4 border-b border-slate-100">
                <div class="w-12 h-12 rounded-[14px] bg-slate-900 text-white font-black text-lg flex items-center justify-center uppercase shadow-md">{{ emp.name.substring(0, 2) }}</div>
                <div class="flex-1 min-w-0">
                    <h3 class="font-black text-base text-slate-800 uppercase truncate">{{ emp.name }}</h3>
                    <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest inline-flex items-center gap-1 mt-1 bg-slate-50 px-2 py-0.5 rounded border border-slate-200">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                        {{ emp.role }}
                    </p>
                </div>
            </div>

            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                <div v-for="d in mingguJadwal" :key="d.tanggal" class="bg-slate-50/80 p-3 rounded-xl border border-slate-100 flex flex-col items-center justify-between">
                    <div class="text-center mb-2">
                        <div class="text-[10px] font-black uppercase text-indigo-600 tracking-widest">{{ d.hari }}</div>
                        <div class="text-[14px] font-black text-slate-800 leading-none mt-0.5">{{ d.tglAngka }}</div>
                    </div>
                    
                    <select 
                        v-model="formJadwal[emp.id || emp.user_id][d.tanggal]"
                        :disabled="!canEditSchedule || formJadwal[emp.id || emp.user_id][d.tanggal].includes('(Approved)')"
                        class="w-full px-1 py-2 text-[9px] font-black uppercase tracking-wider rounded-lg border-2 text-center transition-all outline-none cursor-pointer appearance-none shadow-sm"
                        :class="{
                            'bg-emerald-50 text-emerald-700 border-emerald-200': formJadwal[emp.id || emp.user_id][d.tanggal].includes('Shift 1'),
                            'bg-blue-50 text-blue-700 border-blue-200': formJadwal[emp.id || emp.user_id][d.tanggal].includes('Shift 2'),
                            'bg-purple-50 text-purple-700 border-purple-200': formJadwal[emp.id || emp.user_id][d.tanggal].includes('Middle'),
                            'bg-slate-100 text-slate-400 border-slate-200': formJadwal[emp.id || emp.user_id][d.tanggal] === 'OFF',
                            'bg-slate-800 text-white border-slate-900': formJadwal[emp.id || emp.user_id][d.tanggal].includes('Shift 3')
                        }">
                        <option value="Shift 1">SHIFT 1</option>
                        <option value="Middle">MIDDLE</option>
                        <option value="Shift 2">SHIFT 2</option>
                        <option value="Shift 3">SHIFT 3</option>
                        <option value="OFF">OFF (LIBUR)</option>
                        <slot v-if="formJadwal[emp.id || emp.user_id][d.tanggal]">
                            <option hidden :value="formJadwal[emp.id || emp.user_id][d.tanggal]">
                                {{ formJadwal[emp.id || emp.user_id][d.tanggal].substring(0,7) }}
                            </option>
                        </slot>
                    </select>

                    <div v-if="formJadwal[emp.id || emp.user_id][d.tanggal] !== 'OFF'" class="w-full mt-1.5 flex justify-center">
                        <span v-if="formJadwal[emp.id || emp.user_id][d.tanggal].includes('(Pending)')" class="text-[8px] font-black uppercase text-amber-600 flex items-center gap-1"><svg xmlns="http://www.w3.org/2000/svg" class="w-2 h-2" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="12" r="10"/></svg> Pending</span>
                        <span v-else-if="formJadwal[emp.id || emp.user_id][d.tanggal].includes('(Approved)')" class="text-[8px] font-black uppercase text-emerald-600 flex items-center gap-1"><svg xmlns="http://www.w3.org/2000/svg" class="w-2 h-2" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="12" r="10"/></svg> Approved</span>
                        <span v-else-if="formJadwal[emp.id || emp.user_id][d.tanggal].includes('(Rejected)')" class="text-[8px] font-black uppercase text-red-600 flex items-center gap-1"><svg xmlns="http://www.w3.org/2000/svg" class="w-2 h-2" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="12" r="10"/></svg> Rejected</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
select { -webkit-appearance: none; -moz-appearance: none; appearance: none; }
</style>