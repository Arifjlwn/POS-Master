<?php

use App\Http\Controllers\DashboardController;
use App\Http\Controllers\AbsensiController;
use App\Http\Controllers\ProfileController;
use App\Http\Controllers\ProductController;
use App\Http\Controllers\StoreController;
use App\Http\Controllers\KaryawanController;
use Illuminate\Foundation\Application;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use Illuminate\Support\Facades\DB;
use Inertia\Inertia;
use App\Models\Transaction;
use App\Models\TransactionDetail;
use App\Models\Product;

// Rute awal langsung arahkan ke login saja biar tidak bingung
Route::get('/', function () {
    return redirect('/login');
});

<<<<<<< HEAD
Route::get('/dashboard', function () {
    return Inertia::render('Dashboard');
})->name('dashboard');
=======
// --- SEMUA RUTE DI BAWAH INI WAJIB LOGIN ---
Route::middleware(['auth', 'verified', 'check.store'])->group(function () {
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970

    // Rute Dashboard
    Route::get('/dashboard', [DashboardController::class, 'index'])->name('dashboard');

    // Rute Setup Toko
    Route::get('/setup-toko', [StoreController::class, 'create'])->name('setup.toko');
    Route::post('/setup-toko', [StoreController::class, 'store'])->name('setup.toko.store');

    // Rute Pengaturan Toko (Khusus Owner)
    Route::get('/pengaturan', [StoreController::class, 'edit'])->name('store.edit');
    Route::post('/pengaturan', [StoreController::class, 'update'])->name('store.update');

    // Rute Karyawan
    Route::get('/karyawan', [KaryawanController::class, 'index'])->name('karyawan.index');
    Route::post('/karyawan', [KaryawanController::class, 'store'])->name('karyawan.store');
    Route::delete('/karyawan/{karyawan}', [KaryawanController::class, 'destroy'])->name('karyawan.destroy');

    // Rute Absensi
    Route::get('/absensi', [AbsensiController::class, 'index'])->name('absensi.index');
    Route::post('/absensi', [App\Http\Controllers\AbsensiController::class, 'store'])->name('absensi.store');

    // Tarik Data CSV Absensi Karyawan
    Route::get('/absensi/export', [App\Http\Controllers\AbsensiController::class, 'export'])->name('absensi.export');

    // Rute POS
    Route::get('/pos', function () {
        return Inertia::render('Pos/Index', [
            'products' => Product::all()
        ]);
    })->name('pos');

    // Rute Checkout Transaksi
    Route::post('/pos/checkout', function (Request $request) {
        DB::transaction(function () use ($request) {
            $transaction = Transaction::create([
                'total_price' => $request->total_price,
                'pay_amount' => $request->pay_amount,
                'payment_method' => $request->payment_method,
                'return_amount' => $request->pay_amount - $request->total_price,
            ]);

            foreach ($request->cart as $item) {
                TransactionDetail::create([
                    'transaction_id' => $transaction->id,
                    'product_id' => $item['id'],
                    'quantity' => $item['qty'],
                    'price' => $item['price'],
                ]);

                $product = Product::find($item['id']);
                $product->decrement('stock', $item['qty']);
            }
        });
        return redirect()->back();
    });

    // Rute Produk
    Route::get('/products/export', [ProductController::class, 'export'])->name('products.export');
    Route::post('/products/import', [ProductController::class, 'import'])->name('products.import');
    Route::resource('/products', ProductController::class);

}); // Penutup middleware auth

require __DIR__ . '/auth.php';
