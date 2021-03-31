package datalayer

import (
	"errors"
	"sync"

	internalio "github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/funkeyfreak/msgraph-openapi/internal/utils"
)

const (
	CacheMaximumChanBufferSize = 10
	CacheMaximumWorkerPoolSize = 10
	MaxWorkerCountPerDataType  = 3
)

var (
	ErrFailedToLoadFileIntoFileSystemCache  = errors.New("failed to load file into file system cache")
	ErrFailedToSaveFileDatModelDataModelSet = errors.New("failed to save data model cache at key is occupied")
	ErrFileNotLoadedInSystemCache           = errors.New("file is not loaded in file system cache")
	ErrModelNotLoadedInSystemCache          = errors.New("model is not loaded in file system cache")
	ErrFileSystemCacheNotAFile              = errors.New("the item provided is not a file")
)

// Caching system used for increased performance
type FileCacheSystem struct {
	fsAbs          *internalio.FileSystem
	isConcurrent   bool
	FileModelCache map[string]DataModel
	FileInfoCache  map[string]File
	mu             sync.Mutex
}

func (fs *FileCacheSystem) concurrencyHandler(f func()) {
	if fs.isConcurrent {
		fs.mu.Lock()
		f()
		fs.mu.Unlock()
	} else {
		f()
	}
}

func (fs *FileCacheSystem) worker(id int, in <-chan workerJob, out chan<- error) {
	for job := range in {
		_, err := fs.LoadFileDataModel(job.dmt, job.file)
		out <- err
	}
}

func (fs *FileCacheSystem) loadDataModelType(dms []DataModelType, files []File) []error {
	//validatedFiles := make([]File, 0)
	errors := make([]error, 0)

	for _, file := range files {
		for _, ds := range dms {
			_, err := fs.LoadFileDataModel(ds, file)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

type workerJob struct {
	file File
	dmt  DataModelType
}

func (fs *FileCacheSystem) loadDataModelTypeConcurrent(dmt []DataModelType, files []File) []error {
	numFiles := len(files)
	// size of buffered chan is the size of our input files
	dataModelTypeCount := numFiles * len(dmt)
	in := make(chan workerJob, numFiles)
	out := make(chan error, numFiles)

	// adjust pool size and count of wokers per thread
	workerCountPerDataType := utils.Min(MaxWorkerCountPerDataType, len(dmt))
	cacheMaximumWorkerPoolSize := utils.Min(CacheMaximumWorkerPoolSize, dataModelTypeCount*workerCountPerDataType)

	workerJobs := make([]workerJob, 0)

	// create array of workerJob
	for _, dt := range dmt {
		for _, file := range files {
			workerJobs = append(workerJobs, workerJob{file: file, dmt: dt})
		}
	}

	// create worker pool
	for w := 1; w <= cacheMaximumWorkerPoolSize; w++ {
		go fs.worker(w, in, out)
	}

	// batch in all files for processing
	for _, job := range workerJobs {
		in <- job
	}
	close(in)

	errors := make([]error, 0)
	for i := 0; i < len(dmt)*len(files); i++ {
		err := <-out
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// Create a new instance of FileCacheSystem
func NewFileCacheSystem(isConcurrent bool, fsAbs *internalio.FileSystem) *FileCacheSystem {
	return &FileCacheSystem{
		isConcurrent:   isConcurrent,
		fsAbs:          fsAbs,
		FileModelCache: make(map[string]DataModel),
		FileInfoCache:  make(map[string]File),
	}
}

// return the number of files in the cache system
func (fs *FileCacheSystem) Len() int {
	return len(fs.FileInfoCache)
}

// return the number of files in the cache system
func (fs *FileCacheSystem) MetadataLen() int {
	return len(fs.FileModelCache)
}

// TODO: Add the following methods:
// Remove file by File object
// Add File Model by path

// Fetch file
func (fs *FileCacheSystem) FetchFileFromCache(filePath string) (*File, error) {
	var file *File
	var err error

	fs.concurrencyHandler(func() {
		if m, ok := fs.FileInfoCache[filePath]; ok {
			file = &m
		} else {
			err = ErrModelNotLoadedInSystemCache
		}
	})

	return file, err
}

func (fs *FileCacheSystem) FetchMultipleFilesFromCache(filePaths ...string) ([]*File, error) {
	var files []*File
	var err error
	for _, filePath := range filePaths {
		var m *File
		m, err = fs.FetchFileFromCache(filePath)
		if err != nil {
			return nil, err
		}
		files = append(files, m)
	}

	return files, err
}

func (fs *FileCacheSystem) FetchDataModelFromCache(filePath string) (*DataModel, error) {
	var model *DataModel
	var err error

	fs.concurrencyHandler(func() {
		if m, ok := fs.FileModelCache[filePath]; ok {
			model = &m
		} else {
			err = ErrModelNotLoadedInSystemCache
		}
	})

	return model, err
}

func (fs *FileCacheSystem) FetchMultipleDataModelsFromCache(filePaths ...string) ([]*DataModel, error) {
	var models []*DataModel
	var err error
	for _, filePath := range filePaths {
		var m *DataModel
		m, err = fs.FetchDataModelFromCache(filePath)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}

	return models, err
}

func (fs *FileCacheSystem) AddMultipleToCache(files ...File) error {
	var err error

	for _, filePath := range files {
		err = fs.AddToCache(filePath)
		if err != nil {
			return err
		}
	}

	return err
}

func (fs *FileCacheSystem) AddToCache(file File) error {
	if file == nil {
		return ErrFileSystemCacheNotAFile
	}

	fs.concurrencyHandler(func() {
		fs.FileInfoCache[file.Path()] = file
		fs.FileModelCache[file.Path()] = NewDataModel()
	})

	return nil
}

func (fs *FileCacheSystem) RemoveFromCache(filePath string) error {
	var err error
	fs.concurrencyHandler(func() {
		if _, ok := fs.FileInfoCache[filePath]; ok {
			delete(fs.FileInfoCache, filePath)
		} else {
			err = ErrFileNotLoadedInSystemCache
		}
		if _, ok := fs.FileModelCache[filePath]; ok {
			delete(fs.FileModelCache, filePath)
		} else {
			err = ErrModelNotLoadedInSystemCache
		}
	})

	return err
}

func (fs *FileCacheSystem) RemoveMultipleFromCache(filePaths ...string) []error {
	errors := make([]error, 0)
	for _, filePath := range filePaths {
		err := fs.RemoveFromCache(filePath)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func (fs *FileCacheSystem) LoadFileDataModel(dmt DataModelType, file File) (interface{}, error) {
	if !dmt.IsValid() {
		return nil, ErrUnhandledDataModelType
	}

	if _, ok := fs.FileModelCache[file.Path()]; !ok {
		return nil, ErrFileNotLoadedInSystemCache
	}

	var err error
	var data interface{}
	switch dmt {
	case DataModelTypeList:
		data, err = file.LoadListModel(fs.fsAbs)
	default:
		err = ErrUnhandledDataModelType
	}

	if err != nil {
		return nil, err
	}

	// make sure we can save to data model - if we cannot, then delete then attempt to re-add
	//	and don't forget to lock!
	fs.concurrencyHandler(func() {
		if !fs.FileModelCache[file.Path()].AddDataModel(data) {
			if !fs.FileModelCache[file.Path()].RemoveModel(dmt) || !fs.FileModelCache[file.Path()].AddDataModel(data) {
				err = ErrFailedToSaveFileDatModelDataModelSet
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs *FileCacheSystem) LoadDataModelTypes(dms []DataModelType, files []File) []error {
	var errors []error

	if fs.isConcurrent {
		errors = fs.loadDataModelTypeConcurrent(dms, files)
	} else {
		errors = fs.loadDataModelType(dms, files)
	}
	return errors
}

func (fs *FileCacheSystem) LoadAllDataModelTypes(dms ...DataModelType) []error {
	var errors []error

	files := make([]File, 0)

	for _, file := range fs.FileInfoCache {
		files = append(files, file)
	}

	if len(files) == 0 {
		return []error{ErrFileNotLoadedInSystemCache}
	}

	if fs.isConcurrent {
		errors = fs.loadDataModelTypeConcurrent(dms, files)
	} else {
		errors = fs.loadDataModelType(dms, files)
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
