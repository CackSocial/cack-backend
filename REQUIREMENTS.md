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
- Notifications System
- Infinite Scrolling
- Enhanced Post Interactions (Reposts & Quotes)
- Mentions (@username)
- Profile Pictures Support

## Future Improvements
- Markdown support for posting

## Bugs
- Removing reposts is not working well
- Likes section in the profile is not working


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