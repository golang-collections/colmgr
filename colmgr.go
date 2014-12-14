package colmgr

import (
	"github.com/anlhord/generic"
	"github.com/anlhord/generic/low"
	genref "github.com/anlhord/generic/reflect"
	"reflect"
//	"fmt"
)

// this is the map that stores the roots of the collections

var collections map[uintptr]Collector

func init() {
	collections = make(map[uintptr]Collector)
}

func Init(handle interface{}, rooter Rooter) {
	p := uintptr(reflect.ValueOf(handle).Pointer())

	//	fmt.Printf("Calling Root to %d.\n", p)

	if colmgr_debug && collections[p] != nil {
		panic("Reinitialized collection")
	}

	collections[p] = rooter.Root()
}

func Destroy(handle interface{}) {
	p := uintptr(reflect.ValueOf(handle).Pointer())

	collections[p].Destroy()
	delete(collections, p)

	//	fmt.Printf("Destroyed %d.\n", p)
	// FIXME: refcounting?
}

////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////

func Inserts(gkey uintptr, slice interface{}, u MkNoder) {
	if colmgr_debug && !low.S(slice) {
		panic("Insert slice is not a slice")
	}

	sll := genref.Stuff(slice)	// FIXME: use low.T here

	if colmgr_gen_debug {	// FIXME: make low.T work

		sl := *low.T(slice)

		if &(sll[0]) != &((*sl)[0]) {
			panic("low.T bug")
		}
		if len(sll) != len((*sl)) {
			panic("low.T len bug")
		}
		if cap(sll) != cap((*sl)) {
			panic("low.T cap bug")
		}
	}

//	fmt.Printf("Inserts %v.\n", &((sll)[0]))
	u.MkNode(gkey, &sll)
}

func Inserti(gkey uintptr, iface interface{}, u MkNoder) {
	u.MkNode(gkey, low.Y(iface))
}

func Insert(gkey uintptr, gvalue interface{}, u MkNoder) {
	if low.S(gvalue) {
		Inserts(gkey, gvalue, u)
	} else {
		Inserti(gkey, gvalue, u)
	}
}

////////////////////////////////////////////////////////////////////////////////

func Append(gvalue interface{}, p Pusher) {
	if low.S(gvalue) {
		Appends(gvalue, p)
	} else {
		Appendi(gvalue, p)
	}
}

func Appendi(iface interface{}, p Pusher) {
	p.Push(low.Y(iface))
}

func Appends(slice interface{}, p Pusher) {
	p.Push(*low.T(slice))
}

////////////////////////////////////////////////////////////////////////////////

type Rooter interface {
	Root() Collector
}

type Collector interface {
	Atterer
	MkNoder
	Dumper
	Destroyer
}

type Destroyer interface {
	Destroy()
}

// Cursor operators:////////////////////////////////////////////////
type Ender interface {
	End() bool
}

type UpdMapper interface {
	Map() *generic.Value
	Upd(*generic.Value)
}

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
	Pusher
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

type Pusher interface {
	Push(*generic.Value)
}

type Fixer interface {
	Fix()
}

type MkNoder interface {
	MkNode(uintptr, *generic.Value)
}
/*
// scaffolding operation
func Insert(handle interface{}, key uintptr, val *generic.Value) {
	if key >= End {
		panic("Key -1 is end. Use smaller")
	}

	p := uintptr(reflect.ValueOf(handle).Pointer())
	collections[p].MkNode(key, val)
}
*/
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
