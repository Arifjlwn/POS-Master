export default [
  {
    path: "/laundry/laporan",
    name: "LaundryLaporan",
    component: () =>
      import("@/modules/jasalayanan/laundry/views/LaporanLaundry.vue"),
    meta: { requiresAuth: true, role: "owner" },
  },
  {
    path: "/laundry/pos",
    name: "LaundryPOS",
    component: () =>
      import("@/modules/jasalayanan/laundry/views/PosLaundry.vue"),
    meta: { requiresAuth: true },
  },
];
