package io

import (
	"errors"
	"os"

	"github.com/spf13/afero"
)

const (
	MkdirAllFSFileMode = 0755 // Default FSFileMode when creating folders
	DefaultTmpDirName  = "cmd-name"
)

// wrapper for os.ReadDir
func (fs FileSystem) ReadDir(path string) (FileInfo, error) {
	return afero.ReadDir(fs.fileSystem, path)
}

// wrapper for os.MkdirAll
func (fs FileSystem) MkdirAll(path string) error {
	return fs.fileSystem.MkdirAll(path, MkdirAllFSFileMode)
}

// wrapper for os.RemoveAll
func (fs FileSystem) RemoveAll(path string) error {
	return fs.fileSystem.RemoveAll(path)
}

// helper method which returns true if the item in a path is a directory
func (fs FileSystem) IsDirectory(path string) bool {
	_, err := fs.ReadDir(path)
	if err != nil {
		return false
	}

	return true
}

func (fs FileInfo) MkdirTemp() (string, error) {
	return os.MkdirTemp("", DefaultTmpDirName)
}

func (fs FileSystem) FetchTempDir() (FileInfo, error) {
	switch fs.fileSystem.(type) {
	case afero.OsFs:
		break
	default:
		return nil, errors.New("Not supported on this file system")
	}
	return fs.ReadDir(os.TempDir())
}
