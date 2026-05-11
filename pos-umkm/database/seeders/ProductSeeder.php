<?php

namespace Database\Seeders;

use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class ProductSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('products')->insert([
            [
                'name' => 'Indomie Goreng Special',
                'price' => 3000,
                'stock' => 100,
                'image' => 'https://placehold.co/200x200?text=Indomie',
                'created_at' => now(),
                'updated_at' => now(),
            ],
            [
                'name' => 'Teh Botol Sosro',
                'price' => 5000,
                'stock' => 50,
                'image' => 'https://placehold.co/200x200?text=Teh+Botol+Sosro',
                'created_at' => now(),
                'updated_at' => now(),
            ],
            [
                'name' => 'Aqua 600ml',
                'price' => 4000,
                'stock' => 200,
                'image' => 'https://placehold.co/200x200?text=Aqua+600ml',
                'created_at' => now(),
                'updated_at' => now(),
            ],
        ]);
    }
}
