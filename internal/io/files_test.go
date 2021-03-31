package io_test

import (
	"os"
	"testing"

	"github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/stretchr/testify/require"
)

const (
	TestFilePath   = "./testdata/testFile.file"
	CreateFilePath = "./testdata/testFile2.file"
	InvalidPath    = "DoesNotExist"
	EmptyPath      = ""
)

/*
	WrireFile
*/
func Test_WriteFile(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, CreateFilePath)
	data := []byte("data")

	t.Run("Succeeds", func(t *testing.T) {
		t.Cleanup(func() { os.Remove(CreateFilePath) })
		err := testFileSystem.WriteFile(path, data)
		require.NoError(t, err)
	})

	t.Run("FailsIfInvalidPath", func(t *testing.T) {
		err := testFileSystem.WriteFile(EmptyPath, data)
		require.Error(t, err)
	})
}

/*
	CreateEmptyFile
*/
func Test_CreateEmptyFile(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, CreateFilePath)

	t.Run("Succeeds", func(t *testing.T) {
		t.Cleanup(func() { os.Remove(path) })
		require.NoError(t, testFileSystem.CreateEmptyFile(path))
	})

	t.Run("FailsWhenPathIsEmpty", func(t *testing.T) {
		require.Error(t, testFileSystem.CreateEmptyFile(EmptyPath))
	})
}

/*
	Create
*/
func Test_Create_Succeeds(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, CreateFilePath)

	t.Run("Succeeds", func(t *testing.T) {
		t.Cleanup(func() { os.Remove(path) })
		file, err := testFileSystem.Create(path)
		defer file.Close()
		require.NoError(t, err)
		require.NotNil(t, file)
	})

	t.Run("FailsWithPathIsEmpty", func(t *testing.T) {
		_, err := testFileSystem.Create(EmptyPath)
		require.Error(t, err)
	})
}

/*
	ReadFile
*/
func Test_ReadFile(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, TestFilePath)

	t.Run("Suceeds", func(t *testing.T) {
		file, err := testFileSystem.ReadFile(path)
		require.NoError(t, err)
		require.NotNil(t, file)
	})

	t.Run("ReturnsExpected", func(t *testing.T) {
		file, _ := testFileSystem.ReadFile(path)
		expected, _ := os.ReadFile(path)
		require.Equal(t, file, expected)
	})

	t.Run("FailsWhenInvalidPath", func(t *testing.T) {
		path = InvalidPath
		_, err := testFileSystem.ReadFile(path)
		require.Error(t, err)
	})

}

/*
	Open
*/
func Test_Open(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, TestFilePath)

	t.Run("Succeeds", func(t *testing.T) {
		file, err := testFileSystem.Open(path)
		defer file.Close()
		require.NoError(t, err)
		require.NotNil(t, file)
	})

	t.Run("FailsWhenInvalidPath", func(t *testing.T) {
		path = "invalidPath"
		_, err := testFileSystem.Open(path)
		require.Error(t, err)
	})
}

/*
	OpenFileForWriting
*/
func Test_OpenFile(t *testing.T) {
	testFileSystem, _ := io.NewFileSystem(io.OSFileSystemType)
	path := getTestFilePath(t, CreateFilePath)
	data := []byte("data")

	t.Run("Succeeds", func(t *testing.T) {
		os.Create(CreateFilePath)
		t.Cleanup(func() { os.Remove(path) })
		file, err := testFileSystem.OpenFile(path)
		defer file.Close()
		require.NoError(t, err)
		require.NotNil(t, file)

		_, err = file.Write(data)
		require.NoError(t, err)
	})
}
