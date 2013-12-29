// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"

	"code.google.com/p/gettext-go/gettext"
	"code.google.com/p/gettext-go/gettext/po"
)

func main() {
	poFile, err := po.Load("../testdata/test.po")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", poFile)

	msg := gettext.Gettext("Hello gettext-go")
	fmt.Printf("msg: %s\n", msg)
}
