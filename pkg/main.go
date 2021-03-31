package main

import (
	"fmt"
	"io/ioutil"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	PathOperationChar = "<>"
)

var (
	msgraph_v1_yaml = "/Users/dalinwilliams/workdir/github/msgraph-openapi/api/ref/v1/openapi.yaml"
)

/*type PathMethod string

func (pm *PathMethod) FetchHttpMethod(pathItem *openapi3.PathItem) *openapi3.Operation {
	method := strings.Split(pm, PathOperationChar)
	return pathItem.GetOperation(method)
}*/

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ByTags() {

}

func contains(haystack []string, needles ...string) bool {
	if len(needles) == 1 {
		for _, val := range haystack {
			if val == needles[0] {
				return true
			}
		}
	} else {
		m := make(map[string]bool, len(needles))
		for _, s := range needles {
			m[s] = false
		}

		for _, val := range haystack {
			if _, ok := m[val]; !ok {
				delete(m, val)
			}
		}

		if len(m) == 0 {
			return true
		}
	}

	return false
}

type Operation openapi3.Operation

func (o *Operation) ContainsTags(tags ...string) bool {
	return contains(o.Tags, tags...)
}

func byTags() {
	yamlFile, err := ioutil.ReadFile(msgraph_v1_yaml)
	check(err)
	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromData(yamlFile)
	check(err)

	tagsToCheck := []string{"deviceAppManagement.Actions"}

	fmt.Printf("%d paths are available\n", len(swagger.Paths))
	for path, pathItem := range swagger.Paths {
		fmt.Printf("%d methods are available in path %v\n", len(pathItem.Operations()), path)
		delete(swagger.Paths, path)
		for method, operation := range pathItem.Operations() {

			op := Operation(*operation)
			if !op.ContainsTags(tagsToCheck...) {
				fmt.Printf("removing path %v:%v \n", path, method)
				pathItem.SetOperation(method, nil)
			}
		}
		if len(pathItem.Operations()) != 0 {
			swagger.Paths[path] = pathItem
		}
		fmt.Printf("%d methods remain in path %v\n", len(pathItem.Operations()), path)
	}

	fmt.Printf("%d paths are remaining\n", len(swagger.Paths))
}

func main() {
	yamlFile, err := ioutil.ReadFile(msgraph_v1_yaml)
	check(err)
	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromData(yamlFile)
	check(err)

	item := swagger.Paths.Find("/deviceAppManagement")
	fmt.Printf("%+v", item)
}
