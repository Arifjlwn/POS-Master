<?php

namespace App\Http\Controllers;

use App\Models\Product;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\Response;
use Illuminate\Database\QueryException;
use Inertia\Inertia;

class ProductController extends Controller
{
    // Menampilkan Halaman Master Produk
    public function index(Request $request)
    {
        // 1. Kerangka query ke tabel produk
        $query = Product::query();

        // 2. Jika ada pencarian dari kasir
        if ($request->search) {
            $query->where('name', 'like', '%' . $request->search . '%')
                ->orWhere('sku', 'like', '%' . $request->search . '%');
        }

        // 3. Jika ada filter kategori
        if ($request->category) {
            $query->where('category', $request->category);
        }

        // 4. Ambil data, urutkan dari yang terbaru, dan potong 10 per halaman
        // withQueryString() berguna agar saat pindah halaman 2, hasil pencariannya tidak hilang
        $products = $query->latest()->paginate(10)->withQueryString();

        // 5. Ambil daftar kategori yang unik langsung dari database untuk dropdown
        $categories = Product::whereNotNull('category')->distinct()->pluck('category');

        return inertia('Products/Index', [
            'products' => $products,
            'filters' =>[
                'search' => $request->search ?? '',
                'category' => $request->category ?? ''
                ],
            'categories' => $categories
        ]);
    }

    // Menyimpan Produk Baru Beserta Gambar
    public function store(Request $request)
    {
        // Validasi Inputan
        $request->validate([
            'name' => 'required|string|max:255',
            'sku' => 'nullable|string|max:100',
            'category' => 'nullable|string|max:100',
            'cost_price' => 'required|numeric',
            'price' => 'required|numeric',
            'stock' => 'required|numeric',
            'image' => 'nullable|image|mimes:jpeg,png,jpg|max:2048', // Maks 2MB
        ]);

        $imagePath = null;
        // Cek apakah ada file di upload
        if ($request->hasFile('image')) {
            // Simpan ke folder 'storage/app/public/products'
            $path = $request->file('image')->store('products', 'public');
            // Jadikan URL public untuk disimpan di DB
            $imagePath = '/storage/' . $path;
        }

        // Masukan ke database
        Product::create([
            'name' => $request->name,
            'sku' => $request->sku,
            'category' => $request->category,
            'cost_price' => $request->cost_price,
            'price' => $request->price,
            'stock' => $request->stock,
            'image' => $imagePath,
        ]);

        return redirect()->back();
    }

    // Update Barang Beserta Gambar
    public function update(Request $request, Product $product)
    {
        $request->validate([
            'name' => 'required|string|max:255',
            'sku' => 'nullable|string|max:100',
            'category' => 'nullable|string|max:100',
            'cost_price' => 'required|numeric',
            'price' => 'required|numeric',
            'stock' => 'required|numeric',
            'image' => 'nullable|image|mimes:jpeg,png,jpg|max:2048', // Maks 2MB
        ]);

        $imagePath = $product->image; // Pertahankan gambar lama jika tidak ada yang baru

        // Jika ada gambar baru, simpan dan hapus gambar lama
        if ($request->hasFile('image')) {
            // Hapus gambar lama
            if ($product->image && str_starts_with($product->image, '/storage/')) {
                Storage::disk('public') -> delete(str_replace('/storage/','',$product->image));
            }

            // Simpan gambar baru
            $path = $request->file('image')->store('products', 'public');
            $imagePath = '/storage/' . $path;
        }

        // Update data di database
        $product->update([
            'name' => $request->name,
            'sku' => $request->sku,
            'category' => $request->category,
            'cost_price' => $request->cost_price,
            'price'=> $request->price,
            'stock' => $request->stock,
            'image' => $imagePath,
        ]);

        return redirect()->back();
    }

    // Fungsi Hapus Barang
    public function destroy(Product $product)
    {
        try {
            // Coba hapus data dari database DULU
            $product->delete();

            // JIKA BERHASIL: baru hapus foto fisiknya dari folder
            if ($product->image && str_starts_with($product->image, '/storage/')) {
                Storage::disk('public')->delete(str_replace('/storage/', '', $product->image));
            }

            return redirect()->back();

        } catch (QueryException $e) {
            // JIKA GAGAL KARENA FOREIGN KEY: Tangkap errornya dan kirim pesan ke Vue
            return redirect()->back()->withErrors([
                'message' => 'GAGAL: Produk ini tidak bisa dihapus karena sudah tercatat di dalam riwayat transaksi (struk)!'
            ]);
        }
    }

    // Fungsi Export Produk ke CSV
    public function export()
    {
        $fileName = 'master_produk_' . date('Y-m-d') . '.csv';
        $products = Product::all();

        $headers = [
            "Content-Type" => "text/csv",
            "Content-Disposition" => "attachment; filename=$fileName",
            "Pragma" => "no-cache",
            "Cache-Control" => "must-revalidate, post-check=0, pre-check=0",
            "Expires" => "0"
        ];

        // Baris pertama (Header Kolom)
        $columns = ['Nama Produk', 'SKU/Barcode', 'Kategori', 'HPP', 'Harga Jual', 'Stok'];

        $callback = function() use($products, $columns) {
            $file = fopen('php://output', 'w');
            fputcsv($file, $columns); // Tulis header ke file

            foreach ($products as $product) {
                // Tulis data per baris
                fputcsv($file, [
                    $product->name,
                    $product->sku,
                    $product->category,
                    $product->cost_price,
                    $product->price,
                    $product->stock
                ]);
            }
            fclose($file);
        };
        return Response::stream($callback, 200, $headers);
    }

    // --- FUNGSI IMPORT CSV ---
    public function import(Request $request)
    {
        // Validasi file
        $request->validate([
            'file' => 'required|mimes:csv,txt'
        ]);

        $file = $request->file('file');

        // Baca isi file
        $fileContents = file($file->getPathname());

        foreach ($fileContents as $key => $line) {
            // Lewati baris pertama (karena itu Header kolom)
            if ($key == 0) continue;

            $data = str_getcsv($line);

            // Pastikan baris tidak kosong
            if (count($data) >= 6) {
                // Logika Pintar: Cari berdasarkan SKU.
                // Jika ketemu -> Update. Jika tidak -> Create baru.
                Product::updateOrCreate(
                    ['sku' => $data[1]], // Acuan pencarian (Barcode)
                    [
                        'name' => $data[0],
                        'category' => $data[2],
                        'cost_price' => $data[3],
                        'price' => $data[4],
                        'stock' => $data[5],
                    ]
                );
            }
        }
        return redirect()->back();
    }
}
