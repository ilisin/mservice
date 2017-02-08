package mservice

import (
	"reflect"
	"fmt"
)

type ErrorCode int

const (
	ErrSuccess ErrorCode = 0		// success
	ErrValidate ErrorCode = 1001	// prase param error
	ErrApi = 1002	// api error
	ErrCode ErrorCode = 1010		// route function error
)

type Response struct {
	Error ErrorCode		`json:"error"`
	Message string 		`json:"message"`
	Data interface{}	`json:"data"`
}

func mapToResponse(vs []reflect.Value) *Response{
	resp := &Response{
		Error:ErrSuccess,
	}
	if len(vs) == 0 {
		resp.Error = ErrCode
		return resp
	}
	last := vs[len(vs)-1]
	if !last.CanInterface() {
		resp.Error = ErrCode
	}else{
		if last.IsNil() {
			if len(vs) > 1 {
				if vs[0].CanInterface() {
					resp.Data = vs[0].Interface()
				}else{
					resp.Error = ErrCode
				}
			}else{
				resp.Data = nil
			}
		}else{
			resp.Error = ErrApi
			resp.Message = fmt.Sprintf("%v",last.Interface())
		}
	}
	return resp
}