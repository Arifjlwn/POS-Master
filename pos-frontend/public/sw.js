// public/sw.js
self.addEventListener('install', (event) => {
	self.skipWaiting();
});

self.addEventListener('activate', (event) => {
	event.waitUntil(clients.claim());
});

self.addEventListener('fetch', (event) => {
	// Biarkan browser mengambil data langsung dari internet/Vercel tanpa caching ribet dulu
	event.respondWith(fetch(event.request));
});
