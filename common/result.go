package common

import (
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data gin.H       `json:"data,omitempty"`
	Body interface{} `json:"body,omitempty"`
}

func (r *Result) WithMsg(message string) *Result {
	r.Msg = message
	return r
}

func (r *Result) WithData(data gin.H) *Result {
	r.Data = data
	return r
}

func (r *Result) WithKeyValue(key string, value interface{}) *Result {
	if r.Data == nil {
		r.Data = gin.H{}
	}
	r.Data[key] = value
	return r
}

func (r *Result) WithBody(body interface{}) *Result {
	r.Body = body
	return r
}

func newResult(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func OK(msg string) *Result {
	return newResult(200, msg)
}

func Err(msg string) *Result {
	return newResult(500, msg)
}

func UnAuth(msg string) *Result {
	return newResult(401, msg)
}
