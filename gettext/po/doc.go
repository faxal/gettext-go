// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package po provides support for reading and writing GNU PO file.

Examples:
	import (
		"code.google.com/p/gettext-go/gettext/po"
	)

	func main() {
		poFile, err := po.Load("../testdata/test.po", nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", poFile)

		msg := poFile.PGettext("", "Hello gettext-go")
		fmt.Printf("msg: %s\n", msg)
	}

The GNU PO file specification is at
http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html.
*/
package po
