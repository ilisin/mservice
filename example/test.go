package main

import (
	"errors"
	"fmt"
	"reflect"
)

func test() {
	var err1, err2 error
	err1 = errors.New("a error")
	err2 = errors.New("b error")

	t1 := reflect.TypeOf(err1)
	t2 := reflect.TypeOf(err2)

	fmt.Printf("%v %v\n", t1, t2)
	fmt.Printf("%v %v\n", t1.Kind(), t2.Kind())
	fmt.Printf("%v\n", t1 == t2)
}
