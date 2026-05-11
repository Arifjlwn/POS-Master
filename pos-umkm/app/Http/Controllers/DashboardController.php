<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Inertia\Inertia;

class DashboardController extends Controller
{
    public function index(Request $request)
    {
        // Pengecekan Role (Hanya Owner)
        if ($request->user()->role !== 'owner') {
            return redirect()->route('absensi.index');
        }

        // --- Nanti data ini ditarik dari Database beneran ---
        // Contoh Data 7 Hari Terakhir
        $labelHari = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu'];
        $dataPenjualan = [1500000, 2100000, 1800000, 2400000, 2900000, 3200000, 3800000];

        // Contoh Top 5 Produk Paling Laku
        $topProduk = ['Indomie Goreng', 'Tolak Angin', 'Roti Aoka', 'Aqua 600ml', 'Beras 5kg'];
        $qtyProduk = [120, 85, 75, 50, 30];

        return Inertia::render('Dashboard', [
            'chartData' => [
                'labels' => $labelHari,
                'sales' => $dataPenjualan,
                'topProducts' => $topProduk,
                'topQty' => $qtyProduk
            ]
        ]);
    }
}
