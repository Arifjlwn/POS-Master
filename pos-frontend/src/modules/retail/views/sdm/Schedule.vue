<script setup>
import { useSchedule } from '../../composables/useSchedule.js';
import Sidebar from '../../components/Sidebar.vue';

// Import Komponen Pintar (Sub-komponen Modular)
import ScheduleHeader from '../../components/schedule/ScheduleHeader.vue';
import ScheduleListMobile from '../../components/schedule/ScheduleListMobile.vue';
import ScheduleTableDesktop from '../../components/schedule/ScheduleTableDesktop.vue';

// Ambil otak & state dari Composable Utama
const {
    isSaving, isLoading, searchQuery, currentUser, 
    mingguJadwal, startDate, endDate, formJadwal, filteredKaryawan,
    canEditSchedule, handleSaveJadwalBulk, handleApproval
} = useSchedule();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <ScheduleHeader 
                v-model:searchQuery="searchQuery"
                :startDate="startDate"
                :endDate="endDate"
                :canEditSchedule="canEditSchedule"
                :isSaving="isSaving"
                @save-jadwal="handleSaveJadwalBulk"
            />

            <div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100 shadow-sm">
                <div class="w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <p class="text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Menarik Data Jadwal...</p>
            </div>
            
            <div v-else-if="filteredKaryawan.length === 0" class="flex flex-col items-center justify-center py-20 bg-white/50 rounded-[32px] border-2 border-dashed border-slate-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-20 h-20 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /><line x1="3" y1="3" x2="21" y2="21" stroke-width="2"/></svg>
                <p class="text-slate-400 font-black text-sm uppercase tracking-widest">Karyawan Tidak Ditemukan</p>
            </div>

            <div v-else>
                <ScheduleListMobile 
                    :filteredKaryawan="filteredKaryawan"
                    :mingguJadwal="mingguJadwal"
                    :formJadwal="formJadwal"
                    :canEditSchedule="canEditSchedule"
                    :currentUser="currentUser"
                    @handle-approval="handleApproval"
                />

                <ScheduleTableDesktop 
                    :filteredKaryawan="filteredKaryawan"
                    :mingguJadwal="mingguJadwal"
                    :formJadwal="formJadwal"
                    :canEditSchedule="canEditSchedule"
                    :currentUser="currentUser"
                    @handle-approval="handleApproval"
                />
            </div>

        </div>
    </Sidebar>
</template>