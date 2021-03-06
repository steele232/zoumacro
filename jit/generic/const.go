/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2018 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * const.go
 *
 *  Created on Jan 23, 2019
 *      Author Massimiliano Ghilardi
 */

package arch

import (
	"fmt"
	"reflect"
)

type Const struct {
	val  int64
	kind Kind
}

func (c Const) String() string {
	return fmt.Sprintf("0x%x/*%v*/", c.kind, c.val)
}

// implement Arg interface
func (c Const) UsedRegId() RegId {
	return NoRegId
}

func (c Const) Kind() Kind {
	return c.kind
}

func (c Const) Const() bool {
	return true
}

// convert Const to a different kind
func (c Const) Cast(to Kind) Const {
	val := c.val
	// sign-extend or zero-extend to 64 bits
	switch c.kind {
	case Bool:
		if val != 0 {
			// non-zero means true => convert to 1
			val = 1
		}
	case Int:
		val = int64(int(val))
	case Int8:
		val = int64(int8(val))
	case Int16:
		val = int64(int16(val))
	case Int32:
		val = int64(int32(val))
	case Int64:
		// nothing to do
	case Uint:
		val = int64(uint(val))
	case Uint8:
		val = int64(uint8(val))
	case Uint16:
		val = int64(uint16(val))
	case Uint32:
		val = int64(uint32(val))
	case Uint64:
		val = int64(uint64(val)) // should be a nop
	case Uintptr:
		val = int64(uintptr(val))
	case Float32, Float64:
		errorf("float constants not supported yet")
	default:
		errorf("invalid constant kind: %v", c)
	}
	// let caller truncate val as needed
	return Const{val: val, kind: to}
}

func MakeConst(val int64, kind Kind) Const {
	return Const{val: val, kind: kind}
}

func ConstInt64(val int64) Const {
	return Const{val: val, kind: Int64}
}

func ConstInterface(ival interface{}) Const {
	v := reflect.ValueOf(ival)
	kind := Kind(v.Kind())
	var val int64
	switch kind {
	case Bool:
		if v.Bool() {
			val = 1
		}
	case Int, Int8, Int16, Int32, Int64:
		val = v.Int()
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		val = int64(v.Uint())
	case Float32, Float64:
		errorf("float constants not supported yet")
	default:
		errorf("invalid constant kind: %v", kind)
	}
	return Const{val: val, kind: kind}
}
