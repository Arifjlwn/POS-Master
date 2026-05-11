<?php

namespace App\Http\Controllers;

use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Inertia\Inertia;

class KaryawanController extends Controller
{
    public function index(Request $request)
    {
        // Ambil user yang store_id sama dengan owner tapi bukan owner itu sendiri
        $karyawan = User::where('store_id', $request->user()->store_id)
        ->where('id', '!=', $request->user()->id)
        ->get();

        return Inertia::render('Karyawan/Index', [
            'karyawan' => $karyawan
        ]);
    }

    public function store(Request $request)
    {
        // 1. Validasi cuma butuh nama dan password sementara
        $request->validate([
            'name' => 'required|string|max:255',
            'password' => 'required|string|min:8',
        ]);

        // 2. LOGIKA AUTO-GENERATE NIK (Format: YYYYXXXX)
        $tahunSekarang = date('Y');

        // Cari karyawan terakhir yang NIK-nya berawalan tahun ini
        $karyawanTerakhir = User::where('nik', 'like', $tahunSekarang . '%')
                                ->orderBy('nik', 'desc')
                                ->first();

        if ($karyawanTerakhir) {
            // Ambil 4 digit terakhir, jadikan angka, tambah 1
            $urutan = (int) substr($karyawanTerakhir->nik, 4);
            $urutanBaru = $urutan + 1;
            // Gabungkan lagi dengan tahun (str_pad untuk nambahin nol di depan)
            $nikBaru = $tahunSekarang . str_pad($urutanBaru, 4, '0', STR_PAD_LEFT);
        } else {
            // Kalau belum ada sama sekali di tahun ini, mulai dari 0001
            $nikBaru = $tahunSekarang . '0001';
        }

        // 3. Simpan ke database
        User::create([
            'name' => $request->name,
            'nik' => $nikBaru, // NIK Otomatis masuk
            'email' => null,   // Email dikosongkan karena kasir pakai NIK
            'password' => \Illuminate\Support\Facades\Hash::make($request->password),
            'role' => 'kasir',
            'store_id' => $request->user()->store_id,
        ]);

        return redirect()->back()->with('message', 'Karyawan berhasil ditambahkan dengan NIK: ' . $nikBaru);
    }

    public function destroy(User $karyawan)
    {
        $karyawan->delete();
        return redirect()->back()->with('message', 'Karyawan berhasil dihapus!');
    }
}
