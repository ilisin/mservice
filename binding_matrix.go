package mservice

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// binding tag with section "req"
//	'': param's name,default value is the lower case of field name
//	notnull : the value can't be empty , default null
//	default : default value , if not contain in request
//-	min : request value must bigger or equal than the value
//	max : request value must small or equal than the value
//	enum : request value must is one of the enum value
//	regex : regex , between in / ... /
type bindingMatrix struct {
	Field   reflect.StructField
	Name    string
	CanNull bool
	Default string
	Min     string
	Max     string
	Enum    string
}

func newBindingMatrix(p reflect.StructField) *bindingMatrix {
	bm := &bindingMatrix{}
	tagstr := p.Tag.Get(BINDING_TAG)
	tags := strings.Split(tagstr, " ")
	bm.Field = p
	bm.CanNull = true
	for _, t := range tags {
		if len(t) > 2 && strings.HasPrefix(t, "'") && strings.HasSuffix(t, "'") {
			bm.Name = string(t[1 : len(t)-1])
		} else {
			t = strings.ToLower(t)
			if t == "notnull" {
				bm.CanNull = false
			} else if strings.HasPrefix(t, "default(") && strings.HasSuffix(t, ")") {
				bm.Default = string(t[8 : len(t)-1])
			} else if strings.HasPrefix(t, "min(") && strings.HasSuffix(t, ")") {
				bm.Min = string(t[4 : len(t)-1])
			} else if strings.HasPrefix(t, "max(") && strings.HasSuffix(t, ")") {
				bm.Max = string(t[4 : len(t)-1])
			} else if strings.HasPrefix(t, "enum(") && strings.HasSuffix(t, ")") {
				bm.Enum = string(t[5 : len(t)-1])
			}
		}
	}
	if len(bm.Name) == 0 {
		bm.Name = strings.ToLower(p.Name)
	}
	return bm
}

func (bm *bindingMatrix) ValueTo(sv reflect.Value, value string) (reflect.Value, error) {
	switch bm.Field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(value) == 0 && !bm.CanNull {
			return sv, fmt.Errorf("param %v's value cann't be null", bm.Name)
		}
		if len(value) == 0 && bm.CanNull && len(bm.Default) > 0 {
			value = bm.Default
		}
		if len(value) == 0 && bm.CanNull && len(bm.Default) == 0 {
			break
		}
		if i64, err := strconv.ParseInt(value, 10, 0); err != nil {
			return sv, fmt.Errorf("param %v's value isn't a integer", bm.Name)
		} else {
			if len(bm.Min) > 0 {
				if min, err := strconv.ParseInt(bm.Min, 10, 0); err == nil {
					if i64 < min {
						return sv, fmt.Errorf("param %v's value cann't smaller than %v", bm.Name, min)
					}
				}
			}
			if len(bm.Max) > 0 {
				if max, err := strconv.ParseInt(bm.Max, 10, 0); err == nil {
					if i64 > max {
						return sv, fmt.Errorf("param %v's value cann't be bigger than %v", bm.Name, max)
					}
				}
			}
			if len(bm.Enum) > 0 {
				bes := strings.Split(bm.Enum, ";")
				b := false
				for _, s := range bes {
					if ei64, err := strconv.ParseInt(s, 10, 0); err == nil {
						if i64 == ei64 {
							b = true
							break
						}
					}
				}
				if !b {
					return sv, fmt.Errorf("param %v's value must be in [%v]", bm.Name, strings.Join(bes, ","))
				}
			}
			sv.SetInt(i64)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(value) == 0 && !bm.CanNull {
			return sv, fmt.Errorf("param %v's value cann't be null", bm.Name)
		}
		if len(value) == 0 && bm.CanNull && len(bm.Default) > 0 {
			value = bm.Default
		}
		if len(value) == 0 && bm.CanNull && len(bm.Default) == 0 {
			break
		}
		if ui64, err := strconv.ParseUint(value, 10, 0); err != nil {
			return sv, fmt.Errorf("param %v's value isn't a positive integer", bm.Name)
		} else {
			if len(bm.Min) > 0 {
				if min, err := strconv.ParseUint(bm.Min, 10, 0); err == nil {
					if ui64 < min {
						return sv, fmt.Errorf("param %v's value cann't smaller than %v", bm.Name, min)
					}
				}
			}
			if len(bm.Max) > 0 {
				if max, err := strconv.ParseUint(bm.Max, 10, 0); err == nil {
					if ui64 > max {
						return sv, fmt.Errorf("param %v's value cann't be bigger than %v", bm.Name, max)
					}
				}
			}
			sv.SetUint(ui64)
		}
	case reflect.String:
	case reflect.Bool:
	case reflect.Array:
	case reflect.Slice:
	case reflect.Float32, reflect.Float64:
	case reflect.Map:
	case reflect.Struct:
	case reflect.Ptr:
	}
	return sv, nil
}
