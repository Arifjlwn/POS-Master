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
        Schema::table('users', function (Blueprint $table) {
            // Tambahkan kolom NIK, jadikan unik, tapi boleh kosong (karena Owner daftar pakai email)
            $table->string('nik', 8)->unique()->nullable()->after('id');

            // Ubah kolom email jadi boleh kosong (nullable)
            $table->string('email')->nullable()->change();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
{
    Schema::table('users', function (Blueprint $table) {
        $table->dropColumn('nik');
        $table->string('email')->nullable(false)->change();
    });
}
};
