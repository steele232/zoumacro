/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http//www.gnu.org/licenses/>.
 *
 * util.go
 *
 *  Created on May 07, 2017
 *      Author Massimiliano Ghilardi
 */

package xreflect

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"
)

func debugf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	fmt.Printf("// debug: %s\n", str)
}

func gtypeToKind(gtype types.Type) reflect.Kind {
	gtype = gtype.Underlying()
	var kind reflect.Kind
	switch gtype := gtype.(type) {
	case *types.Array:
		kind = reflect.Array
	case *types.Basic:
		kind = gbasickindToKind(gtype.Kind())
	case *types.Chan:
		kind = reflect.Chan
	case *types.Signature:
		kind = reflect.Func
	case *types.Interface:
		kind = reflect.Interface
	case *types.Map:
		kind = reflect.Map
	case *types.Pointer:
		kind = reflect.Ptr
	case *types.Slice:
		kind = reflect.Slice
	case *types.Struct:
		kind = reflect.Struct
	// case *types.Named: // impossible, handled above
	default:
		errorf("unsupported types.Type: %v", gtype)
	}
	// debugf("gtypeToKind(%T) -> %v", gtype, kind)
	return kind
}

func gbasickindToKind(gkind types.BasicKind) reflect.Kind {
	var kind reflect.Kind
	switch gkind {
	case types.Bool:
		kind = reflect.Bool
	case types.Int:
		kind = reflect.Int
	case types.Int8:
		kind = reflect.Int8
	case types.Int16:
		kind = reflect.Int16
	case types.Int32:
		kind = reflect.Int32
	case types.Int64:
		kind = reflect.Int64
	case types.Uint:
		kind = reflect.Uint
	case types.Uint8:
		kind = reflect.Uint8
	case types.Uint16:
		kind = reflect.Uint16
	case types.Uint32:
		kind = reflect.Uint32
	case types.Uint64:
		kind = reflect.Uint64
	case types.Uintptr:
		kind = reflect.Uintptr
	case types.Float32:
		kind = reflect.Float32
	case types.Float64:
		kind = reflect.Float64
	case types.Complex64:
		kind = reflect.Complex64
	case types.Complex128:
		kind = reflect.Complex128
	case types.String:
		kind = reflect.String
	case types.UnsafePointer:
		kind = reflect.UnsafePointer
	default:
		errorf("unsupported types.BasicKind: %v", gkind)
	}
	return kind
}

func dirToGdir(dir reflect.ChanDir) types.ChanDir {
	var gdir types.ChanDir
	switch dir {
	case reflect.RecvDir:
		gdir = types.RecvOnly
	case reflect.SendDir:
		gdir = types.SendOnly
	case reflect.BothDir:
		gdir = types.SendRecv
	}
	return gdir
}

func toReflectTypes(ts []Type) []reflect.Type {
	rts := make([]reflect.Type, len(ts))
	for i, t := range ts {
		rts[i] = t.ReflectType()
	}
	return rts
}

func toGoParam(t Type) *types.Var {
	return types.NewParam(token.NoPos, nil, "", t.GoType())
}

func toGoParams(ts []Type) []*types.Var {
	vars := make([]*types.Var, len(ts))
	for i, t := range ts {
		vars[i] = toGoParam(t)
	}
	return vars
}

func toGoTuple(ts []Type) *types.Tuple {
	vars := toGoParams(ts)
	return types.NewTuple(vars...)
}