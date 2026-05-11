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
        Schema::create('transactions', function (Blueprint $table) {
            $table->id();
            $table->foreignId('member_id')->nullable()->constrained('members')->nullOnDelete(); //Boleh kosong jika transaksi tanpa member
            $table->integer('total_price');
            $table->integer('pay_amount'); //Jumlah uang yang dibayarkan pelanggan
            $table->integer('return_amount'); //Jumlah uang kembali ke pelanggan jika ada kembalian
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('transactions');
    }
};
