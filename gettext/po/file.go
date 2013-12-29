// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// File represents an PO File.
//
// See http://www.gnu.org/software/gettext/manual/html_node/
type File struct {
	MimeHeader Header
	Messages   []Message
}

func Load(name string) (*File, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadData(data)
}

func LoadData(data []byte) (*File, error) {
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
	return &file, nil
}

func (f *File) Save(name string) error {
	return ioutil.WriteFile(name, []byte(f.String()), 0666)
}

func (f *File) Data(name string) []byte {
	return nil
}

func (f *File) PGettext(msgctxt, msgid string) string {
	return msgid
}

func (f *File) PNGettext(msgctxt, msgid, msgidPlural string, n int) string {
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
