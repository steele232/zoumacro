// this file was generated by gomacro command: import _i "github.com/steele232/zoumacro/base/untyped"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package untyped

import (
	r "reflect"
	"github.com/steele232/zoumacro/imports"
)

// reflection: allow interpreted code to import "github.com/steele232/zoumacro/base/untyped"
func init() {
	imports.Packages["github.com/steele232/zoumacro/base/untyped"] = imports.Package{
	Binds: map[string]r.Value{
		"Bool":	r.ValueOf(Bool),
		"Complex":	r.ValueOf(Complex),
		"ConvertLiteralCheckOverflow":	r.ValueOf(ConvertLiteralCheckOverflow),
		"Float":	r.ValueOf(Float),
		"GoUntypedToKind":	r.ValueOf(GoUntypedToKind),
		"Int":	r.ValueOf(Int),
		"MakeKind":	r.ValueOf(MakeKind),
		"MakeLit":	r.ValueOf(MakeLit),
		"Marshal":	r.ValueOf(Marshal),
		"None":	r.ValueOf(None),
		"Rune":	r.ValueOf(Rune),
		"String":	r.ValueOf(String),
		"Unmarshal":	r.ValueOf(Unmarshal),
		"UnmarshalVal":	r.ValueOf(UnmarshalVal),
	}, Types: map[string]r.Type{
		"Kind":	r.TypeOf((*Kind)(nil)).Elem(),
		"Lit":	r.TypeOf((*Lit)(nil)).Elem(),
		"Val":	r.TypeOf((*Val)(nil)).Elem(),
	}, 
	}
}
