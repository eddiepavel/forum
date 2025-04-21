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
    const selectedCategories = new Set();

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
});

document.addEventListener('DOMContentLoaded', () => {
    const deleteButton = document.querySelector('button[data-action="delete-post"]');
    if (!deleteButton) return;

    deleteButton.addEventListener('click', () => {
        const postID = deleteButton.dataset.postId;

        if (confirm("Are you sure you want to delete this post? This action cannot be undone.")) {
            fetch(`/view?id=${postID}`, {
                method: 'DELETE',
                headers: {
                    'csrf': document.getElementById('csrf').value,
                },
            })
                .then(response => {
                    if (response.ok) {
                        alert("Post deleted successfully.");
                        window.location.href = '/home';
                    } else {
                        alert("Failed to delete the post.");
                    }
                })
                .catch(error => {
                    console.error("Error:", error);
                    alert("An error occurred while deleting the post.");
                });
        }
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const pagePrev = document.getElementById('pagePrev');
    const pageNext = document.getElementById('pageNext');

    const updatePageQuery = (increment) => {
        const url = new URL(window.location.href);
        const currentPage = parseInt(url.searchParams.get('page')) || 1; // Default to page 1 if not present
        const newPage = currentPage + increment;

        if (newPage < 1) return; // Prevent negative or zero pages

        // Update the page query parameter
        url.searchParams.set('page', newPage);

        // Redirect to the updated URL
        window.location.href = url.toString();
    };

    if (pagePrev) {
        pagePrev.addEventListener('click', (event) => {
            event.preventDefault(); // Prevent default link behavior
            updatePageQuery(-1); // Decrement the page
        });
    }

    if (pageNext) {
        pageNext.addEventListener('click', (event) => {
            event.preventDefault(); // Prevent default link behavior
            updatePageQuery(1); // Increment the page
        });
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const sidebarToggle = document.getElementById('sidebar-toggle');
    const sidebar = document.getElementById('sidebar');

    // Define a media query for Tailwind's `lg` breakpoint (min-width: 1024px)
    const isDesktop = window.matchMedia('(min-width: 1024px)');

    // Function to open the sidebar
    const openSidebar = (event) => {
        sidebar.classList.remove('hidden', 'pointer-events-none', 'invisible'); // Make visible
        setTimeout(() => {
            sidebar.classList.remove('-translate-x-full'); // Slide in
        }, 20); // Slight delay to ensure transition applies
        event.stopPropagation(); // Prevent the click from propagating to the document
    };

    // Function to close the sidebar
    const closeSidebar = () => {
        sidebar.classList.add('-translate-x-full'); // Slide out
        setTimeout(() => {
            sidebar.classList.add('hidden', 'pointer-events-none', 'invisible'); // Hide after transition
        }, 200); // Match the transition duration
    };

    // Open the sidebar with a transition
    sidebarToggle.addEventListener('click', (event) => {
        if (!isDesktop.matches) {
            // Only open the sidebar for mobile viewports
            openSidebar(event);
        }
    });

    // Close the sidebar when clicking outside
    document.addEventListener('click', (event) => {
        const isClickInsideSidebar = sidebar.contains(event.target);
        const isClickOnToggle = sidebarToggle.contains(event.target);

        if (!isClickInsideSidebar && !isClickOnToggle && !isDesktop.matches) {
            // Only close the sidebar for mobile viewports
            closeSidebar();
        }
    });

    // Prevent clicks inside the sidebar from closing it
    sidebar.addEventListener('click', (event) => {
        event.stopPropagation(); // Prevent the click from propagating to the document
    });

    // Optional: Add a listener to handle viewport changes dynamically
    isDesktop.addEventListener('change', (e) => {
        if (e.matches) {
            // If switching to desktop, ensure the sidebar is always visible
            sidebar.classList.remove('hidden', '-translate-x-full', 'pointer-events-none', 'invisible');
        } else {
            // If switching to mobile, hide the sidebar initially
            sidebar.classList.add('hidden', '-translate-x-full', 'pointer-events-none', 'invisible');
        }
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const categoryToggle = document.getElementById('category-toggle');
    const categoryToggleRow = document.querySelector('h2'); // Select the entire row (h2 tag)
    const categoryContainer = document.getElementById('category-container');
    const categoryArrow = document.getElementById('category-arrow');

    // Ensure the list is collapsed initially
    categoryContainer.style.maxHeight = '0px';
    categoryContainer.style.overflow = 'hidden'; // Ensure content is hidden when collapsed
    categoryContainer.style.transition = 'max-height 0.3s ease-in-out'; // Add smooth transition

    // Toggle the category list visibility with a smooth transition
    const toggleDropdown = () => {
        if (categoryContainer.style.maxHeight === '0px' || !categoryContainer.style.maxHeight) {
            // Expand the list
            categoryContainer.style.maxHeight = categoryContainer.scrollHeight + 'px'; // Set to full height
            categoryArrow.classList.add('rotate-180'); // Rotate the arrow
        } else {
            // Collapse the list
            categoryContainer.style.maxHeight = '0px'; // Collapse to 0 height
            categoryArrow.classList.remove('rotate-180'); // Reset the arrow rotation
        }
    };

    // Add event listener to the button
    categoryToggle.addEventListener('click', (event) => {
        event.stopPropagation(); // Prevent the event from propagating to the parent <h2>
        toggleDropdown();
    });

    // Add event listener to the entire row
    categoryToggleRow.addEventListener('click', () => {
        toggleDropdown();
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const categoryList = document.getElementById('filter-list');
    const likedButton = Array.from(document.querySelectorAll('h2')).find(h2 => h2.textContent.trim() === "Liked");
    const createdButton = Array.from(document.querySelectorAll('h2')).find(h2 => h2.textContent.trim() === "Created");
    
    // Add click event listeners to category items
    categoryList.querySelectorAll('li').forEach((item) => {
        item.addEventListener('click', () => {
            const selectedCategory = item.querySelector('span.block').textContent.trim();
            const url = new URL(window.location.href);
            url.searchParams.set('category', selectedCategory);
            window.location.href = url.toString();
        });
    });

    // Add click event listener for "Liked"
    likedButton.addEventListener('click', () => {
        const url = new URL(window.location.href);
        url.searchParams.set('category', 'Liked');
        window.location.href = url.toString();
    });

    // Add click event listener for "Created"
    createdButton.addEventListener('click', () => {
        const url = new URL(window.location.href);
        url.searchParams.set('category', 'Created');
        window.location.href = url.toString();
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const resetFilterButton = document.getElementById('reset-filter');

    // Check if the URL has a "category" query parameter
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has('category')) {
        resetFilterButton.classList.remove('hidden'); // Show the button
    }

    // Add click event listener to reset the filter
    resetFilterButton.addEventListener('click', () => {
        const url = new URL(window.location.href);
        url.searchParams.delete('category'); // Remove the "category" query parameter
        resetFilterButton.classList.add('hidden'); // Hide the button
        url.searchParams.delete('page'); // Remove the "page" query parameter
        window.location.href = url.toString(); // Redirect to the updated URL
    });
});

document.addEventListener('DOMContentLoaded', () => {
    // Initialize button styles based on data-vote-state
    document.querySelectorAll('[data-vote-state]').forEach(button => {
        const state = button.dataset.voteState;
        if (state === 'upvote' && button.id.startsWith('upvote-')) {
            button.classList.add('bg-blue-300');
        } else if (state === 'downvote' && button.id.startsWith('downvote-')) {
            button.classList.add('bg-red-200');
        }
    });

    // Function to send a vote and update the UI dynamically
    const sendVote = (id, isPost, voteType) => {
        const formData = new FormData();
        const csrf = document.getElementById('csrf');
        formData.append(isPost ? 'post_id' : 'comment_id', id);
        formData.append('vote_type', voteType);
        formData.append('csrf', csrf.value);
    
        fetch('/vote', {
            method: 'POST',
            body: formData,
        }).then(response => {
            if (response.ok) {
                response.json().then(data => {
                    // Update the vote counts dynamically
                    const upvoteButton = document.getElementById(`upvote-${isPost ? 'post' : 'comment'}-${id}`);
                    const downvoteButton = document.getElementById(`downvote-${isPost ? 'post' : 'comment'}-${id}`);
    
                    // Update the vote counts inside the <span> elements
                    upvoteButton.querySelector('.vote-count').textContent = data.upvotes;
                    downvoteButton.querySelector('.vote-count').textContent = data.downvotes;
    
                    // Reset button styles
                    upvoteButton.classList.remove('bg-blue-300');
                    downvoteButton.classList.remove('bg-red-200');
    
                    // Apply the new style based on the vote state
                    if (data.user_vote === 'upvote') {
                        upvoteButton.classList.add('bg-blue-300');
                    } else if (data.user_vote === 'downvote') {
                        downvoteButton.classList.add('bg-red-200');
                    }
                });
            } else {
                console.error('Failed to register vote');
                alert('Failed to register your vote. Please try again.');
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('An error occurred while processing your vote.');
        });
    };

    // Add event listeners for upvote buttons
    document.querySelectorAll('[id^="upvote-"]').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.id.split('-')[2];
            const isPost = button.id.includes('post');
            sendVote(id, isPost, 'upvote');
        });
    });

    // Add event listeners for downvote buttons
    document.querySelectorAll('[id^="downvote-"]').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.id.split('-')[2];
            const isPost = button.id.includes('post');
            sendVote(id, isPost, 'downvote');
        });
    });
});

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