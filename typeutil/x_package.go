// this file was generated by gomacro command: import _i "github.com/steele232/zoumacro/typeutil"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package typeutil

import (
	r "reflect"

	"github.com/steele232/zoumacro/imports"
)

// reflection: allow interpreted code to import "github.com/steele232/zoumacro/typeutil"
func init() {
	imports.Packages["github.com/steele232/zoumacro/typeutil"] = imports.Package{
		Binds: map[string]r.Value{
			"Identical":           r.ValueOf(Identical),
			"IdenticalIgnoreTags": r.ValueOf(IdenticalIgnoreTags),
			"MakeHasher":          r.ValueOf(MakeHasher),
		},
		Types: map[string]r.Type{
			"Hasher": r.TypeOf((*Hasher)(nil)).Elem(),
			"Map":    r.TypeOf((*Map)(nil)).Elem(),
		},
		Proxies: map[string]r.Type{}}
}
