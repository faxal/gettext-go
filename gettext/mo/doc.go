// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mo provides support for reading and writing GNU MO file.

Examples:
	import (
		"code.google.com/p/gettext-go/gettext/mo"
	)

	func main() {
		moFile, err := mo.Load("test.mo")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", moFile)
	}

The GNU MO file specification is at
http://www.gnu.org/software/gettext/manual/html_node/MO-Files.html.
*/
package mo
