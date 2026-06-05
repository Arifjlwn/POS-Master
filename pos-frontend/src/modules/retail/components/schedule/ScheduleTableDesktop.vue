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
    <div class="hidden lg:block bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden">
        <div class="p-6 border-b border-slate-50 bg-slate-50/50">
            <h3 class="font-black text-slate-800 text-sm flex items-center gap-2 uppercase tracking-widest">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                Matriks Penjadwalan Toko
            </h3>
        </div>

        <div class="overflow-x-auto custom-scrollbar">
            <table class="w-full text-left whitespace-nowrap border-collapse">
                <thead>
                    <tr class="bg-white border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                        <th class="px-6 py-5 sticky left-0 bg-white z-20 border-r border-slate-100 shadow-[4px_0_12px_rgba(0,0,0,0.03)] min-w-[200px]">Nama Karyawan</th>
                        <th v-for="d in mingguJadwal" :key="d.tanggal" class="px-4 py-5 text-center min-w-[160px] border-l border-slate-50">
                            <div class="text-indigo-600 font-black text-xs">{{ d.hari }}</div>
                            <div class="text-slate-400 text-[9px] font-bold mt-1 tracking-widest">{{ d.tanggal }}</div>
                        </th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-50">
                    <tr v-for="emp in filteredKaryawan" :key="emp.id" class="hover:bg-slate-50/50 transition-colors group">
                        <td class="px-6 py-5 sticky left-0 bg-white group-hover:bg-slate-50 z-10 border-r border-slate-100 shadow-[4px_0_12px_rgba(0,0,0,0.03)] transition-colors">
                            <div class="font-black text-slate-800 text-sm uppercase tracking-tighter">{{ emp.name }}</div>
                            <div class="text-[9px] font-bold text-slate-400 mt-1 uppercase tracking-widest">{{ emp.role }}</div>
                        </td>

                        <td v-for="d in mingguJadwal" :key="d.tanggal" class="p-4 text-center border-l border-slate-50">
    <div class="flex flex-col items-center gap-2">
        <select 
            v-model="formJadwal[emp.public_id || emp.id][d.tanggal]"
            :disabled="!canEditSchedule || formJadwal[emp.public_id || emp.id][d.tanggal].includes('(Approved)')"
            class="w-full px-3 py-2.5 text-[11px] font-black uppercase tracking-widest rounded-xl border-2 text-center transition-all outline-none cursor-pointer appearance-none shadow-sm"
            :class="{
                'bg-emerald-50 text-emerald-700 border-emerald-200 focus:border-emerald-500': formJadwal[emp.public_id || emp.id][d.tanggal].includes('Shift 1'),
                'bg-blue-50 text-blue-700 border-blue-200 focus:border-blue-500': formJadwal[emp.public_id || emp.id][d.tanggal].includes('Shift 2'),
                'bg-purple-50 text-purple-700 border-purple-200 focus:border-purple-500': formJadwal[emp.public_id || emp.id][d.tanggal].includes('Middle'),
                'bg-slate-50 text-slate-400 border-slate-200 focus:border-slate-400': formJadwal[emp.public_id || emp.id][d.tanggal] === 'OFF',
                'bg-slate-900 text-white border-slate-800': formJadwal[emp.public_id || emp.id][d.tanggal].includes('Shift 3')
            }">
            <option value="Shift 1">SHIFT 1</option>
            <option value="Middle">MIDDLE</option>
            <option value="Shift 2">SHIFT 2</option>
            <option value="Shift 3">SHIFT 3</option>
            <option value="OFF">LIBUR (OFF)</option>
            <slot v-if="formJadwal[emp.public_id || emp.id][d.tanggal]">
                <option hidden :value="formJadwal[emp.public_id || emp.id][d.tanggal]">
                    {{ formJadwal[emp.public_id || emp.id][d.tanggal] }}
                </option>
            </slot>
        </select>

        <div class="flex justify-center w-full min-h-[24px]" v-if="formJadwal[emp.public_id || emp.id][d.tanggal] !== 'OFF'">
            <span v-if="formJadwal[emp.public_id || emp.id][d.tanggal].includes('(Pending)')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-amber-50 text-amber-600 border border-amber-100 text-[9px] font-black uppercase tracking-widest shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg> Pending
            </span>
            <span v-else-if="formJadwal[emp.public_id || emp.id][d.tanggal].includes('(Approved)')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100 text-[9px] font-black uppercase tracking-widest shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg> Approved
            </span>
            <span v-else-if="formJadwal[emp.public_id || emp.id][d.tanggal].includes('(Rejected)')" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-50 text-red-600 border border-red-100 text-[9px] font-black uppercase tracking-widest shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg> Rejected
            </span>
            <span v-else class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-100 text-slate-500 text-[9px] font-black uppercase tracking-widest border border-slate-200">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg> Draft
            </span>
        </div>

        <div v-if="currentUser.role === 'owner' && formJadwal[emp.public_id || emp.id][d.tanggal].includes('(Pending)')" class="flex gap-2 mt-1 justify-center w-full">
            <button @click="emit('handle-approval', emp.public_id || emp.id, d.tanggal, 'approve')" class="flex-1 p-2 rounded-xl bg-emerald-100 text-emerald-600 hover:bg-emerald-500 hover:text-white transition-all shadow-sm border border-emerald-200 flex items-center justify-center" title="Setujui">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
            </button>
            <button @click="emit('handle-approval', emp.public_id || emp.id, d.tanggal, 'reject')" class="flex-1 p-2 rounded-xl bg-red-100 text-red-600 hover:bg-red-500 hover:text-white transition-all shadow-sm border border-red-200 flex items-center justify-center" title="Tolak">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            </button>
        </div>
    </div>
</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
select { -webkit-appearance: none; -moz-appearance: none; appearance: none; }
</style>