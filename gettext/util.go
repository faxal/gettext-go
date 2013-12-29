// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"runtime"
)

func callerName(calldepth int) string {
	pc, _, _, _ := runtime.Caller(calldepth)
	return runtime.FuncForPC(pc).Name()
}
