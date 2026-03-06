package dto

type SuggestedUserResponse struct {
	ID                  string `json:"id"`
	Username            string `json:"username"`
	DisplayName         string `json:"display_name"`
	Bio                 string `json:"bio"`
	AvatarURL           string `json:"avatar_url"`
	FollowerCount       int64  `json:"follower_count"`
	FollowingCount      int64  `json:"following_count"`
	MutualFollowerCount int64  `json:"mutual_follower_count"`
}
