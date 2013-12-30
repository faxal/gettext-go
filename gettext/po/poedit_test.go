// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package po

import (
	"reflect"
	"testing"
)

var (
	testPoEditPoFile = "../testdata/poedit-1.5.7-zh_CN.po"
	testPoEditMoFile = "../testdata/poedit-1.5.7-zh_CN.mo"
)

func _TestPoEditPoFile(t *testing.T) {
	po, err := Load(testPoEditPoFile, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&po.MimeHeader, &poEditFile.MimeHeader) {
		t.Fatalf("expect = %v, got = %v", &poEditFile.MimeHeader, &po.MimeHeader)
	}
	if len(po.MessageMap) != len(poEditFile.MessageMap) {
		t.Fatal("size not equal")
	}
	for k, v0 := range po.MessageMap {
		v1, ok := poEditFile.MessageMap[k]
		if !ok {
			t.Fatalf("key %q not exists", k)
		}
		if !reflect.DeepEqual(&v0, &v1) {
			t.Fatalf("%s: expect = %v, got = %v", k, v1, v0)
		}
	}
}

var poEditFile = &File{}
