// this file was generated by gomacro command: import "unsafe"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"unsafe"
)

func init() {
	Binds["unsafe"] = map[string]Value{
	}
	Types["unsafe"] = map[string]Type{
		"Pointer":	TypeOf((*unsafe.Pointer)(nil)).Elem(),
	}
	Proxies["unsafe"] = map[string]Type{
	}
}