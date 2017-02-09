package mservice

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const META_PREFIX = "MS"

type HTTPMethod string

const (
	CONNECT HTTPMethod = echo.CONNECT
	DELETE  HTTPMethod = echo.DELETE
	GET     HTTPMethod = echo.GET
	HEAD    HTTPMethod = echo.HEAD
	OPTIONS HTTPMethod = echo.OPTIONS
	PATCH   HTTPMethod = echo.PATCH
	POST    HTTPMethod = echo.POST
	PUT     HTTPMethod = echo.PUT
	TRACE   HTTPMethod = echo.TRACE
)

// MSHandlerFunc defines a function to server HTTP requests.
type MSHandlerFunc func(*Context) error

// return http method,http route path and http request function
type NSAutoHandler func() (HTTPMethod, string, interface{})

type Model struct {
	// route path prefix
	Prefix string
	// struct model
	Prototype interface{}
}

type Meta struct {
	Method      HTTPMethod
	Path        string
	Description string
	Handler     MSHandlerFunc
}

func NewModel(prefix string, proto interface{}) *Model {
	return &Model{
		Prefix:    prefix,
		Prototype: proto,
	}
}

func canExport(funcName string) bool {
	if string(funcName[:1]) > "Z" || string(funcName[:1]) < "A" {
		return false
	}
	return true
}

func autoHandlerCheck(v reflect.Value, m reflect.Method) (*Meta, bool) {
	if !canExport(m.Name) {
		return nil, false
	}
	// none param in is void param
	if m.Type.NumIn() != 1 || m.Type.NumOut() != 4 {
		return nil, false
	}
	if m.Type.Out(0).Kind() != reflect.String || m.Type.Out(1).Kind() != reflect.String ||
		m.Type.Out(2).Kind() != reflect.String {
		return nil, false
	}
	rets := v.Call([]reflect.Value{})
	meta := &Meta{}
	ok := false
	if meta.Method, ok = rets[0].Interface().(HTTPMethod); !ok {
		return nil, false
	}
	if meta.Path, ok = rets[1].Interface().(string); !ok {
		return nil, false
	}
	if meta.Description, ok = rets[2].Interface().(string); !ok {
		return nil, false
	}
	h, err := wrapperHandler(rets[3])
	if err != nil {
		logrus.Error(err)
		return nil, false
	}
	meta.Handler = h
	return meta, true
}

func wrapperHandler(v reflect.Value) (MSHandlerFunc, error) {
	if v.Kind() != reflect.Func {
		return nil, errors.New("the value isn't a function")
	}
	t := reflect.TypeOf(v.Interface())
	//ints := make([]reflect.Type, t.NumIn())
	bindings := make([]*Binding, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		tt := t.In(i)
		if tt.Kind() != reflect.Ptr || tt.Elem().Kind() != reflect.Struct {
			return nil, errors.New("the return function's in param isn't a struct pointer")
		}
		//ints[i] = tt.Elem()
		val := reflect.New(tt.Elem())
		binding := newBinding(tt.Elem())
		if _, ok := val.Interface().(*Context); ok {
			binding.IsContext = true
		}
		bindings[i] = binding
	}
	// (error) ro (*{},error)
	if t.NumOut() != 1 && t.NumOut() != 2 {
		return nil, errors.New("the return function's out param isn't a error or pointer and error")
	}

	return func(c *Context) error {
		pins := make([]reflect.Value, len(bindings))
		for i, bd := range bindings {
			// context
			if bd.IsContext {
				val := reflect.New(bd.Typ)
				nc, _ := val.Interface().(*Context)
				nc.Context = c
				pins[i] = reflect.ValueOf(nc)
			} else {
				vv, err := bd.MapTo(c)
				// map error
				if err != nil {
					return c.JSON(http.StatusOK, &Response{
						Error:   ErrValidate,
						Message: fmt.Sprintf("%v", err),
						Data:    nil,
					})
				}
				pins[i] = vv
			}
		}
		rets := v.Call(pins)
		response := mapToResponse(rets)
		return c.JSON(http.StatusOK, response)
	}, nil
}

func (m *Model) ReadHandlers() (metas []*Meta) {
	metas = make([]*Meta, 0)
	typ := reflect.TypeOf(m.Prototype)
	val := reflect.ValueOf(m.Prototype)
	if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < val.NumMethod(); i++ {
		vval := val.Method(i)
		ttyp := typ.Method(i)
		if !strings.HasPrefix(ttyp.Name, META_PREFIX) {
			continue
		}
		meta, ok := autoHandlerCheck(vval, ttyp)
		if !ok {
			continue
		}
		metas = append(metas, meta)
	}
	return
}
