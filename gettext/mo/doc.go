// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mo implements a GNU MO file decoder and encoder.

Examples:
	import (
		"code.google.com/p/gettext-go/gettext/mo"
	)

	func main() {
		moFile, err := mo.Load("../testdata/test.mo", nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", moFile)

		msg := moFile.PGettext("", "Hello gettext-go")
		fmt.Printf("msg: %s\n", msg)
	}

The GNU MO file specification is at
http://www.gnu.org/software/gettext/manual/html_node/MO-Files.html.
*/
package mo
