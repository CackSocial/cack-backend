// Package user implements the business logic for user registration,
// authentication, and profile management.
package user

import (
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
	"github.com/CackSocial/cack-backend/pkg/auth"
	"github.com/CackSocial/cack-backend/pkg/hash"
	"github.com/CackSocial/cack-backend/internal/domain"
)

// UserUseCase encapsulates all user-related business logic including
// registration, login, profile retrieval, and profile updates.
type UserUseCase struct {
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
	jwtSecret  string
	jwtExpiry  int
}

// NewUserUseCase creates a new UserUseCase with the given dependencies.
func NewUserUseCase(userRepo repository.UserRepository, followRepo repository.FollowRepository, jwtSecret string, jwtExpiry int) *UserUseCase {
	return &UserUseCase{
		userRepo:   userRepo,
		followRepo: followRepo,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
	}
}

// Register creates a new user account. It checks that the username is not
// already taken, hashes the password, persists the user, generates a JWT
// token, and returns a LoginResponse.
func (uc *UserUseCase) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Check if username is already taken.
	existing, _ := uc.userRepo.GetByUsername(req.Username)
	if existing != nil {
		return nil, ucerrors.ErrUsernameTaken
	}

	// Hash the password.
	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:    req.Username,
		Password:    hashed,
		DisplayName: req.DisplayName,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token.
	token, err := auth.GenerateToken(user.ID, uc.jwtSecret, uc.jwtExpiry)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserProfile{
			ID:             user.ID,
			Username:       user.Username,
			DisplayName:    user.DisplayName,
			Bio:            user.Bio,
			AvatarURL:      user.AvatarURL,
			FollowerCount:  0,
			FollowingCount: 0,
			IsFollowing:    false,
		},
	}, nil
}

// Login authenticates a user with username and password. It returns a JWT
// token and the user profile on success, or an invalid credentials error.
func (uc *UserUseCase) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userRepo.GetByUsername(req.Username)
	if err != nil || user == nil {
		return nil, ucerrors.ErrInvalidCredentials
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, ucerrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, uc.jwtSecret, uc.jwtExpiry)
	if err != nil {
		return nil, err
	}

	// Fetch follower/following counts.
	_, followerCount, _ := uc.followRepo.GetFollowers(user.ID, 1, 1)
	_, followingCount, _ := uc.followRepo.GetFollowing(user.ID, 1, 1)

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserProfile{
			ID:             user.ID,
			Username:       user.Username,
			DisplayName:    user.DisplayName,
			Bio:            user.Bio,
			AvatarURL:      user.AvatarURL,
			FollowerCount:  followerCount,
			FollowingCount: followingCount,
			IsFollowing:    false,
		},
	}, nil
}

// GetProfile retrieves a user's public profile by username. It includes
// follower/following counts and whether the current user (if authenticated)
// is following the profile owner.
func (uc *UserUseCase) GetProfile(username string, currentUserID string) (*dto.UserProfile, error) {
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, ucerrors.ErrUserNotFound
	}

	_, followerCount, _ := uc.followRepo.GetFollowers(user.ID, 1, 1)
	_, followingCount, _ := uc.followRepo.GetFollowing(user.ID, 1, 1)

	var isFollowing bool
	if currentUserID != "" {
		isFollowing, _ = uc.followRepo.IsFollowing(currentUserID, user.ID)
	}

	return &dto.UserProfile{
		ID:             user.ID,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Bio:            user.Bio,
		AvatarURL:      user.AvatarURL,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		IsFollowing:    isFollowing,
	}, nil
}

// UpdateProfile updates the authenticated user's profile fields (display name,
// bio). Only non-nil fields in the request are applied. Returns the updated
// user profile.
func (uc *UserUseCase) UpdateProfile(userID string, req *dto.UpdateProfileRequest) (*dto.UserProfile, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, ucerrors.ErrUserNotFound
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	_, followerCount, _ := uc.followRepo.GetFollowers(user.ID, 1, 1)
	_, followingCount, _ := uc.followRepo.GetFollowing(user.ID, 1, 1)

	return &dto.UserProfile{
		ID:             user.ID,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Bio:            user.Bio,
		AvatarURL:      user.AvatarURL,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		IsFollowing:    false,
	}, nil
}
