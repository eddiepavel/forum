document.addEventListener('DOMContentLoaded', () => {
    const dropDownButton = document.getElementById('drop-down');
    const combobox = document.getElementById('combobox');
    const categoryList = document.getElementById('category-list');
    const listItems = document.querySelectorAll('#category-list li');
    const hiddenCategoryInput = document.getElementById('selected-category');
    const form = document.getElementById('post-form');
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

document.addEventListener('DOMContentLoaded', () => {
    let voteState = {}; // Tracks the current vote state for each post/comment by ID

    const updateVotes = (type, button, countElement, id) => {
        if (!voteState[id]) voteState[id] = null; // Initialize vote state for this ID

        if (type === 'upvote') {
            if (voteState[id] === 'upvote') {
                countElement.textContent = parseInt(countElement.textContent) - 1;
                voteState[id] = null;
                button.classList.remove('bg-blue-100');
                button.classList.add('bg-white');
            } else {
                if (voteState[id] === 'downvote') {
                    const downvoteButton = document.getElementById(`downvote-${id}`);
                    const downvoteCount = downvoteButton.querySelector('p');
                    downvoteCount.textContent = parseInt(downvoteCount.textContent) - 1;
                    downvoteButton.classList.remove('bg-red-100');
                    downvoteButton.classList.add('bg-white');
                }
                button.classList.remove('bg-white');
                countElement.textContent = parseInt(countElement.textContent) + 1;
                voteState[id] = 'upvote';
                button.classList.add('bg-blue-100');
            }
        } else if (type === 'downvote') {
            if (voteState[id] === 'downvote') {
                countElement.textContent = parseInt(countElement.textContent) - 1;
                voteState[id] = null;
                button.classList.remove('bg-red-100');
                button.classList.add('bg-white');
            } else {
                if (voteState[id] === 'upvote') {
                    const upvoteButton = document.getElementById(`upvote-${id}`);
                    const upvoteCount = upvoteButton.querySelector('p');
                    upvoteCount.textContent = parseInt(upvoteCount.textContent) - 1;
                    upvoteButton.classList.remove('bg-blue-100');
                    upvoteButton.classList.add('bg-white');
                }
                button.classList.remove('bg-white');
                countElement.textContent = parseInt(countElement.textContent) + 1;
                voteState[id] = 'downvote';
                button.classList.add('bg-red-100');
            }
        }
    };

    // Add event listeners to all upvote and downvote buttons
    document.querySelectorAll('[id^="upvote-"]').forEach((button) => {
        const id = button.id.split('-')[1]; // Extract the index from the button ID
        const countElement = button.querySelector('p');
        button.addEventListener('click', () => {
            updateVotes('upvote', button, countElement, id);
        });
    });

    document.querySelectorAll('[id^="downvote-"]').forEach((button) => {
        const id = button.id.split('-')[1]; // Extract the index from the button ID
        const countElement = button.querySelector('p');
        button.addEventListener('click', () => {
            updateVotes('downvote', button, countElement, id);
        });
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const loginLinks = document.querySelectorAll('a[href="/login"]');
    loginLinks.forEach(link => {
        const currentURL = window.location.pathname + window.location.search;
        if (currentURL.includes('view')) {
            link.href = `/login?redirect=${encodeURIComponent(currentURL)}`;
        } else {
            link.href = '/login';
        }
    });
});