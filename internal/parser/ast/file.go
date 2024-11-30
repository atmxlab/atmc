package ast

type File struct {
	node
	imports []Import
	object  Object
}
