package filter

import (
	"container/list"
	"fmt"

	"github.com/funkeyfreak/msgraph-openapi/internal/openapi"
)

type Filter struct {
	Swagger openapi.Swagger
	Flags   FilterFlag
}

func NewFilter(swagger openapi.Swagger, filterJobs list.List) *Filter {
	return &Filter{}
}

func (f *Filter) ByPath() error {
	f.Swagger.Paths.Find()
}

func (f *Filter) ByTag() error {
	for path, pathItem := range f.Swagger.Paths {
		fmt.Printf("%d methods are available in path %v\n", len(pathItem.Operations()), path)
		delete(f.Swagger.Paths, path)
		for method, operation := range pathItem.Operations() {
			op := openapi.Operation(*operation)
			if !op.ContainsTags(f.Flags.FilterData...) {
				fmt.Printf("removing path %v:%v \n", path, method)
				pathItem.SetOperation(method, nil)
			}
		}
		if len(pathItem.Operations()) != 0 {
			f.Swagger.Paths[path] = pathItem
		}
		fmt.Printf("%d methods remain in path %v\n", len(pathItem.Operations()), path)
	}

	return nil
}

func (f *Filter) byPath(paths []string) {

}

// filter by openapi reference
func FilterJob(in <-chan Filter, out chan<- Filter, errors chan<- error) {
	for job := range in {
		var err error
		switch job.Flags.FilterType {
		case FilterFlagPath:

		case FilterFlagTag:
			err = job.ByTag()
		}

		if err != nil {
			errors <- err
		}

		out <- job
	}
}
