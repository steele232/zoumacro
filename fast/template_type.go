/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017-2018 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * template_type.go
 *
 *  Created on Jun 06, 2018
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	"bytes"
	"go/ast"
	"go/token"
	r "reflect"

	"github.com/steele232/zoumacro/base"
	"github.com/steele232/zoumacro/base/output"
	xr "github.com/steele232/zoumacro/xreflect"
)

// a template type declaration.
// either general, or partially specialized or fully specialized
type TemplateTypeDecl struct {
	Decl   ast.Expr   // type declaration body. use an ast.Expr because we will compile it with Comp.Type()
	Alias  bool       // true if declaration is an alias: 'type Foo = ...'
	Params []string   // template param names
	For    []ast.Expr // for partial or full specialization
}

type TemplateType struct {
	Master    TemplateTypeDecl            // master (i.e. non specialized) declaration
	Special   map[string]TemplateTypeDecl // partially or fully specialized declarations. key is TemplateTypeDecl.For converted to string
	Instances map[I]xr.Type               // cache of instantiated types. key is [N]interface{}{T1, T2...}
}

func (t *TemplateType) String() string {
	if t == nil {
		return "<nil>"
	}
	var buf bytes.Buffer // strings.Builder requires Go >= 1.10
	buf.WriteString("template[")
	decl := t.Master
	for i, param := range decl.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param)
	}
	buf.WriteString("] type ")
	if decl.Alias {
		buf.WriteString("= ")
	}
	var str string
	switch e := decl.Decl.(type) {
	case *ast.ArrayType:
		if e.Len == nil {
			str = "slice"
		} else {
			str = "array"
		}
	case *ast.ChanType:
		str = "chan"
	case *ast.FuncType:
		str = "func"
	case *ast.InterfaceType:
		str = "interface"
	case *ast.MapType:
		str = "map"
	case *ast.StructType:
		str = "struct"
	default:
		(*output.Stringer).Fprintf(nil, &buf, "%v", decl.Decl)
	}
	buf.WriteString(str)
	return buf.String()
}

// DeclTemplateType stores a template type declaration
// for later instantiation
func (c *Comp) DeclTemplateType(spec *ast.TypeSpec) {

	lit, _ := spec.Type.(*ast.CompositeLit)
	if lit == nil {
		c.Errorf("invalid template type declaration: expecting an *ast.CompositeLit, found %T: %v",
			spec.Type, spec)
	}
	expr := lit.Type
	if _, ok := expr.(*ast.CompositeLit); ok {
		c.Errorf("invalid template type declaration: expecting an *ast.CompositeLit, found &ast.CompositeLit{Type: &ast.CompositeLit{}}: %v",
			spec)
	}
	params, fors := c.templateParams(lit.Elts, "type", spec)

	tdecl := TemplateTypeDecl{
		Decl:   lit.Type,
		Alias:  spec.Assign != token.NoPos,
		Params: params,
		For:    fors,
	}
	name := spec.Name.Name

	if len(fors) == 0 {
		// master (i.e. not specialized) declaration
		if len(params) == 0 {
			c.Errorf("cannot declare template type with zero template parameters: %v", spec)
		}

		bind := c.NewBind(name, TemplateTypeBind, c.TypeOfPtrTemplateType())
		// a template type declaration has no runtime effect:
		// it merely creates the bind for on-demand instantiation by other code

		bind.Value = &TemplateType{
			Master:    tdecl,
			Special:   make(map[string]TemplateTypeDecl),
			Instances: make(map[I]xr.Type),
		}
		return
	}

	// partially or fully specialized declaration
	bind := c.Binds[name]
	if bind == nil {
		c.Errorf("undefined identifier: %v", name)
	}
	typ, ok := bind.Value.(*TemplateType)
	if !ok {
		c.Errorf("symbol is not a template type, cannot declare type specializations on it: %s // %v", name, bind.Type)
	}
	key := c.Globals.Sprintf("%v", &ast.IndexExpr{X: spec.Name, Index: &ast.CompositeLit{Elts: fors}})
	if len(typ.Master.Params) != len(fors) {
		c.Errorf("template type specialization for %d parameters, expecting %d: %s", len(fors), len(typ.Master.Params), key)
	}
	if _, ok := typ.Special[key]; ok {
		c.Warnf("redefined template type specialization: %s", key)
	}
	typ.Special[key] = tdecl
}

// TemplateType compiles a template type name#[T1, T2...] instantiating it if needed.
func (c *Comp) TemplateType(node *ast.IndexExpr) xr.Type {
	maker := c.templateMaker(node, TemplateTypeBind)
	if maker == nil {
		return nil
	}
	typ := maker.ifun.(*TemplateType)
	key := maker.ikey

	g := &c.Globals
	debug := g.Options&base.OptDebugTemplate != 0

	instance, _ := typ.Instances[key]
	if instance != nil {
		if debug {
			g.Debugf("found instantiated template type %v", maker)
		}
	} else {
		if debug {
			g.Debugf("instantiating template type %v", maker)
		}
		// hard part: instantiate the template type.
		// must be instantiated in the same *Comp where it was declared!
		instance = maker.instantiateType(typ, node)
	}
	return instance
}

// instantiateTemplateType instantiates and compiles a template function.
// node is used only for error messages
func (maker *templateMaker) instantiateType(typ *TemplateType, node *ast.IndexExpr) xr.Type {

	// choose the specialization to use
	_, special := maker.chooseType(typ)

	// create a new nested Comp
	c := NewComp(maker.comp, nil)
	c.UpCost = 0
	c.Depth--

	// and inject template arguments in it
	special.injectBinds(c)

	key := maker.ikey
	panicking := true
	defer func() {
		if panicking {
			delete(typ.Instances, key) // remove the cached instance if present
			c.ErrorAt(node.Pos(), "error instantiating template type: %v\n\t%v", maker, recover())
		}
	}()
	// compile the type instantiation
	//
	var t xr.Type
	if !special.decl.Alias && maker.sym.Name != "_" {
		if c.Globals.Options&base.OptDebugTemplate != 0 {
			c.Debugf("forward-declaring template type before instantiation: %v", maker)
		}
		// support for template recursive types, as for example
		//   template[T] type List struct { First T; Rest *List#[T] }
		// requires to cache List#[T] as instantiated **before** actually instantiating it.
		//
		// This is similar to the technique used for non-template recursive types, as
		//    type List struct { First int; Rest *List }
		// with the difference that the cache is typ.Instances[key] instead of Comp.Types[name]
		t = c.Universe.NamedOf(maker.String(), c.FileComp().Path, r.Invalid /*kind not yet known*/)
		typ.Instances[key] = t
		u := c.Type(special.decl.Decl)
		c.SetUnderlyingType(t, u)
	} else {
		// either the template type is an alias, or name == "_" (discards the result of type declaration)
		t = c.Type(special.decl.Decl)
		typ.Instances[key] = t
	}
	panicking = false
	return t
}
