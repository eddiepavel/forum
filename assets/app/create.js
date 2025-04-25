document.addEventListener('DOMContentLoaded', () => {
    const dropDownButton = document.getElementById('drop-down');
    const combobox = document.getElementById('combobox');
    const categoryList = document.getElementById('category-list');
    const listItems = document.querySelectorAll('#category-list li');
    const hiddenCategoryInput = document.getElementById('selected-category');
    const form = document.getElementById('post-form');
    const outputSelected = document.getElementById('output-selected');
    const selectedCategories = new Set();

    const updateSelectedCategories = () => {
        outputSelected.textContent = `Selected categories: ${Array.from(selectedCategories).join(', ') || 'None'}`;
    }

    // Function to toggle the dropdown
    const toggleDropdown = () => {
        categoryList.classList.toggle('hidden');
    };

    // Add event listeners for dropdown toggle
    dropDownButton.addEventListener('click', toggleDropdown);
    combobox.addEventListener('click', toggleDropdown);

    // Close the dropdown when clicking outside
    document.addEventListener('click', (event) => {
        const isClickInside = combobox.contains(event.target) || categoryList.contains(event.target) || dropDownButton.contains(event.target);
        if (!isClickInside) {
            categoryList.classList.add('hidden'); // Collapse the dropdown
        }
    });

    form.addEventListener('submit', (event) => {
        // Check the number of selected categories
        if (selectedCategories.size === 0) {
            selectedCategories.add('General');
            hiddenCategoryInput.value = Array.from(selectedCategories);
        }

        if (selectedCategories.size > 4) {
            alert("You can select up to 4 categories.");
            event.preventDefault(); // Prevent form submission
        }
    });

    // Add click event listeners to list items
    listItems.forEach((item) => {
        item.addEventListener('click', () => {
            const selectedCategory = item.querySelector('span.block').textContent.trim();

            // Toggle selection of the category
            if (selectedCategories.has(selectedCategory)) {
                selectedCategories.delete(selectedCategory);
                item.classList.remove('bg-indigo-100'); // Remove highlight
            } else {
                selectedCategories.add(selectedCategory);
                item.classList.add('bg-indigo-100'); // Highlight selected item
            }
    
            // Update the hidden input with selected categories
            hiddenCategoryInput.value = Array.from(selectedCategories).join(',');

            updateSelectedCategories();
    
            // Toggle the purple tick visibility
            const checkIcon = item.querySelector('span[id^="check-"]');
            if (checkIcon) {
                if (selectedCategories.has(selectedCategory)) {
                    checkIcon.classList.remove('text-white');
                    checkIcon.classList.add('text-indigo-600'); // Purple tick
                } else {
                    checkIcon.classList.add('text-white');
                    checkIcon.classList.remove('text-indigo-600');
                }
            }
        });
    });

    // Filter categories based on combobox input
    combobox.addEventListener('input', () => {
        const query = combobox.value.toLowerCase().trim();

        listItems.forEach((item) => {
            const category = item.querySelector('span.block').textContent.toLowerCase();
            if (category.includes(query)) {
                item.classList.remove('hidden');
            } else {
                item.classList.add('hidden');
            }
        });

        // Show the dropdown if it's hidden
        if (categoryList.classList.contains('hidden')) {
            categoryList.classList.remove('hidden');
        }
    });

    updateSelectedCategories();
});