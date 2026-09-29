package filetree

import "lazymr/pkg/gitlab"

type Node struct {
	Name     string
	Path     string
	File     *gitlab.File
	Children []*Node
}

func (n *Node) IsFile() bool {
	return n != nil && n.File != nil
}

func (n *Node) IsDirectory() bool {
	return n != nil && n.File == nil
}
