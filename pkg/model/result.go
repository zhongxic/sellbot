package model

import "github.com/zhongxic/sellbot/pkg/errorcode"

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

func FailedWithCode(code *errorcode.ErrorCode, args ...any) *Result[any] {
	return &Result[any]{
		Success:      false,
		ErrorCode:    code.Code(),
		ErrorMessage: code.Message(args...),
		Data:         nil,
	}
}
