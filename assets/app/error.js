document.addEventListener('DOMContentLoaded', () => {
    const backButton = document.getElementById('back-btn');
    const homeButton = document.getElementById('home-btn');

    backButton.addEventListener('click', () => {
        console.log('Back button clicked');
        window.history.back();
    });

    homeButton.addEventListener('click', () => {
        console.log('Home button clicked');
        window.location.href = '/';
    });
});