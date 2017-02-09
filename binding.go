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
	"errors"
	"reflect"
)

const BINDING_TAG = "req"

type Binding struct {
	//Ctx *Context
	Typ            reflect.Type
	IsContext      bool
	BindingMatrix []bindingMatrix
	//Val reflect.Value
}

func newBinding(typ reflect.Type) *Binding {
	return &Binding{
		//Ctx: ctx,
		Typ: typ,
	}
}

func (b *Binding) MapTo(ctx *Context) (reflect.Value, error) {
	v := reflect.New(b.Typ)
	//v := b.Val
	if b.Typ.Kind() != reflect.Struct {
		return v, errors.New("reflect type error")
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < b.Typ.NumField(); i++ {
		bm := newBindingMatrix(b.Typ.Field(i))
		vv := v.Field(i)
		_, err := bm.ValueTo(vv, ctx.QueryParam(bm.Name))
		if err != nil {
			return v.Addr(), err
		}
	}
	return v.Addr(), nil
}

//
//func (b *Binding) fillValue(v reflect.Value) error {
//	typ := b.Typ
//	val := b.Val
//	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
//		typ = typ.Elem()
//		val = val.Elem()
//	}
//	if typ.Kind() != reflect.Struct {
//		return errors.New("param's type error")
//	}
//
//	return nil
//}
