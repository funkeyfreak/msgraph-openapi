package filequery

type Query interface {
	Execute(string) []byte
}
