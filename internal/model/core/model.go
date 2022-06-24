package core

import "time"

type Forum struct {
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts"`
	Title   string `json:"title"`
	Threads int64  `json:"threads"`
}

type Post struct {
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	ID       int64     `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Thread   int64     `json:"thread"`
	Created  time.Time `json:"created"`
}

type ServiceInfo struct {
	Forum  int64 `json:"forum"`
	Thread int64 `json:"thread"`
	User   int64 `json:"user"`
	Post   int64 `json:"post"`
}

type Thread struct {
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int64     `json:"votes"`
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

type User struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type Vote struct {
	Nickname string
	ThreadID int64
	Voice    int64
}
