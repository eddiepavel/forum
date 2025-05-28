document.addEventListener('DOMContentLoaded', () => {
    const hiddenPreselected = document.getElementById('preselected-categories');
    const dropDownButton = document.getElementById('drop-down');
    const combobox = document.getElementById('combobox');
    const categoryList = document.getElementById('category-list');
    const listItems = document.querySelectorAll('#category-list li');
    const hiddenCategoryInput = document.getElementById('selected-category');
    const form = document.getElementById('post-form');
    const outputSelected = document.getElementById('output-selected');
    const selectedCategories = new Set();

    // Preselect categories before rendering/output
    if (hiddenPreselected && hiddenPreselected.value) {
        hiddenPreselected.value
            .replace(/[\[\]]/g, '') // Remove [ and ] everywhere!
            .split(' ')
            .map(s => s.trim())
            .filter(s => s.length > 0)
            .forEach(cat => selectedCategories.add(cat));
    }

    if (window.location.pathname === '/edit') {
        const form = document.getElementById('post-form');
        if (form) {
            form.addEventListener('submit', function(e) {
                e.preventDefault(); // Prevent default POST
                const formData = new FormData(form);
                fetch(form.action, {
                    method: 'PUT',
                    body: formData,
                    credentials: 'same-origin'
                }).then(response => {
                    if (response.redirected) {
                        window.location.href = response.url;
                    } else {
                        return response.text();
                    }
                }).then(text => {
                    // Optional: Handle errors/messages here
                });
            });
        }
    }



const updateSelectedCategories = () => {
        outputSelected.textContent = `Selected categories: ${Array.from(selectedCategories).join(', ') || 'None'}`;
    };

    // Function to set highlight and tick state
    const setItemSelectedState = (item, isSelected) => {
        if (isSelected) {
            item.classList.add('bg-indigo-100');
        } else {
            item.classList.remove('bg-indigo-100');
        }
        const checkIcon = item.querySelector('span[id^="check-"]');
        if (checkIcon) {
            if (isSelected) {
                checkIcon.classList.remove('text-white');
                checkIcon.classList.add('text-indigo-600');
            } else {
                checkIcon.classList.add('text-white');
                checkIcon.classList.remove('text-indigo-600');
            }
        }
    };

    // Apply initial highlight/tick for preselected
    listItems.forEach((item) => {
        const name = item.querySelector('span.block').textContent.trim();
        setItemSelectedState(item, selectedCategories.has(name));
    });

    // Dropdown logic
    const toggleDropdown = () => {
        categoryList.classList.toggle('hidden');
    };

    dropDownButton.addEventListener('click', toggleDropdown);
    combobox.addEventListener('click', toggleDropdown);

    // Close dropdown when clicking outside
    document.addEventListener('click', (event) => {
        const isClickInside = combobox.contains(event.target) || categoryList.contains(event.target) || dropDownButton.contains(event.target);
        if (!isClickInside) {
            categoryList.classList.add('hidden');
        }
    });

    form.addEventListener('submit', (event) => {
        if (selectedCategories.size === 0) {
            selectedCategories.add('General');
            hiddenCategoryInput.value = Array.from(selectedCategories);
        }

        if (selectedCategories.size > 4) {
            alert("You can select up to 4 categories.");
            event.preventDefault();
        }
    });

    // Click listeners for category selection
    listItems.forEach((item) => {
        item.addEventListener('click', () => {
            const selectedCategory = item.querySelector('span.block').textContent.trim();

            // Toggle selection
            if (selectedCategories.has(selectedCategory)) {
                selectedCategories.delete(selectedCategory);
            } else {
                selectedCategories.add(selectedCategory);
            }
            hiddenCategoryInput.value = Array.from(selectedCategories).join(',');

            // Update UI
            setItemSelectedState(item, selectedCategories.has(selectedCategory));
            updateSelectedCategories();
        });
    });

    // Filter categories
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

        // Show dropdown if hidden
        if (categoryList.classList.contains('hidden')) {
            categoryList.classList.remove('hidden');
        }
    });

    updateSelectedCategories();
    // Set hidden input with preselected categories on initial load
    hiddenCategoryInput.value = Array.from(selectedCategories).join(',');

});