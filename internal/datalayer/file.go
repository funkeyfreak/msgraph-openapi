package datalayer

import (
	"container/list"
	"errors"
	"io"
	"os"

	internalio "github.com/funkeyfreak/msgraph-openapi/internal/io"
)

var (
	ErrFileInvalidFilePath = errors.New("provided file path is invalid")
	ErrFileInvalidFileType = errors.New("unhandled file type")
	FileTypeNotCreated     = errors.New("cannot create provided file type")
)

type FileType string

const (
	FileTypeJson FileType = "json"
	FileTypeYaml FileType = "yaml"
)

func (ft FileType) IsValid() bool {
	switch ft {
	case FileTypeJson, FileTypeYaml:
		return true
	}
	return false
}

func HandledFileTypes() []FileType {
	return []FileType{
		FileTypeJson,
		FileTypeYaml,
	}
}

type fileMetaData struct {
	FileName string
	FileType FileType
	FilePath string
}

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}

type File interface {
	Name() string
	Path() string

	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error

	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder

	LoadListModel(fs *internalio.FileSystem) (*list.List, error)
}

func NewFile(filePath string, fsAbs *internalio.FileSystem) (File, error) {
	// shortest file path is of len = 5 - t.yml
	if len(filePath) < 5 {
		return nil, ErrFileInvalidFilePath
	}
	exists, err := fsAbs.Exists(filePath)
	fileType := FileType([]byte(fsAbs.Ext(filePath))[1:])

	if !exists {
		if err != nil {
			return nil, err
		}
		return nil, os.ErrNotExist
	} else if err == nil && exists && !fileType.IsValid() {
		return nil, ErrFileInvalidFileType // fmt.Errorf("file %v with extension %v is none of %+v", filePath, string(fileType), HandledFileTypes())
	}

	if err != nil {
		return nil, err
	}

	fmd := fileMetaData{
		FileName: fsAbs.Base(filePath),
		FileType: fileType,
		FilePath: filePath,
	}
	switch fileType {
	case FileTypeJson:
		//TODO
	case FileTypeYaml:
		return newYamlFile(fmd), nil
	}

	return nil, FileTypeNotCreated
}
