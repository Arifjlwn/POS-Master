<?php

namespace App\Http\Controllers;

use App\Models\User;
use App\Models\Absensi;
use Carbon\Carbon;
use Illuminate\Http\Request;
use Inertia\Inertia;

class AbsensiController extends Controller
{
    public function index(Request $request)
    {
        $karyawan = User::where('store_id', $request->user()->store_id)->get();

        // 1. Tangkap filter tanggal (Kalau tidak ada filter, pakai hari ini)
        $tanggalFilter = $request->query('tanggal', Carbon::today()->toDateString());

        // 2. Tarik riwayat berdasarkan tanggal yang dipilih
        $riwayatAbsensi = Absensi::with('user:id,name,nik,role')
            ->where('store_id', $request->user()->store_id)
            ->where('tanggal', $tanggalFilter)
            ->orderBy('created_at', 'desc')
            ->get();

        return Inertia::render('Absensi/Index', [
            'karyawan' => $karyawan,
            'riwayat' => $riwayatAbsensi,
            'filterTanggal' => $tanggalFilter // Lempar tanggal ke Vue
        ]);
    }

    public function store(Request $request)
    {
        $userId = $request->user_id;
        $storeId = $request->user()->store_id;
        $today = Carbon::today()->toDateString();
        $now = Carbon::now()->toTimeString();

        $absensi = Absensi::where('user_id', $userId)->where('tanggal', $today)->first();

        if ($request->jenis === 'Masuk') {
            if ($absensi) return redirect()->back()->with('error', 'Sudah absen masuk.');
            Absensi::create(['user_id' => $userId, 'store_id' => $storeId, 'tanggal' => $today, 'jam_masuk' => $now]);
            return redirect()->back()->with('message', 'Berhasil Absen Masuk!');
        } else {
            if (!$absensi) return redirect()->back()->with('error', 'Belum absen masuk.');
            if ($absensi->jam_pulang) return redirect()->back()->with('error', 'Sudah absen pulang.');
            $absensi->update(['jam_pulang' => $now]);
            return redirect()->back()->with('message', 'Berhasil Absen Pulang!');
        }
    }

    // FITUR KHUSUS OWNER: DOWNLOAD CSV
    public function export(Request $request)
    {
        if ($request->user()->role !== 'owner') {
            abort(403, 'Akses Ditolak.');
        }

        $bulan = $request->query('bulan', date('m'));
        $tahun = $request->query('tahun', date('Y'));

        $data = Absensi::with('user')->where('store_id', $request->user()->store_id)
                ->whereMonth('tanggal', $bulan)->whereYear('tanggal', $tahun)
                ->orderBy('tanggal', 'asc')->get();

        $csv = "Tanggal,NIK,Nama Karyawan,Jam Masuk,Jam Pulang,Status\n";
        foreach($data as $row) {
            $nik = $row->user->nik ?? '-';
            $nama = $row->user->name;
            $csv .= "{$row->tanggal},{$nik},{$nama},{$row->jam_masuk},{$row->jam_pulang},{$row->status}\n";
        }

        return response($csv)
            ->header('Content-Type', 'text/csv')
            ->header('Content-Disposition', "attachment; filename=Laporan_Absen_{$tahun}_{$bulan}.csv");
    }
}
