// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type moHeader struct {
	MagicNumber  uint32
	MajorVersion uint16
	MinorVersion uint16
	MsgIdCount   uint32
	MsgIdOffset  uint32
	MsgStrOffset uint32
	HashSize     uint32
	HashOffset   uint32
}

type moStrPos struct {
	Addr uint32
	Size uint32
}

func encodeFile(f *File) []byte {
	hdr := &moHeader{
		MagicNumber: moMagicLittleEndian,
	}
	data := encodeData(hdr, f)
	data = append(encodeHeader(hdr), data...)
	return data
}

// encode data and init moHeader
func encodeData(hdr *moHeader, f *File) []byte {
	var msgList = []Message{mimeHeaderToMessage(f.MimeHeader)}
	for _, v := range f.MessageMap {
		msgList = append(msgList, v)
	}
	var buf bytes.Buffer
	var msgIdPosList = make([]moStrPos, len(msgList))
	var msgStrPosList = make([]moStrPos, len(msgList))
	for i, v := range msgList {
		// write msgid
		msgId := encodeMsgId(v)
		msgIdPosList[i].Addr = uint32(buf.Len() + moHeaderSize)
		msgIdPosList[i].Size = uint32(len(msgId))
		buf.WriteString(msgId)
		// write msgstr
		msgStr := encodeMsgStr(v)
		msgStrPosList[i].Addr = uint32(buf.Len() + moHeaderSize)
		msgStrPosList[i].Size = uint32(len(msgStr))
		buf.WriteString(msgStr)
	}

	hdr.MsgIdOffset = uint32(buf.Len() + moHeaderSize)
	binary.Write(&buf, binary.LittleEndian, msgIdPosList)
	hdr.MsgStrOffset = uint32(buf.Len() + moHeaderSize)
	binary.Write(&buf, binary.LittleEndian, msgStrPosList)

	hdr.MsgIdCount = uint32(len(msgList))
	return buf.Bytes()
}

// must called after encodeData
func encodeHeader(hdr *moHeader) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, hdr)
	return buf.Bytes()
}

func encodeMsgId(v Message) string {
	if v.MsgContext != "" && v.MsgIdPlural != "" {
		return v.MsgContext + EotSeparator + v.MsgId + NulSeparator + v.MsgIdPlural
	}
	if v.MsgContext != "" && v.MsgIdPlural == "" {
		return v.MsgContext + EotSeparator + v.MsgId
	}
	if v.MsgContext == "" && v.MsgIdPlural != "" {
		return v.MsgId + NulSeparator + v.MsgIdPlural
	}
	return v.MsgId
}

func encodeMsgStr(v Message) string {
	if v.MsgIdPlural != "" {
		return strings.Join(v.MsgStrPlural, NulSeparator)
	}
	return v.MsgStr
}

func mimeHeaderToMessage(mimeHeader map[string]string) Message {
	var buf bytes.Buffer
	for k, v := range mimeHeader {
		fmt.Fprintf(&buf, "%s: %s\n", k, v)
	}
	return Message{
		MsgStr: buf.String(),
	}
}

type byMessages []Message

func (d byMessages) Len() int {
	return len(d)
}
func (d byMessages) Less(i, j int) bool {
	if a, b := d[i].MsgContext, d[j].MsgContext; a != b {
		return a < b
	}
	if a, b := d[i].MsgId, d[j].MsgId; a != b {
		return a < b
	}
	if a, b := d[i].MsgIdPlural, d[j].MsgIdPlural; a != b {
		return a < b
	}
	return false
}
func (d byMessages) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}