var CACHE_NAME = 'octo-board-cache-v1';
var urlsToCache = [
  '/',
  '/search',
  '/static/style/main.css',
  '/static/style/search.css',
  '/static/style/home.css',
  '/static/js/index.js',
  '/img/comments.svg',
  '/img/github-logo.svg',
  '/img/octocat.png',
  '/img/spidertocat.png',
  'img/favicon/favicon.ico',
  'img/favicon/ftkdict-120-115567.png',
  'img/favicon/ftkdict-152-115567.png',
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