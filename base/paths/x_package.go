// this file was generated by gomacro command: import _i "github.com/steele232/zoumacro/base/paths"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package paths

import (
	r "reflect"
	"github.com/steele232/zoumacro/imports"
)

// reflection: allow interpreted code to import "github.com/steele232/zoumacro/base/paths"
func init() {
	imports.Packages["github.com/steele232/zoumacro/base/paths"] = imports.Package{
	Binds: map[string]r.Value{
		"DirName":	r.ValueOf(DirName),
		"FileName":	r.ValueOf(FileName),
		"GoSrcDir":	r.ValueOf(&GoSrcDir).Elem(),
		"GomacroDir":	r.ValueOf(&GomacroDir).Elem(),
		"RemoveLastByte":	r.ValueOf(RemoveLastByte),
		"Subdir":	r.ValueOf(Subdir),
		"UserHomeDir":	r.ValueOf(UserHomeDir),
	}, 
	}
}
