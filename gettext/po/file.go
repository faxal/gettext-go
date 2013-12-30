// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"code.google.com/p/gettext-go/gettext/plural"
)

// File represents an PO File.
//
// See http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
type File struct {
	MimeHeader    Header
	Messages      []Message // discard
	MessageMap    map[string]Message
	PluralFormula func(n int) int
}

// Load loads a named po file.
func Load(name string, pluralFormula func(n int) int) (*File, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadData(data, pluralFormula)
}

// LoadData loads po file format data.
func LoadData(data []byte, pluralFormula func(n int) int) (*File, error) {
	var file File
	r := newLineReader(string(data))
	for {
		var entry Message
		if err := entry.readPoEntry(r); err != nil {
			if err == io.EOF {
				return &file, nil
			}
			return nil, err
		}
		if entry.MsgId == "" {
			file.MimeHeader.parseHeader(&entry)
			continue
		}
		file.Messages = append(file.Messages, entry)
	}
	file.PluralFormula = pluralFormula
	if file.PluralFormula == nil {
		if lang := file.MimeHeader.Language; lang != "" {
			file.PluralFormula = plural.Formula(lang)
		} else {
			file.PluralFormula = plural.Formula("??")
		}
	}

	return &file, nil
}

// Save saves a po file.
func (f *File) Save(name string) error {
	return ioutil.WriteFile(name, []byte(f.String()), 0666)
}

// Save returns a po file format data.
func (f *File) Data(name string) []byte {
	// sort the massge as ReferenceFile/ReferenceLine field
	panic("TODO")
}

// PGettext attempt to translate a text string,
// by looking up the translation in current po file.
//
// Examples:
//	func Foo() {
//		msg := poFile.PGettext("gettext-go.example", "Hello") // msgctxt is "gettext-go.example"
//	}
func (f *File) PGettext(msgctxt, msgid string) string {
	return f.PNGettext(msgctxt, msgid, "", 0)
}

// PNGettext attempt to translate a text string,
// by looking up the translation in current po file.
// catalog.
//
// Examples:
//	func Foo() {
//		msg := poFile.PNGettext("gettext-go.example", "%d people", "%d peoples", 2)
//	}
func (f *File) PNGettext(msgctxt, msgid, msgidPlural string, n int) string {
	panic("TODO")
}

// String returns the po format file string.
func (f *File) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", f.MimeHeader.String())
	for i := 0; i < len(f.Messages); i++ {
		fmt.Fprintf(&buf, "%s\n", f.Messages[i].String())
	}
	return buf.String()
}
