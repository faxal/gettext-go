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
// See http://www.gnu.org/software/gettext/manual/html_node/
type File struct {
	MimeHeader    Header
	Messages      []Message // discard
	MessageMap    map[string]Message
	PluralFormula func(n int) int
}

func Load(name string, pluralFormula func(n int) int) (*File, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadData(data, pluralFormula)
}

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

func (f *File) Save(name string) error {
	return ioutil.WriteFile(name, []byte(f.String()), 0666)
}

func (f *File) Data(name string) []byte {
	return nil
}

func (f *File) PGettext(msgctxt, msgid string) string {
	return f.PNGettext(msgctxt, msgid, "", 0)
}

func (f *File) PNGettext(msgctxt, msgid, msgidPlural string, n int) string {
	n = f.PluralFormula(n)
	for i := 0; i < len(f.Messages); i++ {
		//
	}
	return msgid
}

func (f *File) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", f.MimeHeader.String())
	for i := 0; i < len(f.Messages); i++ {
		fmt.Fprintf(&buf, "%s\n", f.Messages[i].String())
	}
	return buf.String()
}
