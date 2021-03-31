package io_test

import (
	"os"
	"testing"

	"github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/stretchr/testify/require"
)

const (
	TestDirPath   = "./testdata"
	MkdirAllPath  = "./testdata/testDataDir/subDir"
	RemoveAllPath = "./testdata/testDataDir"
	InvalidDir    = "./testdata/DoesNotExist"
	EmptyDir      = ""
)

/*
	ReadDir
*/
func Test_ReadDir_Succeds(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, TestDirPath)

	t.Run("Succeeds", func(t *testing.T) {
		folder, err := testFileSystem.ReadDir(path)
		require.NoError(t, err)
		require.NotNil(t, folder)
	})

	t.Run("FailsIfInvalidDir", func(t *testing.T) {
		path := getTestFilePath(t, InvalidDir)
		_, err := testFileSystem.ReadDir(path)
		require.Error(t, err)
	})
}

/*
	MkdirAll
*/
func Test_Mkdir(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, MkdirAllPath)

	t.Run("Succeeds", func(t *testing.T) {
		t.Cleanup(func() {
			path = getTestFilePath(t, RemoveAllPath)
			os.RemoveAll(path)
		})
		require.NoError(t, testFileSystem.MkdirAll(path))
	})

	t.Run("FailsIfEmptyFolder", func(t *testing.T) {
		require.Error(t, testFileSystem.MkdirAll(EmptyDir))
	})
}

/*
	RemoveAll
*/
func Test_RemoveAll(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, RemoveAllPath)

	t.Run("Succeeds", func(t *testing.T) {
		os.MkdirAll(getTestFilePath(t, MkdirAllPath), os.ModeAppend)
		require.NoError(t, testFileSystem.RemoveAll(path))
	})
}

/*
	IsDirectory
*/
func Test_IsDirectory(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, TestDirPath)

	t.Run("Succeeds", func(t *testing.T) {
		t.Run("TrueWhenDirectory", func(t *testing.T) {
			isDir := testFileSystem.IsDirectory(path)
			require.True(t, isDir)
		})

		t.Run("FalseWhenNotDirectory", func(t *testing.T) {
			path = getTestFilePath(t, TestFilePath)
			isDir := testFileSystem.IsDirectory(path)
			require.False(t, isDir)
		})
	})

}
