package colmgr

import (
	"reflect"
	"fmt"
)

// this is the map that stores the roots of the collections

var collections map[uintptr]Collector
func init() {
	collections = make(map[uintptr]Collector)
}

func Init(handle interface{}, rooter Rooter) {
	p := uintptr(reflect.ValueOf(handle).Pointer())

	fmt.Printf("Calling Root to %d.\n", p)

	collections[p] = rooter.Root()
}

func Destroy(handle interface{}) {
	p := uintptr(reflect.ValueOf(handle).Pointer())
	delete(collections, p)

	fmt.Printf("Destroyed %d.\n", p)
	// FIXME: refcounting?
}

type Rooter interface {
	Root() Collector
}

type Collector interface {
	Atterer	// Cursor operator - upcoming
	MkNoder	// SCAFFOLDING OPERATOR
}

// Cursor operators:... upcoming////////////////////////////////////////////////
type Ender interface {
	End() bool
}

type Nexter interface {
	Ender
	Next()
}
type Atter interface {
	Ender
	Next() Nexter
}
type Atterer interface {
	At(uintptr) Atter
}
func At(handle interface{}, key uintptr) Atter {
	p := uintptr(reflect.ValueOf(handle).Pointer())
	return collections[p].At(key)
}
// SCAFFOLDING OPERATORS:///////////////////////////////////////////////////////
type MkNoder interface {
	MkNode(uintptr, []byte)
}
func MkNode(handle interface{}, key uintptr, val []byte) {
	p := uintptr(reflect.ValueOf(handle).Pointer())
	collections[p].MkNode(key, val)
}
