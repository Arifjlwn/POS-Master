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
        Schema::table('products', function (Blueprint $table) {
            //Menambah 3 kolom baru setelah 'name'
            $table->string('sku')->nullable()->after('name');//untuk barcode
            $table->string('category')->nullable()->after('sku');//untuk kategori produk
            $table->integer('cost_price')->default(0)->after('category'); //untuk harga modal
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('products', function (Blueprint $table) {
            $table->dropColumn(['sku', 'category', 'cost_price']);
        });
    }
};
