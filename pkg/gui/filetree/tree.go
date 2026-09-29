package filetree

import (
	"sort"
	"strings"

	"lazymr/pkg/gitlab"
)

type CollapsedPaths struct {
	paths map[string]bool
}

func NewCollapsedPaths() *CollapsedPaths {
	return &CollapsedPaths{
		paths: make(map[string]bool),
	}
}

func (c *CollapsedPaths) IsCollapsed(path string) bool {
	if c == nil {
		return false
	}

	return c.paths[path]
}

func (c *CollapsedPaths) Toggle(path string) {
	if c == nil || path == "" {
		return
	}

	if c.paths[path] {
		delete(c.paths, path)
		return
	}

	c.paths[path] = true
}

func (c *CollapsedPaths) Collapse(path string) {
	if c == nil || path == "" {
		return
	}

	c.paths[path] = true
}

func (c *CollapsedPaths) Expand(path string) {
	if c == nil {
		return
	}

	delete(c.paths, path)
}

func Build(files []gitlab.File) *Node {
	root := &Node{}

	for i := range files {
		file := &files[i]

		path := strings.Trim(file.NewPath, "/")
		if path == "" {
			continue
		}

		parts := strings.Split(path, "/")

		current := root
		currentPath := ""

		for i, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath += "/" + part
			}

			child := findChild(current, part)

			if child == nil {
				child = &Node{
					Name: part,
					Path: currentPath,
				}

				current.Children = append(
					current.Children,
					child,
				)
			}

			current = child

			if i == len(parts)-1 {
				current.File = file
			}
		}
	}

	sortNodes(root)

	return root
}

func findChild(parent *Node, name string) *Node {
	for _, child := range parent.Children {
		if child.Name == name {
			return child
		}
	}

	return nil
}

func sortNodes(node *Node) {
	sort.Slice(
		node.Children,
		func(i, j int) bool {
			left := node.Children[i]
			right := node.Children[j]

			leftIsDirectory := left.IsDirectory()
			rightIsDirectory := right.IsDirectory()

			if leftIsDirectory != rightIsDirectory {
				return leftIsDirectory
			}

			return left.Name < right.Name
		},
	)

	for _, child := range node.Children {
		sortNodes(child)
	}
}

func VisibleNodes(
	root *Node,
	collapsedPaths *CollapsedPaths,
) []*Node {
	if root == nil {
		return nil
	}

	result := make([]*Node, 0)

	var walk func(*Node)

	walk = func(node *Node) {
		if node != root {
			result = append(result, node)
		}

		if node.IsFile() {
			return
		}

		if collapsedPaths != nil &&
			collapsedPaths.IsCollapsed(node.Path) {
			return
		}

		for _, child := range node.Children {
			walk(child)
		}
	}

	walk(root)

	return result
}

func Files(root *Node) []*Node {
	if root == nil {
		return nil
	}

	result := make([]*Node, 0)

	var walk func(*Node)

	walk = func(node *Node) {
		if node.File != nil {
			result = append(result, node)
			return
		}

		for _, child := range node.Children {
			walk(child)
		}
	}

	walk(root)

	return result
}
