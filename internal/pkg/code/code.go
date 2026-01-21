package code

import (
	"net/http"

	"github.com/HappyLadySauce/errors"
	"github.com/novalagung/gubrak"
)

var codeMessage = map[int]string{}

// Message returns the registered external message for a given error code.
// If the code is unknown, it returns an empty string.
func Message(c int) string {
	return codeMessage[c]
}

// ErrCode implements `github.com/marmotedu/errors`.Coder interface.
type ErrCode struct {
	// C refers to the code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

var _ errors.Coder = &ErrCode{}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return coder.Ref
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

// nolint: unparam
func register(code int, httpStatus int, message string, refs ...string) {
	// 允许的 HTTP 状态码集合，需要与各业务错误码保持一致
	allowedStatusCodes := []int{
		http.StatusOK,                  // 200
		http.StatusBadRequest,          // 400
		http.StatusUnauthorized,        // 401
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusConflict,            // 409
		http.StatusInternalServerError, // 500
	}

	found, _ := gubrak.Includes(allowedStatusCodes, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 409, 500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	// Cache message for reuse across layers (avoid duplicating strings).
	codeMessage[code] = message
	errors.MustRegister(coder)
}
