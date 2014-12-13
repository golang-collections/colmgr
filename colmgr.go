package colmgr

import (
	"github.com/anlhord/generic/low"
	"github.com/anlhord/generic"
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

func Puts(slice interface{}, u UpdMapper) {
	if colmgr_debug && !low.S(slice) {
		panic("Put slice is not a slice")
	}
	u.Upd(*low.T(slice))
}

func Gets(slice interface{}, u UpdMapper) {
	if colmgr_debug && !low.S(slice) {
		panic("Get slice is not a slice")
	}
	*low.T(slice) = u.Map()
}

func Puti(iface interface{}, u UpdMapper) {
	if colmgr_debug && low.S(iface) {
		panic("Get iface is not an iface")
	}
	u.Upd(low.Y(iface))
}

func Geti(iface interface{}, u UpdMapper) {
	if colmgr_debug && low.S(iface) {
		panic("Get iface is not an iface")
	}
	*low.I(iface) = u.Map()
}

func Put(gvalue interface{}, u UpdMapper) {
	if low.S(gvalue) {
		Puts(gvalue, u)
	} else {
		Puti(gvalue, u)
	}
}

func Get(gvalue interface{}, u UpdMapper) {
	if low.S(gvalue) {
		Gets(gvalue, u)
	} else {
		Geti(gvalue, u)
	}
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
	Map() *generic.Value
	Upd(*generic.Value)
}

type Pad uint16
const Begin = uintptr(0)
const Root = ^Begin
const End = ^uintptr(1)

type Nexter interface {
	Atterer
	Ender
	Next()
}
type Atter interface {
	Atterer
	UpdMapper
	Ender
	Nexterer
	MkNoder
	Appender
	Fixer
}
type Atterer interface {
	At(uintptr) Atter
}

type Nexterer interface {
	Next() Nexter
}

func At(handle interface{}, key uintptr) Atter {
	p := uintptr(reflect.ValueOf(handle).Pointer())
	return collections[p].At(key)
}

type Appender interface {
	Append(generic.Value)
}

type Fixer interface {
	Fix()
}

// SCAFFOLDING OPERATORS:/DO NOT USE IN PRODUCTION FOR TESTING PURPOSE ONLY/////
type MkNoder interface {
	MkNode(uintptr, generic.Value)
}

func MkNode(handle interface{}, key uintptr, val generic.Value) {
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
