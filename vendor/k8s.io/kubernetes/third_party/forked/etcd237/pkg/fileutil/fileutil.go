// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fileutil implements utility functions related to files and paths.
package fileutil

import (
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/coreos/pkg/capnslog"
)

const (
	privateFileMode = 0600
	// owner can make/remove files inside the directory
	privateDirMode = 0700
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/etcd/pkg", "fileutil")
)

// IsDirWriteable checks if dir is writable by writing and removing a file
// to dir. It returns nil if dir is writable.
func IsDirWriteable(dir string) error {
	f := path.Join(dir, ".touch")
	if err := ioutil.WriteFile(f, []byte(""), privateFileMode); err != nil {
		return err
	}
	return os.Remove(f)
}

// ReadDir returns the filenames in the given directory in sorted order.
func ReadDir(dirpath string) ([]string, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

// TouchDirAll is similar to os.MkdirAll. It creates directories with 0700 permission if any directory
// does not exists. TouchDirAll also ensures the given directory is writable.
func TouchDirAll(dir string) error {
	err := os.MkdirAll(dir, privateDirMode)
	if err != nil && err != os.ErrExist {
		return err
	}
	return IsDirWriteable(dir)
}

func Exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
