package dto

type HttpResponse[TData any] struct {
	HttpStatus int    `json:"http_status"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Data       TData  `json:"data"`
}

type Page[T any] struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Items []T `json:"items"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
