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

document.addEventListener('DOMContentLoaded', function () {
    const notifBtn = document.getElementById('notification-btn');
    const notifDropdown = document.getElementById('notification-dropdown');

    function setEnvelopeIcon(hasUnread) {
        const plain = document.getElementById('notif-envelope');
        const alert = document.getElementById('notif-envelope-exclamation');
        if (hasUnread) {
            plain.classList.add('hidden');
            alert.classList.remove('hidden');
        } else {
            plain.classList.remove('hidden');
            alert.classList.add('hidden');
        }
    }


    async function fetchNotifications() {
        try {
            const response = await fetch('/unread', {
                credentials: 'same-origin'
            });
            if (!response.ok) {
                throw new Error('Failed to fetch');
            }
            const data = await response.json();
            const listNotif = Object.values(data).flat();
            renderNotifications(listNotif);
            setEnvelopeIcon(listNotif.length > 0);
        } catch (error) {
            renderNotifications([]);
            setEnvelopeIcon(false);
        }
    }

    function renderNotifications(listNotif) {
        if (!notifDropdown) return;
        console.log(listNotif);
        if (!Array.isArray(listNotif) || listNotif.length === 0) {
            console.log("gay");
            notifDropdown.innerHTML = `
                <div class="px-4 py-2 text-gray-400 text-sm">
                    No unread notifications
                </div>
                <div class="flex w-full border-t pl-1 py-1 justify-center items-center">
                    <a href="/notifications" class="text-sm/6 font-semibold text-gray-900">Check all notifications</a>
                </div>
            `;
            return;
        }
        notifDropdown.innerHTML =
            listNotif.map(n => {
                    // Add other notification types as needed
                    if (n.type === "comment") {
                        // Make whole notification clickable
                        return `<a href="/view?id=${n.postID}" class="block px-4 py-2 border-b border-gray-100 hover:bg-gray-50 transition-colors duration-150">
                        <div class="text-sm text-gray-800">${n.actor?.username || "Someone"} commented ${n.content || ""}</div>
                        <div class="text-xs text-gray-400">${n.time || ""}</div>
                    </a>`;
                    } else {
                        // Add other notification types as needed
                        return `<a href="/view?id=${n.postID}" class="block px-4 py-2 border-b border-gray-100 hover:bg-gray-50 transition-colors duration-150">
                        <div class="px-4 py-2 border-b border-gray-100 hover:bg-gray-50">
                        <div class="text-sm text-gray-800">${n.actor?.username || "Someone"} ${n.type}d your post</div>
                        <div class="text-xs text-gray-400">${n.time || ""}</div>
                        </div>
                    </a>`;
                    }


            }).join('') +
            `<div class="flex w-full border-t pl-1 py-1 justify-center items-center">
            <a href="/notifications" class="text-sm/6 font-semibold text-gray-900">Check all notifications</a>
        </div>`;
    }

    // Fetch once on page load
    fetchNotifications();
    // Fetch every 30 seconds
    setInterval(fetchNotifications, 30000);

    if (notifBtn && notifDropdown) {
        notifBtn.addEventListener('click', async function (e) {
            e.stopPropagation();
            notifDropdown.classList.toggle('hidden');

            await fetchNotifications();

            if (!notifDropdown.classList.contains('hidden')) {
                await fetch('/unread', {
                    method: 'PUT',
                    credentials: 'same-origin'
                });
                setEnvelopeIcon(false);
            }
        });

        // Hide dropdown on click outside
        document.addEventListener('click', function (e) {
            if (!notifDropdown.classList.contains('hidden')) {
                notifDropdown.classList.add('hidden');
            }
        });

        notifDropdown.addEventListener('click', function (e) {
            e.stopPropagation(); // Let links work, but doesn't close dropdown
        });
    }
});
