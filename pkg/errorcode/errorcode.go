package errorcode

import (
	"fmt"
	"strings"
)

var (
	ParamsError = New("ERR0001", "params error: {0}")
	SystemError = New("ERR0002", "system error: {0}")
)

type ErrorCode struct {
	code    string
	message string
}

func New(code, message string) *ErrorCode {
	return &ErrorCode{
		code:    code,
		message: message,
	}
}

func (c ErrorCode) Code() string {
	return c.code
}

func (c ErrorCode) Message(args ...any) string {
	formatted := c.message
	for i := range args {
		placeholder := fmt.Sprintf("{%d}", i)
		replaced := fmt.Sprintf("%v", args[i])
		formatted = strings.ReplaceAll(formatted, placeholder, replaced)
	}
	return formatted
}
