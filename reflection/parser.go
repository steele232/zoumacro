// this file was generated by gomacro command: import "github.com/cosmos72/gomacro/parser"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package reflection

import (
	. "reflect"

	. "github.com/cosmos72/gomacro/imports"
	"github.com/cosmos72/gomacro/parser"
)

// reflection: allow interpreted code to import "github.com/cosmos72/gomacro/parser"
func init() {
	Packages["github.com/cosmos72/gomacro/parser"] = Package{
		Binds: map[string]Value{
			"AllErrors":         ValueOf(parser.AllErrors),
			"DeclarationErrors": ValueOf(parser.DeclarationErrors),
			"ImportsOnly":       ValueOf(parser.ImportsOnly),
			"PackageClauseOnly": ValueOf(parser.PackageClauseOnly),
			"ParseComments":     ValueOf(parser.ParseComments),
			"SpuriousErrors":    ValueOf(parser.SpuriousErrors),
			"Trace":             ValueOf(parser.Trace),
		},
		Types: map[string]Type{
			"Bailout": TypeOf((*parser.Bailout)(nil)).Elem(),
			"Mode":    TypeOf((*parser.Mode)(nil)).Elem(),
			"Parser":  TypeOf((*parser.Parser)(nil)).Elem(),
		},
		Proxies: map[string]Type{}}
}