// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gettext  implements a basic GNU's gettext library.

Example:
	import (
		"code.google.com/p/gettext-go/gettext"
	)

	func main() {
		gettext.SetLocale("zh_CN")
		gettext.BindTextdomain("hello", "local")
		gettext.Textdomain("hello")

		fmt.Println(gettext.Gettext("Hello, world!"))
		// Output: 你好, 世界!
	}

See:
	http://en.wikipedia.org/wiki/Gettext
	http://www.gnu.org/software/gettext/manual/html_node
	http://www.gnu.org/software/gettext/manual/html_node/Header-Entry.html
	http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
	http://www.gnu.org/software/gettext/manual/html_node/MO-Files.html
	http://www.poedit.net/
*/
package gettext
