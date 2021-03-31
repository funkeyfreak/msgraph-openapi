//TODO: Test returned error types
package datalayer_test

import (
	"container/list"
	"testing"

	"github.com/funkeyfreak/msgraph-openapi/internal/datalayer"
	"github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/stretchr/testify/require"
)

func Test_FileCacheSystem_NewFileCacheFileSystem(t *testing.T) {
	osFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)

	t.Run("Succeeds_IsConcurrent", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(true, osFileSystem)
		require.NotNil(t, cache)
	})

	t.Run("Succeeds", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		require.NotNil(t, cache)
	})
}

func Test_FileCacheSystem_AddToCache(t *testing.T) {
	osFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	yamlFile, _ := datalayer.NewFile(TestFilePathYaml, osFileSystem)

	t.Run("Succeeds_IsConcurrent", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(true, osFileSystem)
		require.NotNil(t, cache)
		err := cache.AddToCache(yamlFile)
		require.NoError(t, err)
		require.Contains(t, cache.FileInfoCache, TestFilePathYaml)
		file := cache.FileInfoCache[TestFilePathYaml]
		require.Equal(t, file.Name(), osFileSystem.Base(TestFilePathYaml))
		require.Equal(t, file.Path(), TestFilePathYaml)
	})

	t.Run("Succeeds", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		require.NotNil(t, cache)
		err := cache.AddToCache(yamlFile)
		require.NoError(t, err)
		require.Contains(t, cache.FileInfoCache, TestFilePathYaml)
		file := cache.FileInfoCache[TestFilePathYaml]
		require.Equal(t, file.Name(), osFileSystem.Base(TestFilePathYaml))
		require.Equal(t, file.Path(), TestFilePathYaml)
	})

	t.Run("Fails", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		err := cache.AddToCache(nil)
		require.Error(t, err)
	})

	t.Run("Fails_IsConcurrent", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		err := cache.AddToCache(nil)
		require.Error(t, err)
	})
}

func Test_FileCacheSystem_AddMultipleToCache(t *testing.T) {
	osFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	yamlFile, _ := datalayer.NewFile(TestFilePathYaml, osFileSystem)

	t.Run("Succeeds_IsConcurrent", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(true, osFileSystem)
		require.NotNil(t, cache)
		err := cache.AddMultipleToCache(yamlFile)
		require.NoError(t, err)
		require.Contains(t, cache.FileInfoCache, TestFilePathYaml)
		file := cache.FileInfoCache[TestFilePathYaml]
		require.Equal(t, file.Name(), osFileSystem.Base(TestFilePathYaml))
		require.Equal(t, file.Path(), TestFilePathYaml)
	})

	t.Run("Succeeds", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		require.NotNil(t, cache)
		err := cache.AddMultipleToCache(yamlFile)
		require.NoError(t, err)
		require.Contains(t, cache.FileInfoCache, TestFilePathYaml)
		file := cache.FileInfoCache[TestFilePathYaml]
		require.Equal(t, file.Name(), osFileSystem.Base(TestFilePathYaml))
		require.Equal(t, file.Path(), TestFilePathYaml)
	})

	t.Run("Fails", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		err := cache.AddMultipleToCache(nil)
		require.Error(t, err)
	})

	t.Run("Fails_IsConcurrent", func(t *testing.T) {
		cache := datalayer.NewFileCacheSystem(false, osFileSystem)
		err := cache.AddMultipleToCache(nil)
		require.Error(t, err)
	})
}

func Test_FileCacheSystem_RemoveFromCache(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		require.NoError(t, cacheSystem.RemoveFromCache(yamlFile.Path()))
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		require.NoError(t, concurrentCacheSystem.RemoveFromCache(yamlFile.Path()))
	})
}

func Test_FileCacheSystem_RemoveMultipleFromCache(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	pathArr := make([]string, 0)
	for _, file := range files {
		pathArr = append(pathArr, file.Path())
	}

	t.Run("Succeeds", func(t *testing.T) {
		require.Equal(t, cacheSystem.Len(), len(filePaths))

		errors := cacheSystem.RemoveMultipleFromCache(pathArr...)
		require.Nil(t, errors)
		require.Zero(t, cacheSystem.Len())
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		require.Equal(t, concurrentCacheSystem.Len(), len(filePaths))

		errors := concurrentCacheSystem.RemoveMultipleFromCache(pathArr...)
		require.Nil(t, errors)
		require.Zero(t, concurrentCacheSystem.Len())
	})

	t.Run("Fails", func(t *testing.T) {
		errors := cacheSystem.RemoveMultipleFromCache(pathArr...)
		require.NotNil(t, errors)
		require.Len(t, errors, len(files))
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		errors := concurrentCacheSystem.RemoveMultipleFromCache(pathArr...)
		require.NotNil(t, errors)
		require.Len(t, errors, len(files))
	})
}

func Test_FileCacheSystem_FetchFileFromCache(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		for _, file := range files {
			f, err := cacheSystem.FetchFileFromCache(file.Path())
			require.NoError(t, err)
			require.NotNil(t, f)
		}
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		for _, file := range files {
			f, err := concurrentCacheSystem.FetchFileFromCache(file.Path())
			require.NoError(t, err)
			require.NotNil(t, f)
		}
	})

	t.Run("Fails", func(t *testing.T) {
		f, err := cacheSystem.FetchFileFromCache(InvalidFilePath)
		require.Error(t, err)
		require.Nil(t, f)
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		f, err := concurrentCacheSystem.FetchFileFromCache(InvalidFilePath)
		require.Error(t, err)
		require.Nil(t, f)
	})
}

func Test_FileCacheSystem_FetchMultipleFilesFromCache(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		f, err := cacheSystem.FetchMultipleFilesFromCache(filePaths...)
		require.NoError(t, err)
		require.Equal(t, len(f), len(filePaths))
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		f, err := concurrentCacheSystem.FetchMultipleFilesFromCache(filePaths...)
		require.NoError(t, err)
		require.Equal(t, len(f), len(filePaths))
	})

	t.Run("Fails", func(t *testing.T) {
		f, err := cacheSystem.FetchMultipleFilesFromCache(InvalidFilePath)
		require.Error(t, err)
		require.Nil(t, f)
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		f, err := concurrentCacheSystem.FetchMultipleFilesFromCache(InvalidFilePath)
		require.Error(t, err)
		require.Nil(t, f)
	})
}

func Test_FileCacheSystem_FetchDataModelFromCache(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		file, err := cacheSystem.FetchDataModelFromCache(yamlFile.Path())
		require.NoError(t, err)
		require.NotNil(t, file)
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		file, err := concurrentCacheSystem.FetchDataModelFromCache(yamlFile.Path())
		require.NoError(t, err)
		require.NotNil(t, file)
	})

	// remove file to test if LoadFileDataModel fails when querying by a file which no
	//	longer is loaded in teh cache
	cacheSystem.RemoveFromCache(yamlFile.Path())
	concurrentCacheSystem.RemoveFromCache(yamlFile.Path())

	t.Run("Fails", func(t *testing.T) {
		file, err := cacheSystem.FetchDataModelFromCache(yamlFile.Path())
		require.Error(t, err)
		require.Nil(t, file)
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		file, err := concurrentCacheSystem.FetchDataModelFromCache(yamlFile.Path())
		require.Error(t, err)
		require.Nil(t, file)
	})
}

func Test_FileCacheSystem_FetchMultipleDataModelsFromCache(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		models, err := cacheSystem.FetchMultipleDataModelsFromCache(filePaths...)
		require.NoError(t, err)
		require.NotNil(t, models)
		require.Len(t, models, cacheSystem.MetadataLen())
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		models, err := concurrentCacheSystem.FetchMultipleDataModelsFromCache(filePaths...)
		require.NoError(t, err)
		require.NotNil(t, models)
		require.Len(t, models, concurrentCacheSystem.MetadataLen())
	})

	// remove file to test if LoadFileDataModel fails when querying by a file which no
	//	longer is loaded in teh cache
	cacheSystem.RemoveFromCache(yamlFile.Path())
	concurrentCacheSystem.RemoveFromCache(yamlFile.Path())

	t.Run("Fails", func(t *testing.T) {
		models, err := cacheSystem.FetchMultipleDataModelsFromCache(filePaths...)
		require.Error(t, err)
		require.Nil(t, models)
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		models, err := cacheSystem.FetchMultipleDataModelsFromCache(filePaths...)
		require.Error(t, err)
		require.Nil(t, models)
	})
}

func Test_FileCacheSystem_LoadFileDataModel(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		for _, file := range files {
			f, err := cacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, file)
			require.NotNil(t, f)
			require.NoError(t, err)
			require.IsType(t, &list.List{}, f)
		}
	})

	t.Run("Succeeds_WhenLoadedBackToBack", func(t *testing.T) {
		for _, file := range files {
			f, err := cacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, file)
			require.NotNil(t, f)
			require.NoError(t, err)
			require.IsType(t, &list.List{}, f)
			// This should call DataModel.Remove() and then call DataModel.Add(), as there is not a DataModel.Update function
			f, err = cacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, file)
			require.NotNil(t, f)
			require.NoError(t, err)
			require.IsType(t, &list.List{}, f)
		}
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		for _, file := range files {
			f, err := concurrentCacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, file)
			require.NotNil(t, f)
			require.NoError(t, err)
			require.IsType(t, &list.List{}, f)
		}
	})

	t.Run("Succeeds_WhenLoadedBackToBack_WhenConcurrent", func(t *testing.T) {
		for _, file := range files {
			f, err := concurrentCacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, file)
			require.NotNil(t, f)
			require.NoError(t, err)
			require.IsType(t, &list.List{}, f)
			// This should call DataModel.Remove() and then call DataModel.Add(), as there is not a DataModel.Update function
			f, err = concurrentCacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, file)
			require.NotNil(t, f)
			require.NoError(t, err)
			require.IsType(t, &list.List{}, f)
		}
	})

	t.Run("Fails", func(t *testing.T) {
		dm, err := cacheSystem.LoadFileDataModel(datalayer.DataModelTypeByteArray, yamlFile)
		require.Error(t, err)
		require.Nil(t, dm)
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		dm, err := concurrentCacheSystem.LoadFileDataModel(datalayer.DataModelTypeByteArray, yamlFile)
		require.Error(t, err)
		require.Nil(t, dm)
	})

	t.Run("Fails_WhenInvalidType", func(t *testing.T) {
		dm, err := cacheSystem.LoadFileDataModel(InvalidDataModelType, yamlFile)
		require.Error(t, err)
		require.EqualError(t, err, datalayer.ErrUnhandledDataModelType.Error())
		require.Nil(t, dm)
	})

	t.Run("Fails_WhenInvalidType_WhenConcurrent", func(t *testing.T) {
		dm, err := concurrentCacheSystem.LoadFileDataModel(InvalidDataModelType, yamlFile)
		require.Error(t, err)
		require.EqualError(t, err, datalayer.ErrUnhandledDataModelType.Error())
		require.Nil(t, dm)
	})

	// remove file to test if LoadFileDataModel fails when querying by a file which no
	//	longer is loaded in teh cache
	cacheSystem.RemoveFromCache(yamlFile.Path())
	concurrentCacheSystem.RemoveFromCache(yamlFile.Path())

	t.Run("Fails_WhenFileIsNotLoadedInCache", func(t *testing.T) {
		dm, err := cacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, yamlFile)
		require.Error(t, err)
		require.EqualError(t, err, datalayer.ErrFileNotLoadedInSystemCache.Error())
		require.Nil(t, dm)
	})

	t.Run("Fails_WhenFileIsNotLoadedInCache_WhenConcurrent", func(t *testing.T) {
		dm, err := concurrentCacheSystem.LoadFileDataModel(datalayer.DataModelTypeList, yamlFile)
		require.Error(t, err)
		require.EqualError(t, err, datalayer.ErrFileNotLoadedInSystemCache.Error())
		require.Nil(t, dm)
	})
}

func Test_FileCacheSystem_LoadDataModelTypes(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		errors := cacheSystem.LoadDataModelTypes([]datalayer.DataModelType{datalayer.DataModelTypeList}, files)
		require.Empty(t, errors)
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		errors := concurrentCacheSystem.LoadDataModelTypes([]datalayer.DataModelType{datalayer.DataModelTypeList}, files)
		require.Empty(t, errors)
	})

	// remove file to test if LoadFileDataModel fails when querying by a file which no
	//	longer is loaded in teh cache
	cacheSystem.RemoveFromCache(yamlFile.Path())
	concurrentCacheSystem.RemoveFromCache(yamlFile.Path())

	t.Run("SucceedsWithErrors_WhenAFileIsRemoved", func(t *testing.T) {
		errors := cacheSystem.LoadDataModelTypes([]datalayer.DataModelType{datalayer.DataModelTypeList}, files)
		require.Len(t, errors, 1)
	})

	t.Run("SucceedsWithErrors_WhenAFileIsRemoved_WhenConcurrent", func(t *testing.T) {
		errors := concurrentCacheSystem.LoadDataModelTypes([]datalayer.DataModelType{datalayer.DataModelTypeList}, files)
		require.Len(t, errors, 1)
	})
}

func Test_FileCacheSystem_LoadAllDataModelTypes(t *testing.T) {
	cacheSystem := datalayer.NewFileCacheSystem(false, osFileSystem)
	concurrentCacheSystem := datalayer.NewFileCacheSystem(true, osFileSystem)

	cacheSystem.AddMultipleToCache(files...)
	concurrentCacheSystem.AddMultipleToCache(files...)

	t.Run("Succeeds", func(t *testing.T) {
		require.Nil(t, cacheSystem.LoadAllDataModelTypes(datalayer.DataModelTypeList))
	})

	t.Run("Succeeds_WhenConcurrent", func(t *testing.T) {
		require.Nil(t, concurrentCacheSystem.LoadAllDataModelTypes(datalayer.DataModelTypeList))
	})

	t.Run("Fails_WhenIncorrectTypeIsGiven", func(t *testing.T) {
		errors := cacheSystem.LoadAllDataModelTypes(InvalidDataModelType)
		require.NotNil(t, errors)
		require.Len(t, errors, 3)
	})

	t.Run("Fails_WhenIncorrectTypeIsGiven_WhenConcurrent", func(t *testing.T) {
		errors := concurrentCacheSystem.LoadAllDataModelTypes(InvalidDataModelType)
		require.NotNil(t, errors)
		require.Len(t, errors, 3)
	})

	// remove all files to verify behaviour when there are nothing in our list
	cacheSystem.RemoveMultipleFromCache(filePaths...)
	concurrentCacheSystem.RemoveMultipleFromCache(filePaths...)

	t.Run("Fails", func(t *testing.T) {
		errors := cacheSystem.LoadAllDataModelTypes(InvalidDataModelType)
		require.NotNil(t, errors)
		require.Len(t, errors, 1)
		require.EqualError(t, errors[0], datalayer.ErrFileNotLoadedInSystemCache.Error())
	})

	t.Run("Fails_WhenConcurrent", func(t *testing.T) {
		errors := concurrentCacheSystem.LoadAllDataModelTypes(datalayer.DataModelTypeList)
		require.NotNil(t, errors)
		require.Len(t, errors, 1)
		require.EqualError(t, errors[0], datalayer.ErrFileNotLoadedInSystemCache.Error())
	})
}
