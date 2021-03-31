package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// Parses an incomming OpenApi document
type Parser interface {
	Read(path string)
}

// Create a new instance of Parser from a local file
func NewParserFromFile(path string) (Parser, error) {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(path)
	if err != nil {
		return nil, ReadFromFileError(path)
	}
	var doc Parser = &openApiDoc{swagger: swagger}

	return doc, nil
}

func (o *openApiDoc) Read(path string) {
	fmt.Println(path)
}

type openApiDoc struct {
	swagger *openapi3.Swagger
}
