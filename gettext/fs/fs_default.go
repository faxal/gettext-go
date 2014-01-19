// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"os"
)

// DefaultFileSystem is a FileSystem implementation backed by the underlying
// operating system's file system.
var DefaultFileSystem FileSystem = defFS{}

type defFS struct{}

func (defFS) Create(name string) (File, error) {
	return os.Create(name)
}

func (defFS) Open(name string) (File, error) {
	return os.Open(name)
}

func (defFS) Remove(name string) error {
	return os.Remove(name)
}

func (defFS) MkdirAll(dir string, perm os.FileMode) error {
	return os.MkdirAll(dir, perm)
}

func (defFS) List(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdirnames(-1)
}

func (defFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
