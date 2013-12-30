// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sort"

	"code.google.com/p/gettext-go/gettext/plural"
)

// File represents an PO File.
//
// See http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
type File struct {
	MimeHeader    Header
	MessageMap    map[string]Message
	PluralFormula func(n int) int
}

// MakeMessageMapKey returns the File.MessageMap key string.
func MakeMessageMapKey(msgctxt, msgid string) string {
	if msgctxt != "" {
		const eotSeparator = "\x04"
		return msgctxt + eotSeparator + msgid
	}
	return msgid
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
	r := newLineReader(string(data))
	var file = File{
		MessageMap: make(map[string]Message),
	}
	for {
		var msg Message
		if err := msg.readPoEntry(r); err != nil {
			if err == io.EOF {
				return &file, nil
			}
			return nil, err
		}
		if msg.MsgId == "" {
			file.MimeHeader.parseHeader(&msg)
			continue
		}
		file.MessageMap[MakeMessageMapKey(msg.MsgContext, msg.MsgId)] = msg
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
func (f *File) Data() []byte {
	// sort the massge as ReferenceFile/ReferenceLine field
	var messages []Message
	for _, v := range f.MessageMap {
		messages = append(messages, v)
	}
	sort.Sort(byMessages(messages))

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", f.MimeHeader.String())
	for i := 0; i < len(messages); i++ {
		fmt.Fprintf(&buf, "%s\n", messages[i].String())
	}
	return buf.Bytes()
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
	n = f.PluralFormula(n)
	if ss := f.findMsgStrPlural(msgctxt, msgid, msgidPlural); len(ss) != 0 {
		if n >= len(ss) {
			n = len(ss) - 1
		}
		return ss[n]
	}
	if msgidPlural != "" && n > 0 {
		return msgidPlural
	} else {
		return msgid
	}
}

func (f *File) findMsgStrPlural(msgctxt, msgid, msgidPlural string) []string {
	if v, ok := f.MessageMap[MakeMessageMapKey(msgctxt, msgid)]; ok {
		if len(v.MsgIdPlural) != 0 {
			if len(v.MsgStrPlural) != 0 {
				return v.MsgStrPlural
			} else {
				return nil
			}
		} else {
			if len(v.MsgStr) != 0 {
				return []string{v.MsgStr}
			} else {
				return nil
			}
		}
	}
	return nil
}

// String returns the po format file string.
func (f *File) String() string {
	return string(f.Data())
}
