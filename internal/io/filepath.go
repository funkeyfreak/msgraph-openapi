package io

import (
	"path/filepath"

	"github.com/spf13/afero"
)

func (fs *FileSystem) IsDir(filePath string) bool {
	return true
}

func (fs *FileSystem) Exists(filePath string) (bool, error) {
	return afero.Exists(fs.fileSystem, filePath)
}

func (fs *FileSystem) Ext(path string) string {
	return filepath.Ext(path)
}

func (fs *FileSystem) Base(path string) string {
	return filepath.Base(path)
}
