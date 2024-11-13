package model

type Result[T any] struct {
	Success      bool   `json:"success"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Data         T      `json:"data"`
}

func Success() *Result[any] {
	return &Result[any]{
		Success:      true,
		ErrorCode:    "",
		ErrorMessage: "",
		Data:         nil,
	}
}

func SuccessWithData[T any](data T) *Result[T] {
	return &Result[T]{
		Success:      true,
		ErrorCode:    "",
		ErrorMessage: "",
		Data:         data,
	}
}

func Failed(code string) *Result[any] {
	return &Result[any]{
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: "",
		Data:         nil,
	}
}

func FailedWithMessage(code, message string) *Result[any] {
	return &Result[any]{
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: message,
		Data:         nil,
	}
}

func FailedWithData[T any](code, message string, data T) *Result[T] {
	return &Result[T]{
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: message,
		Data:         data,
	}
}
