document.addEventListener('DOMContentLoaded', () => {
    const loginLinks = document.querySelectorAll('a[href="/login"]');
    loginLinks.forEach(link => {
        const currentURL = window.location.pathname + window.location.search;
        if (currentURL.includes('view') || currentURL.includes('home')) {
            link.href = `/login?redirect=${encodeURIComponent(currentURL)}`;
        } else {
            link.href = '/login';
        }
    });
});