package dto

type HttpResponse[TData any] struct {
	HttpStatus int    `json:"http_status"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Data       TData  `json:"data,omitempty"`
}

type OffsetPageResult[TEntity any] struct {
	Items []TEntity `json:"items"`
	Page  int8      `json:"page"`
	Size  int8      `json:"size"`
	Total int64     `json:"total"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
