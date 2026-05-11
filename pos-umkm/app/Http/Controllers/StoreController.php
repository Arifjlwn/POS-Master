<?php

namespace App\Http\Controllers;

use App\Models\Store;
use Illuminate\Http\Request;
use Inertia\Inertia;

class StoreController extends Controller
{
    public function create()
    {
        return Inertia::render('SetupToko');
    }

    public function store(Request $request)
    {
        // 1. Validasi Super Ketat (Termasuk File Gambar dan Array Fitur)
        $request->validate([
            'nama_toko' => 'required|string|max:255',
            'alamat_toko' => 'required|string',
            'telepon' => 'nullable|string|max:20',
            'qris_image' => 'nullable|image|mimes:jpeg,png,jpg|max:2048', // Maksimal 2MB
            'fitur_opsional' => 'required|array', // Harus berupa array dari Vue
        ]);

        // 2. Proses Upload QRIS (Kalau ada)
        $qrisPath = null;
        if ($request->hasFile('qris_image')) {
            $qrisPath = $request->file('qris_image')->store('qris_merchants', 'public');
        }

        // 3. Simpan Data Toko
        $store = Store::create([
            'nama_toko' => $request->nama_toko,
            'alamat_toko' => $request->alamat_toko,
            'telepon' => $request->telepon,
            'qris_image' => $qrisPath,
            'fitur_opsional' => $request->fitur_opsional, // Tersimpan otomatis jadi JSON
        ]);

        // 4. Ikatkan User sebagai 'Owner' di toko tersebut
        $request->user()->update([
            'store_id' => $store->id,
            'role' => 'owner'
        ]);

        // 5. Gas ke Kasir!
        return redirect()->route('pos');
    }

    // Form Pengaturan Toko
    public function edit(Request $request)
    {
        // Hanya Owner Yang bisa akses
        if ($request->user()->role !== 'owner') {
            abort(403, 'Akses Di Tolak.');
        }
        return Inertia::render('Store/Edit', [
            'store' => $request->user()->store
        ]);
    }

    // Proses Update Pengaturan Toko
    public function update(Request $request)
    {
        if ($request->user()->role !== 'owner') {
            abort(403, 'Akses Ditolak.');
        }

        $request->validate([
            'nama_toko' => 'required|string|max:255',
            'alamat_toko' => 'required|string',
            'telepon' => 'nullable|string|max:20',
            'fitur_opsional' => 'required|array',
        ]);

        $store = $request->user()->store;
        $store->update([
            'nama_toko' => $request->nama_toko,
            'alamat_toko' => $request->alamat_toko,
            'telepon' => $request->telepon,
            'fitur_opsional' => $request->fitur_opsional,
        ]);
        return redirect()->back()->with('message', 'Pengaturan Toko Berhasil Diperbarui!');
    }
}
