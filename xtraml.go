package xtraml

import "reflect"

var registry map[string]reflect.Type

func Register(first reflect.Type, types ...reflect.Type) {
	register(first)
	for _, t := range types {
		register(t)
	}
}

func register(t reflect.Type) {
	key := t.PkgPath() + "/" + t.Name()
	registry[key] = t
}
