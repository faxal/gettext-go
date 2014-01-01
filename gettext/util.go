// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"regexp"
	"runtime"
)

var (
	reInit    = regexp.MustCompile(`init·\d+$`) // main.init·1
	reClosure = regexp.MustCompile(`func·\d+$`) // main.func·001
)

func callerName(calldepth int) string {
	// caller types:
	// runtime.goexit
	// runtime.main
	// main.init
	// main.init·1
	// main.main
	// main.func·001
	// code.google.com/p/gettext-go/gettext.TestCallerName
	// ...
	for {
		pc, _, _, ok := runtime.Caller(calldepth)
		if !ok {
			return ""
		}
		name := runtime.FuncForPC(pc).Name()
		if reInit.MatchString(name) {
			return reInit.ReplaceAllString(name, "init")
		}
		if reClosure.MatchString(name) {
			calldepth++
			continue
		}
		return name
	}
}
