# SocialConnect Frontend Integration Guide

## Base URL

```
http://localhost:8080/api/v1
```

In production (via nginx reverse proxy), the app is served on port 80/443. All API calls are under `/api/v1`.

Static/uploaded files are served from:
```
http://localhost:8080/uploads/<filename>
```

---

## Authentication

### Token Storage
- JWT is returned on login/register as `token` (string).
- Token lifetime: **72 hours** (configurable via `JWT_EXPIRY_HOURS`).
- Store the token in `localStorage` or a secure cookie.

### Sending the Token
All protected endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <token>
```

### Optional Auth Endpoints
Some public endpoints accept an optional token. If provided, the response will include context-sensitive data (e.g., `is_liked`, `is_following`). If not provided, these fields default to `false`.

---

## Response Envelope

### Success
```json
{
  "success": true,
  "data": { ... }
}
```

### Error
```json
{
  "success": false,
  "message": "error description"
}
```

### Paginated
```json
{
  "success": true,
  "data": [ ... ],
  "page": 1,
  "limit": 20,
  "total": 142
}
```

---

## Pagination

All list endpoints accept query parameters:
- `page` — page number, minimum 1, default 1
- `limit` — items per page, 1–100, default 20

Example: `GET /api/v1/users/john/posts?page=2&limit=10`

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad request / validation error / business rule violation |
| 401 | Unauthenticated or invalid credentials |
| 404 | Resource not found |
| 409 | Conflict (e.g., username already taken) |
| 500 | Internal server error |

---

## Auth Endpoints

### POST `/auth/register`
Create a new account. No token required.

**Request (JSON):**
```json
{
  "username": "johndoe",       // required, 3–50 chars, alphanumeric
  "password": "secret123",     // required, 6–100 chars
  "display_name": "John Doe"   // optional, max 100 chars
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "token": "<jwt>",
    "user": { ... UserProfile ... }
  }
}
```

**Errors:**
- `400` — validation failure
- `409` — username already taken

---

### POST `/auth/login`
Authenticate an existing user.

**Request (JSON):**
```json
{
  "username": "johndoe",
  "password": "secret123"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "token": "<jwt>",
    "user": { ... UserProfile ... }
  }
}
```

**Errors:**
- `401` — invalid credentials

---

## User Endpoints

### UserProfile Shape
```json
{
  "id": "uuid",
  "username": "johndoe",
  "display_name": "John Doe",
  "bio": "Hello world",
  "avatar_url": "",
  "follower_count": 10,
  "following_count": 5,
  "is_following": false    // only meaningful when a logged-in user views someone else
}
```
> `is_following` is `false` for guests and `true/false` for authenticated viewers.

---

### GET `/users/:username`
Get a user's public profile. Optional auth (send token for `is_following` context).

**Response (200):** `data` = UserProfile

**Errors:**
- `404` — user not found

---

### PUT `/users/me` 🔒
Update the authenticated user's profile.

**Request (JSON):** All fields optional.
```json
{
  "display_name": "New Name",   // max 100 chars
  "bio": "New bio text"          // max 500 chars
}
```

**Response (200):** `data` = updated UserProfile

> **Note:** Avatar upload is not supported via this endpoint. There is no avatar upload endpoint currently.

---

## Post Endpoints

### PostResponse Shape
```json
{
  "id": "uuid",
  "content": "Hello #world",
  "image_url": "http://localhost:8080/uploads/abc.jpg",  // empty string if no image
  "author": { ... UserProfile ... },
  "tags": ["world"],
  "like_count": 5,
  "comment_count": 3,
  "is_liked": false,           // true if authenticated user has liked this post
  "created_at": "2024-01-15T10:30:00Z"
}
```
> Hashtags in `content` (e.g. `#golang`) are automatically parsed and stored. `tags` contains the parsed tag names without `#`.

---

### POST `/posts` 🔒
Create a post. Uses `multipart/form-data` (not JSON).

**Request (multipart/form-data):**
- `content` (text) — required if no image, max 5000 chars
- `image` (file) — optional image attachment, max 10 MB

**Response (201):** `data` = PostResponse

---

### GET `/posts/:id`
Get a single post. Optional auth.

**Response (200):** `data` = PostResponse

**Errors:**
- `404` — post not found

---

### DELETE `/posts/:id` 🔒
Delete a post. Must be the author.

**Response (200):**
```json
{ "success": true, "data": { "message": "post deleted" } }
```

**Errors:**
- `401` — not the author
- `404` — post not found

---

### GET `/users/:username/posts`
Get all posts by a user. Optional auth. Paginated.

**Response (200):** `data` = PostResponse[]

---

## Follow Endpoints

### POST `/users/:username/follow` 🔒
Follow a user.

**Response (200):** `{ "message": "followed" }`

**Errors:**
- `400` — cannot follow yourself / already following
- `404` — user not found

---

### DELETE `/users/:username/follow` 🔒
Unfollow a user.

**Response (200):** `{ "message": "unfollowed" }`

---

### GET `/users/:username/followers`
Get paginated list of followers.

**Response (200):** `data` = UserProfile[]

---

### GET `/users/:username/following`
Get paginated list of users the given user follows.

**Response (200):** `data` = UserProfile[]

---

## Timeline Endpoint

### GET `/timeline` 🔒
Get the authenticated user's feed: posts from people they follow. Paginated.

**Response (200):** `data` = PostResponse[]

> An empty timeline means the user doesn't follow anyone yet.

---

## Like Endpoints

### POST `/posts/:id/like` 🔒
Like a post.

**Response (200):** `{ "message": "liked" }`

**Errors:**
- `400` — already liked this post
- `404` — post not found

---

### DELETE `/posts/:id/like` 🔒
Unlike a post.

**Response (200):** `{ "message": "unliked" }`

---

### GET `/posts/:id/likes`
Get paginated list of users who liked a post.

**Response (200):** `data` = UserProfile[]

---

## Comment Endpoints

### CommentResponse Shape
```json
{
  "id": "uuid",
  "content": "Great post!",
  "author": { ... UserProfile ... },
  "created_at": "2024-01-15T10:35:00Z"
}
```

---

### POST `/posts/:id/comments` 🔒
Add a comment to a post.

**Request (JSON):**
```json
{
  "content": "Great post!"   // required, max 2000 chars
}
```

**Response (201):** `data` = CommentResponse

**Errors:**
- `404` — post not found

---

### GET `/posts/:id/comments`
Get paginated comments for a post.

**Response (200):** `data` = CommentResponse[]

---

### DELETE `/comments/:id` 🔒
Delete a comment. Must be the author.

**Response (200):** `{ "message": "comment deleted" }`

**Errors:**
- `401` — not the author
- `404` — comment not found

---

## Tag Endpoints

### GET `/tags/trending`
Get top 10 trending tags.

**Response (200):**
```json
{
  "success": true,
  "data": [
    { "name": "golang", "post_count": 42 }
  ]
}
```

---

### GET `/tags/:name/posts`
Get paginated posts with a specific tag. Optional auth.

**Response (200):** `data` = PostResponse[]

**Errors:**
- `404` — tag not found

---

## Message Endpoints (REST)

### MessageResponse Shape
```json
{
  "id": "uuid",
  "sender_id": "uuid",
  "receiver_id": "uuid",
  "content": "Hey!",
  "image_url": "",              // empty string if no image
  "read_at": null,              // ISO timestamp or null if unread
  "created_at": "2024-01-15T10:00:00Z"
}
```

### ConversationListResponse Shape
```json
{
  "user": { ... UserProfile ... },
  "last_message": { ... MessageResponse ... },
  "unread_count": 3
}
```

---

### GET `/messages/conversations` 🔒
Get all conversations for the current user. Paginated.

**Response (200):** `data` = ConversationListResponse[]

> Conversations are sorted by most recent message.

---

### GET `/messages/:username` 🔒
Get message history with a specific user. Paginated.

**Response (200):** `data` = MessageResponse[]

**Errors:**
- `404` — user not found

---

### POST `/messages/:username` 🔒
Send a message via REST (alternative to WebSocket). Uses `multipart/form-data`.

**Request (multipart/form-data):**
- `content` (text) — required if no image, max 5000 chars
- `image` (file) — optional image attachment, max 10 MB

**Response (201):** `data` = MessageResponse

---

## WebSocket — Real-Time Messaging

### Connection

Connect to:
```
ws://localhost:8080/api/v1/ws?token=<jwt>
```

The JWT is passed as a **query parameter**, not a header (WebSocket limitation).

If the token is missing or invalid, the server closes the connection with `401`.

---

### Sending a Message (Client → Server)

Send JSON text frames:
```json
{
  "type": "message",
  "receiver_id": "uuid-of-recipient",
  "content": "Hello!",
  "image_url": ""            // use REST POST /messages/:username for image messages
}
```

- `receiver_id` is the **UUID** of the recipient (not username).
- Either `content` or `image_url` must be non-empty.
- Image messages should use the REST endpoint; WebSocket only supports `image_url` referencing an already-uploaded file.

---

### Receiving a Message (Server → Client)

Both sender and receiver receive this frame when a message is delivered:
```json
{
  "type": "message",
  "id": "uuid",
  "sender_id": "uuid",
  "receiver_id": "uuid",
  "content": "Hello!",
  "image_url": "",
  "created_at": "2024-01-15T10:00:00Z"
}
```

> Messages are also persisted to the database. REST and WebSocket share the same message store.

---

### Reconnection

- If a user reconnects with the same JWT, the previous connection is closed and replaced.
- Implement automatic reconnection with exponential backoff on disconnect.

---

## File Uploads

- Accepted image formats: any common image (server does not enforce MIME type, only size).
- Max file size: **10 MB** (configurable via `MAX_UPLOAD_SIZE_MB`).
- After upload, the image is served at `http://localhost:8080/uploads/<filename>`.
- The full URL is returned in `image_url` fields in responses.
- Endpoints that accept image uploads use `multipart/form-data` (not `application/json`):
  - `POST /posts`
  - `POST /messages/:username`

---

## CORS

CORS is open in the current configuration (all origins allowed). In production, the nginx proxy handles origin restrictions.

---

## Key Behaviors & Edge Cases

| Scenario | Behavior |
|----------|----------|
| Viewing own profile | `is_following` is `false` (you don't follow yourself) |
| Like a post you already liked | `400` error |
| Follow yourself | `400` error |
| Follow someone you already follow | `400` error |
| Delete another user's post | `401` error |
| Delete another user's comment | `401` error |
| Post with no content and no image | `400` validation error |
| Message with no content and no image | `400` validation error |
| Posting with `#hashtag` in content | Tags are auto-extracted; returned in `tags[]` |
| `read_at` on MessageResponse | `null` = unread; ISO string = when it was read |
| WebSocket receiver offline | Message is persisted; delivered via REST on next fetch |

---

## Swagger / Interactive Docs

The full OpenAPI spec is available at:
```
http://localhost:8080/swagger/index.html
```

Use this to explore all endpoints, request/response shapes, and test calls interactively.

---

## Running the Backend Locally

```bash
# With Docker (recommended)
docker compose up -d --build
# API: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html

# Without Docker
cp .env.example .env   # edit DB credentials
make run
```
