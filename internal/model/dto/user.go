package dto

type CreateUserRequest struct {
	Nickname string `path:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type GetUserProfileRequest struct {
	Nickname string `path:"nickname"`
}

type GetUserProfileResponse struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type EditUserProfileRequest struct {
	Nickname string `path:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
