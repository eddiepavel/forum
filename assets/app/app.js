
document.addEventListener('DOMContentLoaded', () => {
    const dropDownButton = document.getElementById('drop-down');
    const combobox = document.getElementById('combobox');
    const categoryList = document.getElementById('category-list');
    const listItems = document.querySelectorAll('#category-list li');
    const form = document.getElementById('post-form');

        // Close the menu when clicking the "menu-untoggle" button
    document.getElementById('menu-untoggle').addEventListener('click', function () {
        const menu = document.getElementById('menu');
        menu.classList.add('hidden', 'pointer-events-none', 'invisible'); // Add classes to hide the menu
    });

    // Open the menu when clicking the "menu-toggle" button
    document.getElementById('menu-toggle').addEventListener('click', function () {
        const menu = document.getElementById('menu');
        menu.classList.remove('hidden', 'pointer-events-none', 'invisible'); // Remove classes to show the menu
    });

    // Function to toggle the dropdown
    const toggleDropdown = () => {
    categoryList.classList.toggle('hidden');
    };

    // Add event listeners for dropdown toggle
    dropDownButton.addEventListener('click', toggleDropdown);
    combobox.addEventListener('click', toggleDropdown);

    // Add click event listeners to list items
    listItems.forEach((item) => {
    item.addEventListener('click', () => {
        const checkIcon = item.querySelector('span#check-1');

        if (checkIcon) {
        // Check if the clicked item is already active
        if (checkIcon.classList.contains('text-indigo-600')) {
            // If active, deactivate it
            checkIcon.classList.remove('text-indigo-600');
            checkIcon.classList.add('text-white');
        } else {
            // Otherwise, deactivate all items and activate the clicked one
            listItems.forEach((li) => {
            const otherCheckIcon = li.querySelector('span#check-1');
            if (otherCheckIcon) {
                otherCheckIcon.classList.remove('text-indigo-600');
                otherCheckIcon.classList.add('text-white');
            }
            });

            // Activate the clicked item
            checkIcon.classList.remove('text-white');
            checkIcon.classList.add('text-indigo-600');
        }
        }
    });
    });

    // Form submission validation
    form.addEventListener('submit', (event) => {
    const isAnyItemActive = Array.from(listItems).some((item) => {
        const checkIcon = item.querySelector('span#check-1');
        return checkIcon && checkIcon.classList.contains('text-indigo-600');
    });

    if (!isAnyItemActive) {
        event.preventDefault(); // Prevent form submission
        alert('Please select at least one category before posting.');
    }
    });
});
