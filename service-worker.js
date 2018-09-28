var CACHE_NAME = 'octo-board-cache-v1';
var urlsToCache = [
  '/',
  '/search',
  '/static/style/main.css',
  '/static/style/search.css',
  '/static/style/home.css',
  '/static/js/index.js',
  '/static/img/comments.svg',
  '/static/img/github-logo.svg',
  '/static/img/octocat.png',
  '/static/img/spidertocat.png',
  '/static/img/favicon/favicon.ico',
  '/static/img/favicon/ftkdict-120-115567.png',
  '/static/img/favicon/ftkdict-152-115567.png',
  '/static/img/favicon/ftkdict-192-115567.png',
  '/static/img/favicon/ftkdict-512-115567.png',
];

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => {
        console.log('Opened cache');
        return cache.addAll(urlsToCache);
      })
  );
});

self.addEventListener('fetch', event => {
  event.respondWith(
    caches.match(event.request)
      .then(response => {
        if (response) {
          return response;
        }
        return fetch(event.request);
      }
    )
  );
});