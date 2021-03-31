package io

import (
	"os"

	"github.com/spf13/afero"
)

const (
	WriteFileFSFileMode = 0644 // Default FSFileMode when writing files
)

// wrapper for fs.WriteFile
func (fs FileSystem) WriteFile(path string, data []byte) error {
	return afero.WriteFile(fs.fileSystem, path, data, WriteFileFSFileMode)
}

// helper for creating an empty file
func (fs FileSystem) CreateEmptyFile(path string) error {
	emptyData := []byte("")

	return fs.WriteFile(path, emptyData)
}

// wrapper for fs.Create
func (fs FileSystem) Create(path string) (File, error) {
	return fs.fileSystem.Create(path)
}

// wrapper for afero.ReadFile
func (fs FileSystem) ReadFile(path string) ([]byte, error) {
	return afero.ReadFile(fs.fileSystem, path)
}

// wrapper for fs.Open
func (fs FileSystem) Open(path string) (File, error) {
	return fs.fileSystem.Open(path)
}

// wrapper for fs.OpenFile with the flags os.O_APPEND|os.O_WRONLY and the perm os.ModeAppend
func (fs FileSystem) OpenFile(path string) (File, error) {
	return fs.fileSystem.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

// wrapper for os.CreateTemp
func (fs FileSystem) CreateTemp(dir string, pattern string) (File, error) {
	return os.CreateTemp(dir, pattern)
}
