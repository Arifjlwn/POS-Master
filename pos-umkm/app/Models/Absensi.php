<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Absensi extends Model
{
    // Nama tabelnya
    protected $table = 'absensi';

    // Kolom yang boleh diisi (Mass Assignment)
    protected $fillable = [
        'user_id',
        'store_id',
        'tanggal',
        'jam_masuk',
        'jam_pulang',
        'status'
    ];

    /**
     * Relasi: Absensi ini milik siapa (User)
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    /**
     * Relasi: Absensi ini di toko mana
     */
    public function store(): BelongsTo
    {
        return $this->belongsTo(Store::class);
    }
}
