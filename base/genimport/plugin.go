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
 * plugin.go
 *
 *  Created on Feb 27, 2017
 *      Author Massimiliano Ghilardi
 */

package genimport

import (
	"io"
	"os/exec"
	r "reflect"

	"github.com/steele232/zoumacro/base/paths"
)

func compilePlugin(o *Output, filepath string, stdout io.Writer, stderr io.Writer) string {
	gosrcdir := paths.GoSrcDir
	gosrclen := len(gosrcdir)
	filelen := len(filepath)
	if filelen < gosrclen || filepath[0:gosrclen] != gosrcdir {
		o.Errorf("source %q is in unsupported directory, cannot compile it: should be inside %q", filepath, gosrcdir)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin")
	cmd.Dir = paths.DirName(filepath)
	cmd.Stdin = nil
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	o.Debugf("compiling %q ...", filepath)
	err := cmd.Run()
	if err != nil {
		o.Errorf("error executing \"go build -buildmode=plugin\" in directory %q: %v", cmd.Dir, err)
	}

	dirname := paths.RemoveLastByte(paths.DirName(filepath))
	// go build uses innermost directory name as shared object name,
	// i.e.	foo/bar/main.go is compiled to foo/bar/bar.so
	filename := paths.FileName(dirname)

	return paths.Subdir(dirname, filename+".so")
}

func (imp *Importer) loadPluginSymbol(soname string, symbolName string) interface{} {
	// use imports.Packages["plugin"].Binds["Open"] and reflection instead of hard-coding call to plugin.Open()
	// reasons:
	// * import ( "plugin" ) does not work on all platforms (creates broken gomacro.exe on Windows/386)
	// * allow caller to provide us with a different implementation,
	//   either in imports.Packages["plugin"].Binds["Open"]
	//   or in Globals.Importer.PluginOpen

	o := imp.output
	if !imp.setPluginOpen() {
		o.Errorf("gomacro compiled without support to load plugins - requires Go 1.8+ and Linux - cannot import packages at runtime")
	}
	if len(soname) == 0 || len(symbolName) == 0 {
		// caller is just checking whether PluginOpen() is available
		return nil
	}
	so, err := reflectcall(imp.PluginOpen, soname)
	if err != nil {
		o.Errorf("error loading plugin %q: %v", soname, err)
	}
	vsym, err := reflectcall(so.MethodByName("Lookup"), symbolName)
	if err != nil {
		o.Errorf("error loading symbol %q from plugin %q: %v", symbolName, soname, err)
	}
	return vsym.Interface()
}

func reflectcall(fun r.Value, arg interface{}) (r.Value, interface{}) {
	vs := fun.Call([]r.Value{r.ValueOf(arg)})
	return vs[0], vs[1].Interface()
}
