// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// PoFile represents an PO PoFile.
//
// See http://www.gnu.org/software/gettext/manual/html_node/
type PoFile struct {
	Header PoHeader
	Entrys []PoEntry
}

func NewPoFile() *PoFile {
	return &PoFile{}
}

func LoadPoFile(name string) (*PoFile, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadPoData(data)
}

func LoadPoData(data []byte) (*PoFile, error) {
	var file PoFile
	r := newLineReader(string(data))
	for {
		var entry PoEntry
		if err := entry.readPoEntry(r); err != nil {
			if err == io.EOF {
				return &file, nil
			}
			return nil, err
		}
		if entry.MsgId == "" {
			file.Header.parseHeader(&entry)
			continue
		}
		file.Entrys = append(file.Entrys, entry)
	}
	return &file, nil
}

func (f *PoFile) SaveFile(name string) error {
	return ioutil.WriteFile(name, []byte(f.String()), 0666)
}

func (f *PoFile) ToMoFile() *MoFile {
	return nil
}

func (f *PoFile) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", f.Header.String())
	for i := 0; i < len(f.Entrys); i++ {
		fmt.Fprintf(&buf, "%s\n", f.Entrys[i].String())
	}
	return buf.String()
}
