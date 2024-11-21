package common

type Response[T any] struct {
	Success bool          `json:"success"`
	Data    *T            `json:"data"`
	Error   *ErrorMessage `json:"error"`
}

type PagingSuccessResponse[T any] struct {
	Success   bool    `json:"success"`
	Data      *T      `json:"data"`
	TotalPage *uint64 `json:"totalPage"`
}

type ErrorMessage struct {
	Errors  *string `json:"errors"`
	Message *string `json:"message"`
}

func NewSuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		Success: true,
		Data:    &data,
		Error:   nil,
	}
}

func NewErrorResponse(err string, message string) Response[any] {
	return Response[any]{
		Success: false,
		Data:    nil,
		Error: &ErrorMessage{
			Errors:  &err,
			Message: &message,
		},
	}
}

func NewPagingSuccessResponse[T any](data T, totalPage uint64) PagingSuccessResponse[T] {
	return PagingSuccessResponse[T]{
		Success:   true,
		Data:      &data,
		TotalPage: &totalPage,
	}
}
