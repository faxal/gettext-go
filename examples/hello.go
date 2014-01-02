// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"code.google.com/p/gettext-go/gettext"
)

func init() {
	// bind app domain
	gettext.BindTextdomain("hello", "local")
	gettext.Textdomain("hello")

	// $(LC_MESSAGES) or $(LANG) or empty
	fmt.Println(gettext.Gettext("Gettext in init."))
	fmt.Println(gettext.PGettext("main.init", "Gettext in init."))
	// Output(depends on local environment):
	// ?
	// ?

	// set simple chinese
	gettext.SetLocale("zh_CN")

	// simple chinese
	fmt.Println(gettext.Gettext("Gettext in init."))
	fmt.Println(gettext.PGettext("main.init", "Gettext in init."))
	// Output:
	// Init函数中的Gettext.
	// Init函数中的Gettext.
}

func main() {
	// simple chinese
	fmt.Println(gettext.Gettext("Hello, world!"))
	fmt.Println(gettext.PGettext("main.main", "Hello, world!"))
	// Output:
	// 你好, 世界!
	// 你好, 世界!

	// set traditional chinese
	gettext.SetLocale("zh_TW")

	// traditional chinese
	func() {
		fmt.Println(gettext.Gettext("Gettext in func."))
		fmt.Println(gettext.PGettext("main.func", "Gettext in func."))
		// Output:
		// 閉包函數中的Gettext.
		// 閉包函數中的Gettext.
	}()
}
