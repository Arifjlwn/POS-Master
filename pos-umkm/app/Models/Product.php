<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Factories\HasFactory;

class Product extends Model
{
    use HasFactory;

    protected $guarded = [];
    protected $fillable = ['name', 'sku', 'category', 'cost_price', 'price', 'stock', 'image'];
}
