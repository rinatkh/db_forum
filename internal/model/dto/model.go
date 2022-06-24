package dto

import (
	"github.com/rinatkh/db_forum/internal/model/core"
	"time"
)

type CreateForumRequest struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type GetForumRequest struct {
	Slug string `path:"slug"`
}

type GetForumThreadsRequest struct {
	Slug  string `path:"slug"`
	Limit int64  `query:"limit"`
	Since string `query:"since"`
	Desc  bool   `query:"desc"`
}

type GetForumUsersRequest struct {
	Slug  string `path:"slug"`
	Limit int64  `query:"limit"`
	Since string `query:"since"`
	Desc  bool   `query:"desc"`
}
type Post struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type PostDetails struct {
	Author *core.User   `json:"author,omitempty"`
	Thread *core.Thread `json:"thread,omitempty"`
	Post   *core.Post   `json:"post"`
	Forum  *core.Forum  `json:"forum,omitempty"`
}

type GetPostDetailsRequest struct {
	Related string `query:"related"`
	ID      int64  `path:"id"`
}

type EditPostRequest struct {
	Message string `json:"message"`
	ID      int64  `path:"id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Data interface{}
	Code int
}
type CreateThreadRequest struct {
	Author  string    `json:"author"`
	Forum   string    `path:"slug"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Created time.Time `json:"created,omitempty"`
}

type EditVoteRequest struct {
	Voice    int64  `json:"voice"`
	Nickname string `json:"nickname"`
}

type EditThreadRequest struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}
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
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type EditUserProfileRequest struct {
	Nickname string `path:"nickname"`
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}
