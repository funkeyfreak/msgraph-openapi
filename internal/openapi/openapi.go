package openapi

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/funkeyfreak/msgraph-openapi/internal/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

const (
	IsExternalRefsAllowed = true
)

var (
	ErrSwaggerNotLoaded error = errors.New("swagger definition is not loaded")
)

type Swagger *openapi3.Swagger

type Operation openapi3.Operation

func (o *Operation) ContainsTags(tags ...string) bool {
	return utils.Contains(o.Tags, tags...)
}

type OpenApi struct {
	swagger *openapi3.Swagger
	loader  *openapi3.SwaggerLoader
}

func NewOpenApi() *OpenApi {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = IsExternalRefsAllowed
	return &OpenApi{
		loader: loader,
	}
}

// LoadSwaggerFromUri loads an OpenApi spec from a remote URL
func (o *OpenApi) LoadSwaggerFromUri(location *url.URL) (Swagger, error) {
	swagger, err := o.loader.LoadSwaggerFromURI(location)
	if err != nil {
		return nil, err
	}
	o.swagger = swagger

	return swagger, nil
}

// LoadSwaggerFromData loads an OpenApi spec from a byte array
func (o *OpenApi) LoadSwaggerFromData(data []byte) (Swagger, error) {
	swagger, err := o.loader.LoadSwaggerFromData(data)
	if err != nil {
		return nil, err
	}
	o.swagger = swagger

	return swagger, nil
}

// LoadSwaggerFromFile loads an OpenApi spec from local file
func (o *OpenApi) LoadSwaggerFromFile(path string) (Swagger, error) {
	swagger, err := o.loader.LoadSwaggerFromFile(path)
	if err != nil {
		return nil, err
	}
	o.swagger = swagger

	return swagger, nil
}

func (o *OpenApi) UnmarshalJSON(data []byte) error {
	return o.swagger.UnmarshalJSON(data)
}

func (o *OpenApi) MarshalJSON() ([]byte, error) {
	return o.swagger.MarshalJSON()
}

func (o *OpenApi) MarshalYAML() ([]byte, error) {
	swag, err := o.swagger.MarshalJSON()
	if err != nil {
		return nil, err
	}
	tmp := map[string]interface{}{}

	err = json.Unmarshal(swag, &tmp)
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(tmp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *OpenApi) UnmarshalYAML(data []byte) error {
	return yaml.Unmarshal(data, o.swagger)
}
