document.addEventListener('DOMContentLoaded', function() {
    const tabButtons = document.querySelectorAll('.tab-button');
    const contentSections = document.querySelectorAll('.content-section');

    function updateTabStyles() {
        const isSmallScreen = window.innerWidth < 768;

        tabButtons.forEach(btn => {
            if (!isSmallScreen) {
                // Remove highlighting on large screens
                btn.classList.remove('text-blue-600', 'border-blue-600');
                btn.classList.add('text-gray-600');
            }
        });

        if (isSmallScreen && !Array.from(tabButtons).some(btn => btn.classList.contains('text-blue-600'))) {
            // Ensure first tab is active on mobile if none are
            tabButtons[0].classList.add('text-blue-600', 'border-blue-600');
            tabButtons[0].classList.remove('text-gray-600');
        }
    }

    function handleScreenSize() {
        const isSmallScreen = window.innerWidth < 768;

        if (isSmallScreen) {
            let activeIndex = Array.from(tabButtons).findIndex(btn =>
                btn.classList.contains('text-blue-600'));
            if (activeIndex === -1) activeIndex = 0;

            contentSections.forEach((section, index) => {
                if (index === activeIndex) {
                    section.classList.remove('hidden');
                } else {
                    section.classList.add('hidden');
                }
            });
        } else {
            contentSections.forEach(section => {
                section.classList.remove('hidden');
            });
        }

        updateTabStyles();
    }

    // Set initial state
    handleScreenSize();

    tabButtons.forEach((button, index) => {
        button.addEventListener('click', () => {
            if (window.innerWidth >= 768) return;

            tabButtons.forEach(btn => {
                btn.classList.remove('text-blue-600', 'border-blue-600');
                btn.classList.add('text-gray-600');
            });

            button.classList.remove('text-gray-600');
            button.classList.add('text-blue-600', 'border-blue-600');

            contentSections.forEach((section, idx) => {
                if (idx === index) {
                    section.classList.remove('hidden');
                } else {
                    section.classList.add('hidden');
                }
            });
        });
    });

    window.addEventListener('resize', handleScreenSize);
});