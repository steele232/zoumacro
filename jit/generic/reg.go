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
 * reg.go
 *
 *  Created on May 24, 2018
 *      Author Massimiliano Ghilardi
 */

package arch

// hardware register. implementation is architecture-dependent
type RegId uint8

func (id RegId) Validate() {
	if !id.Valid() {
		errorf("invalid register: %v", id)
	}
}

// register + kind
type Reg struct {
	id   RegId
	kind Kind // defines width and signedness
}

func MakeReg(id RegId, kind Kind) Reg {
	return Reg{id: id, kind: kind}
}

// implement Arg interface
func (r Reg) UsedRegId() RegId {
	return r.id
}

func (r Reg) Kind() Kind {
	return r.kind
}

func (r Reg) Const() bool {
	return false
}

func (r Reg) Valid() bool {
	return r.id.Valid()
}

func (r Reg) Validate() {
	r.id.Validate()
}

// ===================================

type RegIds [RHi + 1]uint32 // Reg -> use count

func newRegs(ids ...RegId) *RegIds {
	var ret RegIds
	for _, id := range ids {
		ret.IncUse(id)
	}
	return &ret
}

func (rs *RegIds) InitLive() {
	*rs = alwaysLiveRegIds
}

func (rs *RegIds) IsUsed(r RegId) bool {
	return r >= RLo && r <= RHi && rs[r] != 0
}

func (rs *RegIds) IncUse(r RegId) {
	if r >= RLo && r <= RHi {
		rs[r]++
	}
}

func (rs *RegIds) DecUse(r RegId) {
	if rs.IsUsed(r) {
		rs[r]--
	}
}

// ===================================

func (asm *Asm) RegIsUsed(id RegId) bool {
	return asm.regIds.IsUsed(id)
}

func (asm *Asm) RegIncUse(id RegId) {
	asm.regIds.IncUse(id)
}

func (asm *Asm) RegDecUse(id RegId) {
	asm.regIds.DecUse(id)
}
