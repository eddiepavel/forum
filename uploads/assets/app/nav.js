document.addEventListener('DOMContentLoaded', () => {

    const menu = document.getElementById('menu');
    const menuToggle = document.getElementById('menu-toggle');
    const menuUntoggle = document.getElementById('menu-untoggle');

    // Close the menu with a transition
    menuUntoggle.addEventListener('click', function () {
        menu.classList.add('translate-x-full'); // Slide out
        setTimeout(() => {
            menu.classList.add('hidden', 'pointer-events-none', 'invisible'); // Hide after transition
        }, 200); // Match the transition duration
    });

    // Open the menu with a transition
    menuToggle.addEventListener('click', function () {
        menu.classList.remove('hidden', 'pointer-events-none', 'invisible'); // Make visible
        setTimeout(() => {
            menu.classList.remove('translate-x-full'); // Slide in
        }, 20); // Slight delay to ensure transition applies
    });

});