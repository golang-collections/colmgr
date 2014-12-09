package colmgr

import (
//	"fmt"
	"reflect"
)

// this is the map that stores the roots of the collections

var collections map[uintptr]Collector

func init() {
	collections = make(map[uintptr]Collector)
}

func Init(handle interface{}, rooter Rooter) {
	p := uintptr(reflect.ValueOf(handle).Pointer())

//	fmt.Printf("Calling Root to %d.\n", p)

	collections[p] = rooter.Root()
}

func Destroy(handle interface{}) {
	p := uintptr(reflect.ValueOf(handle).Pointer())

	collections[p].Destroy()
	delete(collections, p)

//	fmt.Printf("Destroyed %d.\n", p)
	// FIXME: refcounting?
}

type Rooter interface {
	Root() Collector
}

type Collector interface {
	Atterer // Cursor operator - upcoming
	MkNoder // SCAFFOLDING OPERATOR
	Dumper	// dumps the collection to a stdout
	Destroyer // this should trigger a collection destruction
}

type Destroyer interface {
	Destroy()
}
// Cursor operators:... upcoming////////////////////////////////////////////////
type Ender interface {
	End() bool
}

type UpdMapper interface {
	Map() []byte
	Upd([]byte)
}

const Collection = 0
const Begin = uintptr(0)
const Root = ^Begin
const End = ^uintptr(1)

type Nexter interface {
	UpdMapper
	Ender
	Next()
}
type Atter interface {
	UpdMapper
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

// SCAFFOLDING OPERATORS:/DO NOT USE IN PRODUCTION FOR TESTING PURPOSE ONLY/////
type MkNoder interface {
	MkNode(uintptr, []byte)
}

func MkNode(handle interface{}, key uintptr, val []byte) {
	if key >= End {
		panic("Key -1 is end. Use smaller")
	}

//	print("MKNODE()\n")

	p := uintptr(reflect.ValueOf(handle).Pointer())
	collections[p].MkNode(key, val)
}

type Dumper interface {
	Dump(byte)
}

func Dump(handle interface{}, format byte) {
	if format > 0 {
		panic("Unsupported format")
	}
	p := uintptr(reflect.ValueOf(handle).Pointer())
	collections[p].Dump(format)
}
