## Architecture
Uses industry standard Clean Architecture principles to ensure separation of concerns and maintainability.
There are 3 main repositories for this project:
1. **Backend**: This repository contains the server-side code, including the API endpoints,
1. **Frontend - Web**: This repository contains the client-side code for the web application, including the user interface and interactions.
1. **Frontend - Mobile**: This repository contains the client-side code for the mobile application, including the user interface and interactions.

## Features
- Posting a message consists of image, text or both of them.
- Following other people
- Direct messaging consists of image, text or both of them. 
- Adding tags to posts.
- Trending tags section.
- Timeline of posts from following people.
- Interactions; likes, and nested comments.
- Deleting a post
- Deleting an account (hard delete)
- Bookmark
- Image preview and zoom
### Future Improvements

   1. Notifications System
      - Why: Currently, users have no way to know if they received a like, a new follower, or a mention without manually
        checking.
      - Implementation: Create a new NotificationsPage, a NotificationsStore, and integrate a real-time WebSocket connection to push notifications instantly. Add a badge to the sidebar icon.

   2. Infinite Scrolling
      - Why: The current implementation likely uses simple page-based fetching or loads everything at once.
      - Implementation: Implement an Intersection Observer-based infinite scroll for the timeline, explore feed, and profile
        feeds to improve UX and initial load times.

   3. Enhanced Post Interactions (Reposts & Quotes)
      - Why: Increases content virality and user interaction.
      - Implementation: Add "Repost" and "Quote Post" actions to the PostCard component. This will require updates to the
        PostComposer to handle quoting existing content.

   4. Mentions (@username)
      - Why: Essential for a social platform.
      - Implementation: Enhance the renderTaggedContent utility and the Textarea component to detect, link, and autocomplete @mentions in posts and comments.

   5. Profile Pictures Support
      - Why: Essential for a social platform.
      - Implementation: Add image support to profile pictures.

   6. Robust Global Error Handling
      - Why: API errors or WebSocket disconnections shouldn't fail silently or break the UI.
      - Implementation: Implement a global Toast/Notification system integrated into the API client and Zustand stores to
        provide user-friendly feedback on failures.

   7. Security Enhancements
      - Why: Storing auth tokens in localStorage can be vulnerable to XSS attacks.
      - Implementation: If the backend supports it, migrate to secure, HttpOnly cookie-based authentication and ensure CSRF
        protection is in place.

   8. Automated Testing
      - Why: The codebase currently lacks unit and integration tests, which makes future refactoring risky.
      - Implementation: Introduce Vitest and React Testing Library to test critical components (like PostComposer), stores
        (authStore), and utility functions.

## Tech Stack
### Backend
- Golang
- PostgreSQL
- Docker for containerization
- RESTful API for communication between frontend and backend
- JWT for authentication and authorization
- GORM for ORM (Object-Relational Mapping)
- Github Actions for CI/CD
- Digital Ocean for image storage
- Nginx for load balancing and reverse proxy
- Prometheus and Grafana for monitoring and logging
- Swagger for API documentation
- Go Modules for dependency management
- Go Test for unit testing

### Frontend - Web
- React.js
- Typescript

### Frontend - Mobile
- React Native