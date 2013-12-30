// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po_test

import (
	"fmt"
	"log"

	"code.google.com/p/gettext-go/gettext/po"
)

func _ExampleFile() {
	testfile := "../../testdata/test.po"
	f, err := po.Load(testfile, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %s\n", "Title", f.PGettext("", "Title"))
	fmt.Printf("%s: %s\n", "%d topic", f.PGettext("", "%d topic"))
	fmt.Printf("%s: %s\n", "%d topic", f.PNGettext("", "%d topic", "%d topics", 0))
	fmt.Printf("%s: %s\n", "%d topics", f.PNGettext("", "%d topic", "%d topics", 1))
	fmt.Printf("%s: %s\n", "%d topics", f.PNGettext("", "%d topic", "%d topics", 2))
	// Output:
	// Title: Título
	// %d topic: %d tema
	// %d topic: %d tema
	// %d topics: %d temas
	// %d topics: %d temas
}
