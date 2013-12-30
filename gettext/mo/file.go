// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"

	"code.google.com/p/gettext-go/gettext/plural"
)

const (
	moHeaderSize        = 28
	moMagicLittleEndian = 0x950412de
	moMagicBigEndian    = 0xde120495

	EotSeparator = "\x04" // msgctxt and msgid separator
	NulSeparator = "\x00" // msgid and msgstr separator
)

// File represents an MO File.
//
// See http://www.gnu.org/software/gettext/manual/html_node/MO-Files.html
type File struct {
	MajorVersion  uint16
	MinorVersion  uint16
	MimeHeader    map[string]string
	MessageMap    map[string]Message
	PluralFormula func(n int) int
}

// A PO file is made up of many entries,
// each entry holding the relation between an original untranslated string
// and its corresponding translation.
//
// See http://www.gnu.org/software/gettext/manual/html_node/PO-Files.html
type Message struct {
	MsgContext   string   // msgctxt context
	MsgId        string   // msgid untranslated-string
	MsgIdPlural  string   // msgid_plural untranslated-string-plural
	MsgStr       string   // msgstr translated-string
	MsgStrPlural []string // msgstr[0] translated-string-case-0
}

// MakeMessageMapKey returns the File.MessageMap key string.
func MakeMessageMapKey(msgctxt, msgid string) string {
	if msgctxt != "" {
		return msgctxt + EotSeparator + msgid
	}
	return msgid
}

// Load loads a named mo file.
func Load(name string, pluralFormula func(n int) int) (*File, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadData(data, pluralFormula)
}

// LoadData loads mo file format data.
func LoadData(data []byte, pluralFormula func(n int) int) (*File, error) {
	r := bytes.NewReader(data)

	var magicNumber uint32
	if err := binary.Read(r, binary.LittleEndian, &magicNumber); err != nil {
		return nil, fmt.Errorf("gettext: %v", err)
	}
	var bo binary.ByteOrder
	switch magicNumber {
	case moMagicLittleEndian:
		bo = binary.LittleEndian
	case moMagicBigEndian:
		bo = binary.BigEndian
	default:
		return nil, fmt.Errorf("gettext: %v", "invalid magic number")
	}

	var header struct {
		MajorVersion uint16
		MinorVersion uint16
		MsgIdCount   uint32
		MsgIdOffset  uint32
		MsgStrOffset uint32
		HashSize     uint32
		HashOffset   uint32
	}
	if err := binary.Read(r, bo, &header); err != nil {
		return nil, fmt.Errorf("gettext: %v", err)
	}
	if v := header.MajorVersion; v != 0 && v != 1 {
		return nil, fmt.Errorf("gettext: %v", "invalid version number")
	}
	if v := header.MinorVersion; v != 0 && v != 1 {
		return nil, fmt.Errorf("gettext: %v", "invalid version number")
	}

	msgIdStart := make([]uint32, header.MsgIdCount)
	msgIdLen := make([]uint32, header.MsgIdCount)
	if _, err := r.Seek(int64(header.MsgIdOffset), 0); err != nil {
		return nil, fmt.Errorf("gettext: %v", err)
	}
	for i := 0; i < int(header.MsgIdCount); i++ {
		if err := binary.Read(r, bo, &msgIdLen[i]); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}
		if err := binary.Read(r, bo, &msgIdStart[i]); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}
	}

	msgStrStart := make([]int32, header.MsgIdCount)
	msgStrLen := make([]int32, header.MsgIdCount)
	if _, err := r.Seek(int64(header.MsgStrOffset), 0); err != nil {
		return nil, fmt.Errorf("gettext: %v", err)
	}
	for i := 0; i < int(header.MsgIdCount); i++ {
		if err := binary.Read(r, bo, &msgStrLen[i]); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}
		if err := binary.Read(r, bo, &msgStrStart[i]); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}
	}

	file := &File{
		MajorVersion:  header.MajorVersion,
		MinorVersion:  header.MinorVersion,
		MimeHeader:    make(map[string]string),
		MessageMap:    make(map[string]Message),
		PluralFormula: pluralFormula,
	}
	for i := 0; i < int(header.MsgIdCount); i++ {
		if _, err := r.Seek(int64(msgIdStart[i]), 0); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}
		msgIdData := make([]byte, msgIdLen[i])
		if _, err := r.Read(msgIdData); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}

		if _, err := r.Seek(int64(msgStrStart[i]), 0); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}
		msgStrData := make([]byte, msgStrLen[i])
		if _, err := r.Read(msgStrData); err != nil {
			return nil, fmt.Errorf("gettext: %v", err)
		}

		if len(msgIdData) == 0 {
			ss := strings.Split(string(msgStrData), "\n")
			for i := 0; i < len(ss); i++ {
				idx := strings.Index(ss[i], ":")
				if idx < 0 {
					continue
				}
				key := strings.TrimSpace(ss[i][:idx])
				val := strings.TrimSpace(ss[i][idx+1:])
				file.MimeHeader[key] = val
			}
		} else {
			var msg = Message{
				MsgId:  string(msgIdData),
				MsgStr: string(msgStrData),
			}
			// Is this a context message?
			if idx := strings.Index(msg.MsgId, EotSeparator); idx != -1 {
				msg.MsgContext, msg.MsgId = msg.MsgId[:idx], msg.MsgId[idx+1:]
			}
			// Is this a plural message?
			if idx := strings.Index(msg.MsgId, NulSeparator); idx != -1 {
				msg.MsgId, msg.MsgIdPlural = msg.MsgId[:idx], msg.MsgId[idx+1:]
				msg.MsgStrPlural = strings.Split(msg.MsgStr, NulSeparator)
				msg.MsgStr = ""
			}
			file.MessageMap[MakeMessageMapKey(msg.MsgContext, msg.MsgId)] = msg
		}
	}
	if file.PluralFormula == nil {
		if lang := file.MimeHeader["Language"]; lang != "" {
			file.PluralFormula = plural.Formula(lang)
		} else {
			file.PluralFormula = plural.Formula("??")
		}
	}

	return file, nil
}

// Save saves a mo file.
func (f *File) Save(name string) error {
	return ioutil.WriteFile(name, f.Data(name), 0666)
}

// Save returns a mo file format data.
func (f *File) Data(name string) []byte {
	return encodeFile(f)
}

// PGettext attempt to translate a text string,
// by looking up the translation in current mo file.
//
// Examples:
//	func Foo() {
//		msg := moFile.PGettext("gettext-go.example", "Hello") // msgctxt is "gettext-go.example"
//	}
func (f *File) PGettext(msgctxt, msgid string) string {
	return f.PNGettext(msgctxt, msgid, "", 0)
}

// PNGettext attempt to translate a text string,
// by looking up the translation in current mo file.
// catalog.
//
// Examples:
//	func Foo() {
//		msg := moFile.PNGettext("gettext-go.example", "%d people", "%d peoples", 2)
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
	key := MakeMessageMapKey(msgctxt, msgid)
	if v, ok := f.MessageMap[key]; ok {
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
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# version: %d.%d\n", f.MajorVersion, f.MinorVersion)
	fmt.Fprintf(&buf, `msgid ""`+"\n")
	fmt.Fprintf(&buf, `msgstr ""`+"\n")
	for k, v := range f.MimeHeader {
		fmt.Fprintf(&buf, `"%s: %s\n"`+"\n", k, v)
	}
	fmt.Fprintf(&buf, "\n")

	for k, v := range f.MessageMap {
		fmt.Fprintf(&buf, `msgid "%s"`+"\n", k)
		fmt.Fprintf(&buf, `msgstr "%s"`+"\n", v.MsgStr)
		fmt.Fprintf(&buf, "\n")
	}

	return buf.String()
}
