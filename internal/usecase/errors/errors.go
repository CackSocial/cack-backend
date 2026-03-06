// Package errors defines domain errors used across the usecase layer.
package errors

import "errors"

// Authentication and authorization errors.
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)

// User-related errors.
var (
	ErrUsernameTaken = errors.New("username already taken")
	ErrUserNotFound  = errors.New("user not found")
)

// Post-related errors.
var (
	ErrPostNotFound    = errors.New("post not found")
	ErrAlreadyReposted = errors.New("already reposted this post")
	ErrRepostNotFound  = errors.New("repost not found")
	ErrCannotRepost    = errors.New("cannot repost a repost")
	ErrContentRequired = errors.New("content is required")
)

// Follow-related errors.
var (
	ErrSelfFollow     = errors.New("cannot follow yourself")
	ErrAlreadyFollowing = errors.New("already following this user")
)

// Like-related errors.
var (
	ErrAlreadyLiked = errors.New("already liked this post")
)

// Comment-related errors.
var (
	ErrCommentNotFound = errors.New("comment not found")
)

// Bookmark-related errors.
var (
	ErrAlreadyBookmarked = errors.New("already bookmarked this post")
)

// Notification-related errors.
var (
	ErrNotificationNotFound = errors.New("notification not found")
)
