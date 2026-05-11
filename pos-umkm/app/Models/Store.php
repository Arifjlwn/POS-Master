<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Store extends Model
{
    use HasFactory;

    protected $guarded = [];

    // cast kolom JSON agar otomatis menjadi array
    protected $casts = [
        'fitur_opsional' => 'array',
    ];

    // 1 toko punya banyak user
    public function users()
    {
        return $this->hasMany(User::class);
    }
}
