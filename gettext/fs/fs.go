// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"io"
	"os"
	"strings"
)

// File is a readable, writable sequence of bytes.
//
// Typically, it will be an *os.File, but test code may choose to substitute
// memory-backed implementations.
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Writer
	Stat() (os.FileInfo, error)
	Sync() error
}

// FileSystem is a namespace for files.
//
// The names are filepath names: they may be / separated or \ separated,
// depending on the underlying operating system.
type FileSystem interface {
	// Create creates the named file for writing, truncating it if it already
	// exists.
	Create(name string) (File, error)

	// Open opens the named file for reading.
	Open(name string) (File, error)

	// Remove removes the named file or directory.
	Remove(name string) error

	// MkdirAll creates a directory and all necessary parents. The permission
	// bits perm have the same semantics as in os.MkdirAll. If the directory
	// already exists, MkdirAll does nothing and returns nil.
	MkdirAll(dir string, perm os.FileMode) error

	// List returns a listing of the given directory. The names returned are
	// relative to dir.
	List(dir string) ([]string, error)

	// Stat returns an os.FileInfo describing the named file.
	Stat(name string) (os.FileInfo, error)
}

var fsDrivers = make(map[string]FileSystem)

// Register makes a file system driver available by the provided name.
// If Register is called twice with the same name or if driver is nil, it panics.
// The emptry prefix is reserved for the default file system.
func Register(prefix string, driver FileSystem) {
	if driver == nil {
		panic("gettext/fs: Register driver is nil")
	}
	if _, dup := fsDrivers[prefix]; dup {
		panic("gettext/fs: Register called twice for driver " + prefix)
	}
	fsDrivers[prefix] = driver
}

// Create creates the named file for writing, truncating it if it already
// exists.
func Create(name string) (File, error) {
	return getFsDriver(name).Create(name)
}

// Open opens the named file for reading.
func Open(name string) (File, error) {
	return getFsDriver(name).Open(name)
}

// Remove removes the named file or directory.
func Remove(name string) error {
	return getFsDriver(name).Remove(name)
}

// MkdirAll creates a directory and all necessary parents. The permission
// bits perm have the same semantics as in os.MkdirAll. If the directory
// already exists, MkdirAll does nothing and returns nil.
func MkdirAll(dir string, perm os.FileMode) error {
	return getFsDriver(dir).MkdirAll(dir, perm)
}

// List returns a listing of the given directory. The names returned are
// relative to dir.
func List(dir string) ([]string, error) {
	return getFsDriver(dir).List(dir)
}

// Stat returns an os.FileInfo describing the named file.
func Stat(name string) (os.FileInfo, error) {
	return getFsDriver(name).Stat(name)
}

func getFsDriver(path string) FileSystem {
	for prefix, driver := range fsDrivers {
		if strings.HasPrefix(path, prefix) {
			return driver
		}
	}
	return DefaultFileSystem
}
