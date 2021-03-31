package main

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

const (
	msgraph_v1_directory_location = "/Users/dalinwilliams/workdir/github/msgraph-openapi/api/gen/"
	msgraph_v1_yaml               = "/Users/dalinwilliams/workdir/github/msgraph-openapi/api/ref/v1/openapi.yaml"
	refgrabber                    = ".. | select(has(\"$ref\"))"
)

var processedNodes *list.List

func check(e error) {
	if e != nil {
		fmt.Printf("error encountered: %v", e)
		panic(e)
		//os.Exit(1)
	}
}
func readDocuments(reader io.Reader, filename string, fileIndex int) (*list.List, error) {
	decoder := yaml.NewDecoder(reader)
	inputList := list.New()
	var currentIndex uint = 0

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			switch reader := reader.(type) {
			case *os.File:
				reader.Close()
			}
			return inputList, nil
		} else if errorReading != nil {
			return nil, errorReading
		}
		candidateNode := &yqlib.CandidateNode{
			Document:         currentIndex,
			Filename:         filename,
			Node:             &dataBucket,
			FileIndex:        fileIndex,
			EvaluateTogether: true,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}
}

func fetchMatchingNodes(expression string, filenames []string) (*list.List, error) {
	fileIndex := 0

	var allDocuments *list.List = list.New()
	for _, filename := range filenames {
		reader, err := os.Open(msgraph_v1_yaml)
		if err != nil {
			return nil, err
		}
		fileDocuments, err := readDocuments(reader, filename, fileIndex)
		if err != nil {
			return nil, err
		}
		allDocuments.PushBackList(fileDocuments)
		fileIndex = fileIndex + 1
	}
	matches, err := yqlib.NewAllAtOnceEvaluator().EvaluateCandidateNodes(expression, allDocuments)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func fetchUniqueRefs() []string {
	res := make([]string, 0)

	// var buf bytes.Buffer = bytes.Buffer{}

	// lo := log.New(&buf, "logger: ", log.Lshortfile)

	//log, err := logging.GetLogger("yq-lib")
	//check(err)

	logging.SetLevel(logging.ERROR, "yq-lib")

	f, err := os.Create("/Users/dalinwilliams/workdir/github/msgraph-openapi/api/ref/test")
	check(err)
	w := bufio.NewWriter(f)
	defer w.Flush()

	//printer := yqlib.NewPrinter(w, false, false, false, 0, false)

	//err = yqlib.NewAllAtOnceEvaluator().EvaluateFiles(refgrabber, []string{msgraph_v1_yaml}, printer)
	//check(err)
	//streamEvaluator := yqlib.NewStreamEvaluator()

	//err = streamEvaluator.EvaluateFiles(refgrabber, []string{msgraph_v1_yaml}, printer)
	//check(err)

	yamlFile, err := ioutil.ReadFile(msgraph_v1_yaml)

	check(err)

	swagger := &openapi3.Swagger{}
	err = yaml.Unmarshal(yamlFile, swagger)
	for _, server := range swagger.Servers {
		fmt.Printf("%+v\n", server.URL)
	}
	check(err)

	/*UNCOMMENT
	list, err := fetchMatchingNodes(refgrabber, []string{msgraph_v1_yaml})
	check(err)

	fmt.Println(list.Len())
	*/

	//printer.PrintResults(list)

	file, err := os.Open(msgraph_v1_yaml)
	check(err)
	// buf := bytes.Buffer{}
	filename := msgraph_v1_yaml
	fileIndex := 0
	inputList := list.New()
	decoder := yaml.NewDecoder(file)
	var currentIndex uint = 0

	for {
		var dataBucket yaml.Node
		errorReading := decoder.Decode(&dataBucket)

		if errorReading == io.EOF {
			//switch reader := reader.(type) {
			//case *os.File:
			//}
			//return inputList, nil
			defer file.Close()
			break
		} else if errorReading != nil {
			//return nil, errorReading
			check(err)
		}
		candidateNode := &yqlib.CandidateNode{
			Document:         currentIndex,
			Filename:         filename,
			Node:             &dataBucket,
			FileIndex:        fileIndex,
			EvaluateTogether: true,
		}

		inputList.PushBack(candidateNode)

		currentIndex = currentIndex + 1
	}

	fmt.Printf("we have read %d nodes\n", inputList.Len())

	/*node := yaml.Node{}
	err = yaml.Unmarshal(yamlFile, &node)
	check(err)*/
	//out, err := yaml.Marshal(node)
	/*for _, value := range node.Content {
		fmt.Printf("%+v \n", value)
		if strings.Contains(value.Tag, "map") {
			for _, val := range value.Content {
				fmt.Printf("%+v \n", val)
			}
		}
	}*/

	//r := yaml.NewDecoder(file).Decode(node)
	//err =  swagger. //node.Decode(swagger) // node.Encode(swagger)
	//check(err)

	evaluator := yqlib.NewAllAtOnceEvaluator()

	//list, err := evaluator.EvaluateNodes(refgrabber, node.Content...)
	//TODO: HAHAHAHHAHHA
	processedNodes = inputList
	list, err := evaluator.EvaluateCandidateNodes(refgrabber, inputList)
	check(err)

	m := make(map[string]bool)

	fmt.Printf("%d <<=== yes\n", list.Len())

	t := 0
	for e := list.Front(); e != nil; e = e.Next() {
		node, ok := e.Value.(*yqlib.CandidateNode)
		if !ok {
			panic("nooooo")
		}

		buf := bytes.Buffer{}
		encoder := yaml.NewEncoder(&buf)
		err = encoder.Encode(node.Node)
		check(err)

		ma := make(map[string]string)
		err = yaml.Unmarshal(buf.Bytes(), ma)
		check(err)

		//fmt.Println(node.Node)
		str := buf.String()

		if _, ok := m[str]; !ok {
			m[str] = true

			var test string = ma["$ref"]
			s := fmt.Sprintf("%v\n", ma["$ref"])
			w.WriteString(s)
			res = append(res, string(test[2:]))

			//fmt.Println(str)
		}
		t++
	}

	/*i := 0
	for val, _ := range m {
		res[i] = val
		i++
	}*/
	return res
}

type Trie struct {
	root *trieNode
}

type trieNode struct {
	val    string
	isFile bool
	m      map[string]*trieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &trieNode{},
	}
}

func (t *Trie) Insert(paths []string) {
	curr := t.root

	for _, val := range paths {
		if _, ok := curr.m[val]; !ok {
			curr.m[val] = &trieNode{val: val}
		}
		curr = curr.m[val]
	}
	curr.isFile = true
}

func openFileWriter(file string) (*os.File, error) {
	return os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

func createEmptyFile(file string) error {
	d := []byte("")

	return ioutil.WriteFile(file, d, 0644)
}

func createFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

func genApiRefFiles(files []string) error {
	for _, val := range files {
		//fmt.Println(val)
		//break
		strArr := strings.Split(val, "/")
		strArrLen := len(strArr)
		strArr = strArr[:strArrLen-1]
		//fmt.Println(msgraph_v1_directory_location + strings.Join(strArr, "/"))
		err := createFolder(msgraph_v1_directory_location + strings.Join(strArr, "/"))
		check(err)

		createEmptyFile(msgraph_v1_directory_location + val + ".yaml")
		check(err)
	}

	return nil
}

func main() {
	err := os.RemoveAll(msgraph_v1_directory_location)
	check(err)

	// decompose the document by refs
	res := fetchUniqueRefs()
	err = genApiRefFiles(res)
	check(err)

	for _, val := range res {
		strArr := strings.Split(val, "/")
		strArrLen := len(strArr)
		queryStr := fmt.Sprintf(".%v.\"%v\"", strings.Join(strArr[:strArrLen-1], "."), strArr[strArrLen-1])
		outFile := msgraph_v1_directory_location + "/" + val + ".yaml"

		f, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		check(err)
		defer f.Close()

		// VERY HACKY - Move this implemntation into a higher level object/struct - use with parser
		//list, err := fetchMatchingNodes(refgrabber, []string{msgraph_v1_yaml})
		evaluator := yqlib.NewAllAtOnceEvaluator()
		//fmt.Println(queryStr)
		list, err := evaluator.EvaluateCandidateNodes(queryStr, processedNodes)
		check(err)

		for e := list.Front(); e != nil; e = e.Next() {
			node, ok := e.Value.(*yqlib.CandidateNode)
			if !ok {
				panic("nooooo")
			}

			buf := bytes.Buffer{}
			encoder := yaml.NewEncoder(&buf)
			err = encoder.Encode(node.Node)
			check(err)

			/*ma := make(map[string]string)
			err = yaml.Unmarshal(buf.Bytes(), ma)
			check(err)*/

			//fmt.Println(node.Node)
			//str := buf.String()
			//f.Write(buf.Bytes())
			_, err := f.WriteString(buf.String())
			check(err)
			//fmt.Println(bLen)
			f.Sync()
			/*if _, ok := m[str]; !ok {
				m[str] = true

				var test string = ma["$ref"]
				s := fmt.Sprintf("%v\n", ma["$ref"])
				w.WriteString(s)
				res = append(res, string(test[2:]))

				//fmt.Println(str)
			}*/
		}

	}

	//

	//swagger := &openapi3.Swagger{}
	yamlFile, err := ioutil.ReadFile(msgraph_v1_yaml)
	check(err)
	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromData(yamlFile)
	check(err)

	//err = yaml.Unmarshal(yamlFile, swagger)
	//check(err)
	pathsDirectoryPath := msgraph_v1_directory_location + "paths"
	err = createFolder(pathsDirectoryPath)
	check(err)
	for path, pathItem := range swagger.Paths {
		//TODO: do not make excessive directories
		pathByteArr := []byte(path)
		// remove first forward slash
		pathByteArr = pathByteArr[1:]
		path = string(pathByteArr)
		// if we have more than one directory, remove the last and create the folder structure
		pathsArr := strings.Split(path, "/")
		if len(pathsArr) > 1 {
			folder := strings.Join(pathsArr[:len(pathsArr)-1], "/")
			err := createFolder(pathsDirectoryPath + "/" + folder)
			check(err)
		}

		// create the output file for the paths yaml
		outFile := pathsDirectoryPath + "/" + path + ".yaml"
		err = createEmptyFile(outFile)
		check(err)

		// TODO: More optimal please! Perhaps find a way to fetch all paths without querying
		// 		for each path individually. Also, move evaluator into seperate class. Finally, fix the way in which we query these paths, this is pretty jank
		// re-query for processedNodes for each path
		/*queryStr := fmt.Sprintf(".path./%v", path)
		fmt.Println(queryStr)
		evaluator := yqlib.NewAllAtOnceEvaluator()
		list, err := evaluator.EvaluateCandidateNodes(queryStr, processedNodes)
		check(err)

		for e := list.Front(); e != nil; e = e.Next() {
			node, ok := e.Value.(*yqlib.CandidateNode)
			if !ok {
				panic("nooooo")
			}

			buf := bytes.Buffer{}
			encoder := yaml.NewEncoder(&buf)
			err = encoder.Encode(node.Node)
			check(err)
			f, err := openFileWriter(outFile)
			check(err)
			defer f.Close()
			bCount, err := f.Write(buf.Bytes())
			check(err)
			fmt.Println(bCount)
			f.Sync()

		}*/

		// old approach - creates invalid yaml
		/*
			byteArr, err := yaml.Marshal(pathItem)
			check(err)
			f, err := openFileWriter(outFile)
			check(err)
			defer f.Close()
			bCount, err := f.Write(byteArr)
			check(err)
			fmt.Println(bCount)
			f.Sync()
		*/

		// Using openapi3 and hacky cast to avoid yaml serilization issues
		jsonBuff, err := pathItem.MarshalJSON()
		check(err)

		f, err := openFileWriter(outFile)
		check(err)
		defer f.Close()

		tmp := map[string]interface{}{}

		err = json.Unmarshal(jsonBuff, &tmp)
		check(err)
		err = yaml.NewEncoder(f).Encode(tmp)
		check(err)
		fi, err := f.Stat()
		check(err)
		fmt.Printf("wrote %d \n", fi.Size())
		f.Sync()

	}
}
