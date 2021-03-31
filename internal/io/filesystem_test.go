package io_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/funkeyfreak/msgraph-openapi/internal/io"
	"github.com/stretchr/testify/require"
)

func getTestFilePath(t *testing.T, path string) string {
	t.Helper()
	resolvedPath, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("failed to load test file path: %w", err))
	}

	return resolvedPath
}

func TestMain(m *testing.M) {
	m.Run()
}

/*
	NewFileSystem
*/
func Test_NewFileSystem_CreateOSFileSystemTypeSucceeds(t *testing.T) {
	testingfileSystem, err := io.NewFileSystem(io.OSFileSystemType)
	require.NoError(t, err)
	require.NotNil(t, testingfileSystem)
}

func Test_NewFileSystem_CreateMemFileSystemTypeSucceeds(t *testing.T) {
	testingfileSystem, err := io.NewFileSystem(io.MemFileSystemType)
	require.NoError(t, err)
	require.NotNil(t, testingfileSystem)
}

func Test_NewFileSystem_CreateIncorrectFileSystemTypeFails(t *testing.T) {
	nonExistantFileSystem := "DoesNotExist"
	testingfileSystem, err := io.NewFileSystem(io.FileSystemType(nonExistantFileSystem))
	require.Error(t, err)
	require.Nil(t, testingfileSystem)
}
