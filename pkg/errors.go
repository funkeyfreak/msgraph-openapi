package main

import "fmt"

type ReadFromFileError string

func (p ReadFromFileError) Error() string {
	return fmt.Sprintf("parser: failed to read file %s", string(p))
}
