package xtraml

import "reflect"

var registry map[string]RegisteredType

func Register(first reflect.Type, types ...reflect.Type) {
	register(first)
	for _, t := range types {
		register(t)
	}
}

func getKey(t reflect.Type) string {
	return t.PkgPath() + "/" + t.Name()
}

func register(t reflect.Type) {
	key := getKey(t)
	_, ok := registry[key]
	if !ok {
		registry[key] = fromType(t)
	}
}

func has(t reflect.Type) bool {
	key := getKey(t)
	_, ok := registry[key]
	return ok
}

type RegisteredType struct {
}

func fromType(t reflect.Type) RegisteredType {

	return RegisteredType{}
}

type XtraMLTag struct {
	IsArray  bool
	NodeType NodeType
	Types    map[string]string
}

type NodeType int

const (
	NodeNameType NodeType = iota
	NodeObjType
	AttrType
	CharDataType
	CommentType
	ProcInstType
)
