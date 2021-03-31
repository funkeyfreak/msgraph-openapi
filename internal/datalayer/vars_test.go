package datalayer_test

import (
	"container/list"
	"path/filepath"

	"github.com/funkeyfreak/msgraph-openapi/internal/datalayer"
	"github.com/funkeyfreak/msgraph-openapi/internal/io"
)

var (
	TestFilePathYaml, _  = filepath.Abs("./testdata/testfile.yaml")
	TestFilePathYaml1, _ = filepath.Abs("./testdata/testfile1.yaml")
	TestFilePathYaml2, _ = filepath.Abs("./testdata/testfile2.yaml")
	TestFilePathJson, _  = filepath.Abs("./testdata/testfile.json")
	TestFilePathJson1, _ = filepath.Abs("./testdata/testfile1.json")
	TestFilePathJson2, _ = filepath.Abs("./testdata/testfile2.json")
	InvalidFilePath      = "/does/not/exist"
	ListDataModel        = &list.List{}
	ByteArrayDataModel   = make([]byte, 0)

	InvalidDataModelType datalayer.DataModelType = "beeps_and_boots"

	filePaths = []string{
		TestFilePathYaml,
		TestFilePathYaml1,
		TestFilePathYaml2,
		//TestFilePathJson,
		//TestFilePathJson1
		//TestFilePathJson2
	}

	// NOTE: since json file type is not implemented, we will uncomment once we have implemented that data structure
	osFileSystem, _ = io.NewFileSystem(io.OSFileSystemType)
	yamlFile, _     = datalayer.NewFile(TestFilePathYaml, osFileSystem)
	yamlFile1, _    = datalayer.NewFile(TestFilePathYaml1, osFileSystem)
	yamlFile2, _    = datalayer.NewFile(TestFilePathYaml2, osFileSystem)
	//jsonFile, _ := datalayer.NewFile(TestFilePathJson, osFileSystem)
	//jsonFile1, _ := datalayer.NewFile(TestFilePathJson1, osFileSystem)
	//jsonFile2, _ := datalayer.NewFile(TestFilePathJson2, osFileSystem)

	files = []datalayer.File{
		yamlFile,
		yamlFile1,
		yamlFile2,
		//jsonFile,
		//jsonFile1,
		//jsonFile2,
	}
)
