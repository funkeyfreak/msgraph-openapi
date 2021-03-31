package yq

import (
	"container/list"
	"io"
	"sync"

	internalio "github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

type FileName string

func (fn FileName) loadFile() (*list.List, error) {

}

func (fn FileName) ToString() string {
	return string(fn)
}

// TODO: Move FileState to state package
type FileSet struct {
	isConcurrent bool
	fileMap      map[FileName]*list.List
	mu           sync.Mutex
}

func NewFiles(makeConcurrent bool, fileNames ...string) *FileSet {
	fileSet := FileSet{isConcurrent: makeConcurrent}
	for _, fileName := range fileNames {
		fileSet.fileMap[FileName(fileName)] = nil
	}

	return &fileSet
}

func (fs *FileSet) concurrencyHandler(f func()) {
	if fs.isConcurrent {
		fs.mu.Lock()
		f()
		fs.mu.Unlock()
	} else {
		f()
	}
}

// TODO: make this handle both yaml and json files
func (fs *FileSet) readYamlFile(reader io.Reader, fileName FileName, fileIndex int) (*list.List, error) {
	decoder := yaml.NewDecoder(reader)
	inputList := list.New()
	var currentIndex uint = 0

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			switch reader := reader.(type) {
			case internalio.File:
				reader.Close()
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &yqlib.CandidateNode{
			Document:         currentIndex,
			Filename:         fileName.ToString(),
			Node:             &dataBucket,
			FileIndex:        fileIndex,
			EvaluateTogether: true,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}

func (fs *FileSet) SetFile(fileName FileName) {
	// TODO: May need to lock here

}

func (fs *FileSet) LoadFiles() {
	fs.concurrencyHandler(func() {
		wg := sync.WaitGroup{}
		id := 1

		// size of buffered chan is the size of our input filess
		in := make(chan FileName, len(fs.fileMap))

		out := make(chan *list.List, len(fs.fileMap))
		for fileName, list := range f {
			if list == nil {
				wg.Add()
				id++
			}
		}
	})
}
