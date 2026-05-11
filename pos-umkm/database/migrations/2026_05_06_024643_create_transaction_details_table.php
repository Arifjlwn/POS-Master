<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('transaction_details', function (Blueprint $table) {
            $table->id();
            $table->foreignId('transaction_id')->constrained('transactions')->cascadeOnDelete(); //Relasi dengan tabel transactions, jika transaksi dihapus maka detailnya juga dihapus
            $table->foreignId('product_id')->constrained('products')->restrictOnDelete(); //Relasi dengan tabel products, jika produk dihapus maka tidak bisa dihapus jika masih ada transaksi yang menggunakan produk tersebut
            $table->integer('quantity'); //Jumlah produk yang dibeli dalam transaksi
            $table->integer('price'); //harga saat barang dibeli (antisipasi harga berubah)
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('transaction_details');
    }
};
