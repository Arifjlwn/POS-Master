<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Factories\HasFactory;

class Transaction extends Model
{
    use HasFactory;

    protected $guarded = [];
    protected $fillable = ['total_price', 'pay_amount', 'payment_method', 'return_amount'];
}
