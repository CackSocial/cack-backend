package dto

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password    string `json:"password" binding:"required,min=6,max=100"`
	DisplayName string `json:"display_name" binding:"max=100"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Bio            string `json:"bio"`
	AvatarURL      string `json:"avatar_url"`
	FollowerCount  int64  `json:"follower_count"`
	FollowingCount int64  `json:"following_count"`
	IsFollowing    bool   `json:"is_following"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
	Bio         *string `json:"bio" binding:"omitempty,max=500"`
}
