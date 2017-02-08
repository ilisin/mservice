// request data model binding
//
//type ReqStatus string
//
//const (
//	ReqStatusOn ReqStatus = "S001"
//	ReqStatusOff ReqStatus = "S002"
//)
//
//type TestModel struct {
//	Keyword string `req:"'keyword' notnull"`
//	Size int `req:"'size' default(0)"`
//	Page int `req:"'page' default(0)"`
//	Count int `req:"min(1) max(1)"`
//	Status ReqStatus `req:"enum(S001,S002)"`
//	Email string `req:"regex/^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$/ error(email error)"`
//}
// binding tag with section "req"
//	'': param's name,default value is the lower case of field name
//	notnull : the value can't be empty , default null
//	default : default value , if not contain in request
//-	min : request value must bigger or equal than the value
//	max : request value must small or equal than the value
//	enum : request value must is one of the enum value
//	regex : regex , between in / ... /
//	error : if validate fail ,and return this error message
package mservice

import (
	"reflect"
	"errors"
)

type Binding struct {
	Ctx *Context
	Typ reflect.Type
	Val reflect.Value
}

func newBinding(ctx *Context,typ reflect.Type,val reflect.Value) *Binding{
	return &Binding{
		Ctx:ctx,
		Typ:typ,
	}
}

func (b *Binding)MapTo() (reflect.Value,error){
	v := reflect.New(b.Typ)
	if b.Typ.Kind() != reflect.Struct {
		return v,errors.New("reflect type error")
	}
	b.fillValue(v)
	return v,nil
}

func (b *Binding)fillValue(v reflect.Value){

}