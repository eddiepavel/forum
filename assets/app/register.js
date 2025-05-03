document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('[data-toggle="password"]').forEach(button => {
        button.addEventListener('click', () => {
            const fieldId = button.getAttribute('data-target');
            const field = document.getElementById(fieldId);
            if (field.type === "password") {
                field.type = "text";
                button.textContent = "Hide";
            } else {
                field.type = "password";
                button.textContent = "Show";
            }
        });
    });
});