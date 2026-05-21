<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Swal from 'sweetalert2'

const route = useRoute()
const router = useRouter()

const isMobileMenuOpen = ref(false)

const userName = ref('Kasir')
const userRole = ref('Staff FnB')
const userInitial = ref('KR')

const storeName = ref('Resto POS')
const storeInitial = ref('RP')

const menuItems = [
    {
        name: 'Kasir POS',
        path: '/fnb/kasir',
        color: 'from-indigo-500 to-indigo-700',
        active: 'bg-indigo-600 text-white shadow-indigo-500/30',
        icon: `
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.3">
                <rect x="3" y="3" width="18" height="18" rx="2"/>
                <path d="M3 9h18"/>
                <path d="M9 21V9"/>
            </svg>
        `
    },
    {
        name: 'Monitor Dapur',
        path: '/fnb/dapur',
        color: 'from-orange-500 to-orange-700',
        active: 'bg-orange-500 text-white shadow-orange-500/30',
        icon: `
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.3">
                <path d="M4 8h16"/>
                <path d="M6 8v10a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V8"/>
                <path d="M9 12h6"/>
            </svg>
        `
    },
    {
        name: 'Master Menu',
        path: '/fnb/master-menu',
        color: 'from-blue-500 to-blue-700',
        active: 'bg-blue-600 text-white shadow-blue-500/30',
        icon: `
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.3">
                <path d="M4 7h16"/>
                <path d="M6 7l1 12h10l1-12"/>
                <path d="M9 11h6"/>
            </svg>
        `
    },
    {
        name: 'Laporan',
        path: '/fnb/laporan',
        color: 'from-emerald-500 to-emerald-700',
        active: 'bg-emerald-500 text-white shadow-emerald-500/30',
        icon: `
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.3">
                <path d="M4 19h16"/>
                <path d="M7 15l3-3 3 2 4-5"/>
            </svg>
        `
    }
]

onMounted(() => {
    const name = localStorage.getItem('name') || 'Kasir Resto'

    userName.value = name
    userRole.value =
        localStorage.getItem('role') || 'F&B Outlet'

    userInitial.value = name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .substring(0, 2)
        .toUpperCase()

    const savedStore = localStorage.getItem('storeName')

    if (savedStore) {
        storeName.value = savedStore

        storeInitial.value = savedStore
            .split(' ')
            .map((n) => n[0])
            .join('')
            .substring(0, 2)
            .toUpperCase()
    }
})

const handleLogout = () => {
    Swal.fire({
        title: 'Akhiri Sesi?',
        text: 'Pastikan operasional aman sebelum logout.',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#ef4444',
        cancelButtonColor: '#475569',
        confirmButtonText: 'Keluar'
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.removeItem('token')

            router.push('/login')
        }
    })
}
</script>

<template>
    <div
        class="relative flex w-full h-screen overflow-hidden bg-slate-100"
    >
        <!-- MOBILE OVERLAY -->
        <transition name="fade">
            <div
                v-if="isMobileMenuOpen"
                @click="isMobileMenuOpen = false"
                class="fixed inset-0 z-40 bg-slate-950/50 backdrop-blur-sm lg:hidden"
            ></div>
        </transition>

        <!-- SIDEBAR -->
        <aside
            :class="
                isMobileMenuOpen
                    ? 'translate-x-0'
                    : '-translate-x-full lg:translate-x-0'
            "
            class="fixed lg:relative inset-y-0 left-0 z-50 w-[290px] bg-[#071120] border-r border-white/5 flex flex-col transition-all duration-300 ease-out shadow-2xl lg:shadow-none"
        >
            <!-- BG -->
            <div
                class="absolute inset-0 overflow-hidden pointer-events-none"
            >
                <div
                    class="absolute top-0 left-0 w-[300px] h-[300px] bg-indigo-500/10 blur-3xl rounded-full"
                ></div>

                <div
                    class="absolute bottom-0 right-0 w-[250px] h-[250px] bg-cyan-400/5 blur-3xl rounded-full"
                ></div>
            </div>

            <!-- DESKTOP HEADER -->
<div
    class="relative items-center justify-between hidden px-6 border-b h-20 border-white/5 shrink-0 lg:flex"
>
    <div class="flex items-center gap-4 min-w-0">

        <div
            class="flex items-center justify-center w-14 h-14 text-lg font-black text-white border shadow-2xl rounded-3xl bg-gradient-to-br from-indigo-500 to-indigo-700 border-white/10 shadow-indigo-900/50 shrink-0"
        >
            {{ storeInitial }}
        </div>

        <div class="min-w-0">

            <h1
                class="text-sm font-black tracking-[0.25em] uppercase text-white truncate"
            >
                {{ storeName }}
            </h1>

            <p
                class="mt-1 text-[10px] uppercase tracking-[0.25em] text-slate-500 font-bold"
            >
                Smart FnB Dashboard
            </p>

        </div>
    </div>
</div>

            <!-- NAV -->
            <div class="flex flex-col flex-1 min-h-0 relative">
                <div class="px-6 pt-6 pb-3 lg:pt-6">
                    <p
                        class="text-[10px] uppercase tracking-[0.3em] text-slate-500 font-black"
                    >
                        Main Navigation
                    </p>
                </div>

                <nav
                    class="flex-1 px-4 pb-6 overflow-y-auto custom-scrollbar"
                >
                    <div class="space-y-2">
                        <router-link
                            v-for="item in menuItems"
                            :key="item.path"
                            :to="item.path"
                            @click="isMobileMenuOpen = false"
                            class="group relative flex items-center gap-4 px-4 py-4 rounded-2xl transition-all duration-300 overflow-hidden"
                            :class="
                                route.path === item.path
                                    ? item.active
                                    : 'text-slate-400 hover:bg-white/[0.04] hover:text-white'
                            "
                        >
                            <!-- ACTIVE GLOW -->
                            <div
                                v-if="route.path === item.path"
                                class="absolute inset-0 opacity-20 bg-gradient-to-r"
                                :class="item.color"
                            ></div>

                            <!-- ICON -->
                            <div
                                class="relative z-10 flex items-center justify-center w-11 h-11 rounded-2xl shrink-0"
                                :class="
                                    route.path === item.path
                                        ? 'bg-white/10'
                                        : 'bg-white/[0.03] group-hover:bg-white/10'
                                "
                                v-html="item.icon"
                            ></div>

                            <!-- TEXT -->
                            <div class="relative z-10 min-w-0">
                                <h2
                                    class="text-xs font-black uppercase tracking-widest truncate"
                                >
                                    {{ item.name }}
                                </h2>

                                <p
                                    class="mt-1 text-[9px] uppercase tracking-[0.2em] opacity-60"
                                >
                                    Management
                                </p>
                            </div>
                        </router-link>
                    </div>
                </nav>

                <!-- USER -->
                <div
                    class="p-4 border-t border-white/5 bg-white/[0.02] shrink-0"
                >
                    <div
                        class="flex items-center gap-4 p-4 rounded-3xl bg-white/[0.03] border border-white/5"
                    >
                        <div
                            class="flex items-center justify-center w-14 h-14 text-sm font-black text-white shadow-lg rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 shrink-0"
                        >
                            {{ userInitial }}
                        </div>

                        <div class="min-w-0">
                            <h3
                                class="text-sm font-black tracking-wide text-white truncate uppercase"
                            >
                                {{ userName }}
                            </h3>

                            <p
                                class="mt-1 text-[10px] uppercase tracking-[0.25em] text-slate-500 font-bold truncate"
                            >
                                {{ userRole }}
                            </p>
                        </div>
                    </div>

                    <!-- LOGOUT -->
                    <button
                        @click="handleLogout"
                        class="group mt-4 w-full flex items-center justify-center gap-3 h-14 rounded-2xl bg-rose-500/10 text-rose-400 hover:bg-rose-500 hover:text-white transition-all duration-300 active:scale-[0.98]"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-5 h-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2.3"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                            />

                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M16 17l5-5-5-5"
                            />

                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M21 12H9"
                            />
                        </svg>

                        <span
                            class="text-xs font-black tracking-[0.2em] uppercase"
                        >
                            Logout
                        </span>
                    </button>
                </div>
            </div>
        </aside>

        <!-- CONTENT -->
        <div class="flex flex-col flex-1 min-w-0 overflow-hidden">
            <!-- MOBILE TOPBAR -->
            <header
                class="lg:hidden h-20 shrink-0 bg-white/80 backdrop-blur-xl border-b border-slate-200 flex items-center justify-between px-5"
            >
                <!-- MENU -->
                <button
                    @click="isMobileMenuOpen = true"
                    class="flex items-center justify-center w-12 h-12 transition-all shadow-sm rounded-2xl bg-slate-100 active:scale-95"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="w-6 h-6 text-indigo-600"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2.5"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M4 7h16"
                        />

                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M4 12h16"
                        />

                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M4 17h16"
                        />
                    </svg>
                </button>

                <!-- TITLE -->
                <div class="text-center min-w-0 px-4">
                    <h2
                        class="text-sm font-black tracking-[0.25em] uppercase text-slate-800 truncate"
                    >
                        {{ storeName }}
                    </h2>

                    <p
                        class="mt-1 text-[9px] uppercase tracking-[0.25em] text-slate-400 font-bold"
                    >
                        POS SYSTEM
                    </p>
                </div>

                <!-- PROFILE -->
                <div
                    class="flex items-center justify-center w-12 h-12 text-xs font-black text-white shadow-lg rounded-2xl bg-gradient-to-br from-slate-800 to-slate-950 shrink-0"
                >
                    {{ userInitial }}
                </div>
            </header>

            <!-- PAGE -->
            <main
                class="flex-1 min-h-0 overflow-y-auto overflow-x-hidden custom-scrollbar"
            >
                <slot />
            </main>
        </div>
    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
    width: 5px;
    height: 5px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(148, 163, 184, 0.3);
    border-radius: 999px;
}

.custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

* {
    -webkit-tap-highlight-color: transparent;
}

button {
    user-select: none;
}
</style>