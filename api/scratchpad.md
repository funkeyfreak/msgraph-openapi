yq reference
    * yq eval '[.. | select(has("$ref"))]' ./api/ref/v1/openapi.yaml

# Tool Behaviour
* All tools should have the ability to be ran in the background using routines

## File Creation
* All files will be created in a temporary space until command has completed
    - consider the default temp directories on linux, then expand to other operating systems. Handle this via configuration set a the tools installation

# Commands

## Batch
Tool should enable commands to be ran in searilization or concurrent "modes"

E.g. serilized mode can look like the following:
``` sh
    <tool-root-cmd> batch openapi.[yaml|json], ...
        filter <filter-params>
        dedupe <dedupe-params>
        decompose <decompose-params>

```
The behaviour should be as follows: all files will have the filter, then dedupe, then decompose commands ran against them 1) in the order that the files are given to the command and 2) in the order of the commands themselves. 

E.g. serilized mode can look like the following:
``` sh
    <tool-root-cmd> -c batch openapi.[yaml|json], ...
        filter <filter-params>
        dedupe <dedupe-params>
        decompose <decompose-params>
```
The behaviour should be as follows: all files are ran in concurrent processes, with the filter, then dedupe, the decompose commands ran against them. The commands will be ran in the order they are given. Note that the "-c" parameter tells our command to run each file concurrently.

command syntax: batch [flags] <commnads...> file
    batch
        [-s|--stateful: share the state! NOTE: should behave different if '-c|--concurrent' is/is not provided. If -c is provided, then the files being processed must wait to complete. NOTE: May be dropped. NOTE: The output state is shared, so this can effectively merge many yaml/json files into one]
        [-c|--concurrent: run each file concurrently! Will consume more memory]

## Filter
Tool should create a sub-set of the document according to some open-api query/yq query. I may need to re-implement several parts of the YQ library, as it is a bit buggy in some regards

The tool will need to return the matched items AND the references required. This tool will not re-arrange any refs - a user can leverage the 'Decompose' tool after running split on a large yaml file

## Compose
Tool should be able to build a large yaml from a doc which contains eiter external or internal ref references

Default behavior: Only compose/synthesize external references - do not synthesize internal references. Output will default to refrence file format

command syntax:
    compose
        [-c/--compose: Also compose/synthesize the internal references to items in comopose NOTE WILL DELETE ITEMS FROM WITHIN COMPOSE]
        [-d/--delete: Delete the original files]
        [-f/--format: The format of the created document - default is the reference file format - e.g. if the reference file is json, then the output will be json]
        openapi.[yaml|json]: the root openapi file on which to operate

## Dedupe
Tool should be able to deduplicate the following openapi doc items:
* paths
    + verb [get, post, put, etc]
        - parameters
            = in
        - responses
* components
    

By default, tool should place deduped items in the components object 
    under a representative path, e.g. "components/parameters", "components/responses", "components/parameters/in", etc

## Decompose
Tool should be able to decompose on the following openapi 3 specific structures and export them into a files in a given folder structure:
* paths
    + verb [get, post, put, etc]
        - parameters
        - responses

By default, the tool should decompose "cleverly," finding the largest chunks of duplicated text, then exporting these out to a given file

command syntax:
    decompose 
        [-p/--path: The path under which to save the generated file(s) - will default to "./gen" direcotry]
        [-A/--all: Decompose on every possible openapi element]
        [-p/--paths: Decompose the paths only]
        []
        openapi.yaml: The yaml file on which to operate
