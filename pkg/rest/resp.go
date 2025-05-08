package rest

import (
	"fmt"
	"net/http"
)

const UnknownErr = -1

type HTTPError struct {
	HttpStatus int         `json:"-"`
	Code       int         `json:"code"`
	Err        error       `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (resp HTTPError) Error() string {
	if resp.Err == nil {
		return ""
	}
	return fmt.Sprintf("err code:%d\nstatus code:%d\nerr message:%s", resp.Code, resp.HttpStatus, resp.Err)
}

func Json(data interface{}) HTTPError {
	return HTTPError{
		HttpStatus: http.StatusOK,
		Data:       data,
	}
}

func ErrorWithStatus(status int, err error) HTTPError {
	return HTTPError{
		HttpStatus: status,
		Code:       UnknownErr,
		Err:        err,
	}
}

type CodeError struct {
	Code int   `json:"code"`
	Err  error `json:"error,omitempty"`
}

func (err CodeError) Error() string {
	if err.Err == nil {
		return ""
	}
	return fmt.Sprintf("err code:%d\nerr message:%s", err.Code, err.Err)
}
