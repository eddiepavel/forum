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
