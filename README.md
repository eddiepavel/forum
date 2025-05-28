# Forum App

Welcome to **Forum App**, a collaborative project built by a team of passionate developers. This project is a fully functional forum application where users can create posts, comment, vote, and interact with each other. It’s a culmination of our learning journey, showcasing various modern web development techniques and best practices.

---

## 🚀 Features

- **Dynamic Voting System**: Users can upvote or downvote posts and comments in real-time using AJAX.
- **Authentication**: Secure login, registration, and session management.
- **Post Creation**: Users can create posts with multiple categories and rich content.
- **Pagination & Filtering**: Posts are paginated and can be filtered by categories or user-specific views (e.g., "Liked" or "Created").
- **Responsive Design**: Built with Tailwind CSS for a seamless experience across devices.
- **Middleware-Driven Architecture**: Modular middleware for authentication, CSRF protection, rate limiting, and more.
- **Comprehensive Routing**: Routes are method-specific (GET, POST, DELETE) and chained with middleware.
- **Partial HTML Templates**: Reusable HTML components for consistent UI and maintainability.
- **Error Handling**: Graceful error pages and flash messages for user feedback.

---

## 🛠️ How It Works

### 1. **AJAX for Real-Time Voting**
   - Voting is handled dynamically using JavaScript and AJAX.
   - Example: When a user clicks an upvote or downvote button, an AJAX request is sent to the `/vote` endpoint. The server processes the vote and returns updated counts, which are then reflected on the UI without a page reload.
   - Relevant Files:
     - JavaScript: [`assets/app/view.js`](assets/app/view.js)
     - Backend: [`handlers/auth/view.go`](handlers/auth/view.go), [`database/post.go`](database/post.go)

### 2. **Middleware-Driven Architecture**
   - Middleware is used to handle common tasks like logging, authentication, CSRF protection, and rate limiting.
   - Example: The `ChainMiddleware` function in [`middleware/middleware.go`](middleware/middleware.go) allows us to stack multiple middleware functions for each route.
   - Key Middlewares:
     - `AuthMiddleware`: Ensures only authenticated users can access certain routes.
     - `CsrfTokenMiddleware`: Protects against CSRF attacks.
     - `RateLimitMiddleware`: Prevents abuse by limiting the number of requests.

### 3. **Tailwind CSS for Styling**
   - We used Tailwind CSS for rapid and consistent styling.
   - Example: Classes like `bg-indigo-600`, `rounded-md`, and `shadow-lg` are used throughout the app for a modern and responsive design.
   - Relevant Files:
     - [`assets/input.css`](assets/input.css)
     - [`tailwind.config.js`](tailwind.config.js)

### 4. **Partial HTML Templates**
   - HTML templates are split into reusable components (partials) for maintainability.
   - Example: The navigation bar is defined in [`assets/partials/nav.html`](assets/partials/nav.html) and included in multiple pages.
   - Key Partials:
     - `nav.html`: Navigation bar.
     - `home.html`: Home page layout.
     - `posts.html`: Post listing.
     - `view.html`: Single post view.

### 5. **Comprehensive Routing**
   - Routes are defined with method-specific handlers and middleware chaining.
   - Example: The `/create` route uses `GET` for rendering the form and `POST` for handling form submissions.
   - Relevant File: [`routes/routes.go`](routes/routes.go)

### 6. **Database Integration**
   - SQLite is used as the database, with migrations for schema management.
   - Example: The `votes` table tracks upvotes and downvotes for posts and comments.
   - Relevant Files:
     - Migrations: [`migrations/post.sql`](migrations/post.sql), [`migrations/votes.sql`](migrations/votes.sql)
     - Queries: [`database/post_helpers.go`](database/post_helpers.go)

---

## 🌟 Meet the Team

This project was brought to life by a group of enthusiastic developers who collaborated, learned, and grew together. Here’s the dream team:

- **Giannis Georgakopoulos**: Backend wizard.
- **Edouardos Pavel**: Frontend enthusiast and Tailwind CSS expert.
- **Stamatis Manousis**: Database architect.
- **Giorgos Pavrianidis**: Deployment officer and upcoming QA engineer.

---