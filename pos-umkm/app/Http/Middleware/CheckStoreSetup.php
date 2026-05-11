<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class CheckStoreSetup
{
    /**
     * Handle an incoming request.
     *
     * @param  Closure(Request): (Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $user = $request->user();

        // Jika user login, dia owner, tapi belom punya toko
        if ($user && $user->role === 'owner' && !$user->store_id) {
            // Biar ga looping, masuk ke setup-toko
            if (!$request->is('setup-toko*') && !$request->is('logout')) {
                return redirect()->route('setup.toko');
            }
        }
        return $next($request);
    }
}
