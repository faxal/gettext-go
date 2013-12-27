// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"reflect"
	"testing"
)

var (
	testPoEditPoFile = "../testdata/poedit-1.5.7-zh_CN.po"
	testPoEditMoFile = "../testdata/poedit-1.5.7-zh_CN.mo"
)

func _TestPoEditPoFile(t *testing.T) {
	po, err := LoadPoFile(testPoEditPoFile)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&po.Header, &poEditFile.Header) {
		t.Fatalf("expect = %v, got = %v", &poEditFile.Header, &po.Header)
	}
	for i := 0; i < len(po.Entrys) && i < len(poEditFile.Entrys); i++ {
		if !reflect.DeepEqual(&po.Entrys[i], &poEditFile.Entrys[i]) {
			t.Fatalf("%d: expect = %v, got = %v", i, poEditFile.Entrys[i], po.Entrys[i])
		}
	}
}

var poEditFile = &PoFile{}
