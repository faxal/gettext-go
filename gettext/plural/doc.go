// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package plural provides standard plural formulas.

Examples:
	import (
		"code.google.com/p/gettext-go/gettext/mo"
		"code.google.com/p/gettext-go/gettext/plural"
	)

	func main() {
		enFormula := plural.Formula("en_US")
		moFile, err := mo.Load("../testdata/test.mo", enFormula)
		if err != nil {
			log.Fatal(err)
		}
		...
	}

See http://www.gnu.org/software/gettext/manual/html_node/Plural-forms.html
*/
package plural
