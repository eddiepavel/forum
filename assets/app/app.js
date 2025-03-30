document.addEventListener('DOMContentLoaded', () => {
    const dropDownButton = document.getElementById('drop-down');
    const combobox = document.getElementById('combobox');
    const categoryList = document.getElementById('category-list');
    const listItems = document.querySelectorAll('#category-list li');
    const hiddenCategoryInput = document.getElementById('selected-category');
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
            const selectedCategory = item.querySelector('span.block').textContent.trim();

            // Update the combobox input with the selected category
            combobox.value = selectedCategory;

            // Update the hidden input with the selected category
            hiddenCategoryInput.value = selectedCategory;

            // Optionally hide the dropdown
            categoryList.classList.add('hidden');

            // Highlight the selected item
            listItems.forEach((li) => {
                const checkIcon = li.querySelector('span[id^="check-"]');
                if (checkIcon) {
                    checkIcon.classList.add('text-white');
                    checkIcon.classList.remove('text-indigo-600');
                }
            });

            const checkIcon = item.querySelector('span[id^="check-"]');
            if (checkIcon) {
                checkIcon.classList.remove('text-white');
                checkIcon.classList.add('text-indigo-600');
            }
        });
    });

    // Form submission validation
    form.addEventListener('submit', (event) => {
        if (!hiddenCategoryInput.value) {
            console.log(hiddenCategoryInput.value)
            event.preventDefault(); // Prevent form submission
            alert('Please select a category before submitting.');
        }
    });
});