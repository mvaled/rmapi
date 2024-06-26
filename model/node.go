package model

import (
	"errors"
	"time"
)

type Node struct {
	Document *Document
	Children map[string]*Node
	Parent   *Node
}

func CreateNode(document Document) Node {
	return Node{&document, make(map[string]*Node, 0), nil}
}

func (node *Node) Name() string {
	return node.Document.VisibleName
}

func (node *Node) Id() string {
	return node.Document.ID
}
func (node *Node) Version() int {
	return node.Document.Version
}

func (node *Node) IsRoot() bool {
	return node.Id() == ""
}

func (node *Node) IsDirectory() bool {
	return node.Document.Type == "CollectionType"
}

func (node *Node) IsFile() bool {
	return !node.IsDirectory()
}

func (node *Node) EntyExists(id string) bool {
	_, ok := node.Children[id]
	return ok
}

func (node *Node) LastModified() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, node.Document.ModifiedClient)
}

func (node *Node) FindByName(name string) (*Node, error) {
	for _, n := range node.Children {
		if n.Name() == name {
			return n, nil
		}
	}
	return nil, errors.New("entry doesn't exist")
}
