package io

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

type FileSystemType string

const (
	OSFileSystemType  FileSystemType = "OSFileSystemType"
	MemFileSystemType FileSystemType = "MemFileSystemType"
)

func (fT FileSystemType) IsValid() error {
	switch fT {
	case OSFileSystemType, MemFileSystemType:
		return nil
	}

	return fmt.Errorf("FileSystemType %v is invalid", fT)
}

type File afero.File

type FileInfo []os.FileInfo

type FileSystem struct {
	fileSystem afero.Fs
}

func NewFileSystem(filesystemType FileSystemType) (*FileSystem, error) {
	var fs *FileSystem
	switch filesystemType {
	case MemFileSystemType:
		fs = &FileSystem{fileSystem: afero.NewMemMapFs()}
	case OSFileSystemType:
		fs = &FileSystem{fileSystem: afero.NewOsFs()}
	}

	return fs, filesystemType.IsValid()
}
