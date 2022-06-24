package dto

type ErrorResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Data interface{}
	Code int
}
