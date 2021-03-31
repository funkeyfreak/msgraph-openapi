package datalayer

import (
	"container/list"
	"errors"
	"io"

	"gopkg.in/yaml.v3"

	internalio "github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

var (
	YamlEncoderNotClosedError error = errors.New("encoder is not closed")
)

type yamlFile struct {
	decoder       *yaml.Decoder
	encoder       *yaml.Encoder
	encoderClosed bool
	*fileMetaData
}

func newYamlFile(fmd fileMetaData) *yamlFile {
	return &yamlFile{
		fileMetaData:  &fmd,
		encoderClosed: true,
	}
}

func (y *yamlFile) LoadListModel(fs *internalio.FileSystem) (*list.List, error) {
	reader, err := fs.Open(y.FilePath)
	if err != nil {
		return nil, err
	}

	inputList := list.New()
	fileIndex := 0
	var currentIndex uint = 0
	decoder := y.NewDecoder(reader)

	for {
		var dataBucket yaml.Node
		err := decoder.Decode(&dataBucket)

		if err == io.EOF {
			switch reader := reader.(type) {
			case internalio.File:
				reader.Close()
			}
			return inputList, nil
		} else if err != nil {
			return nil, err
		}
		candidateNode := &yqlib.CandidateNode{
			Document:         currentIndex,
			Filename:         y.FileName,
			Node:             &dataBucket,
			FileIndex:        fileIndex,
			EvaluateTogether: true,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}

func (y *yamlFile) Name() string {
	return y.fileMetaData.FileName
}

func (y *yamlFile) Path() string {
	return y.fileMetaData.FilePath
}

func (y *yamlFile) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *yamlFile) Unmarshal(data []byte, v interface{}) error {
	return y.Unmarshal(data, v)
}

func (y *yamlFile) NewEncoder(w io.Writer) Encoder {
	return yaml.NewEncoder(w)
}

func (y *yamlFile) NewDecoder(r io.Reader) Decoder {
	return yaml.NewDecoder(r)
}
