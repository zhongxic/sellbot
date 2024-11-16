package errorcode

import (
	"testing"
)

func TestErrorCodeMessage(t *testing.T) {
	code := "ERR0000"
	message := "some err: {0}"
	errorCode := New(code, message)

	if errorCode.Code() != code {
		t.Errorf("expected error code [%v] actual [%v]", code, errorCode.Code())
	}
	if errorCode.Message() != message {
		t.Errorf("expected error message [%v] actual [%v]", message, errorCode.Message())
	}
	expectedMessage := "some err: internal server error"
	if errorCode.Message("internal server error") != expectedMessage {
		t.Errorf("expected error message [%v] actual [%v]", expectedMessage, errorCode.Message("internal server error"))
	}
}
