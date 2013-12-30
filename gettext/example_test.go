// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext_test

import (
	"fmt"

	"code.google.com/p/gettext-go/gettext"
)

func _Example() {
	gettext.SetLocale("zh_CN")
	gettext.BindTextdomain("hello", "../examples/local")
	gettext.Textdomain("hello")

	fmt.Println(gettext.PGettext("main.main", "Hello, world!"))
	// Output: 你好, 世界!
}
