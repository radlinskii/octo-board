const registerServiceWorker = () => {
  if (navigator.serviceWorker) {
    navigator.serviceWorker.register('/service-worker.js')
    .then(() => {
      console.log("Service Worker registered.")
    })
    .catch(() => {
      console.err("Service Worker failed to register.")
    });
  }
};

window.onload = () => {
  registerServiceWorker();
}


