package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
	"lazymr/pkg/gui/filetree"
)

const (
	fileTreeExpandedIcon  = "▼"
	fileTreeCollapsedIcon = "▶"
	fileTreeDirectoryIcon = "\uf07b" // 
	fileTreeFileIcon      = "\uf15b" // 
)

var filesCollapsedPaths = filetree.NewCollapsedPaths()

func createFilesView(
	g *gocui.Gui,
	files *gitlab.Files,
) error {
	v, err := createView(
		g,
		filesView,
		"[3]-Files",
	)
	if v == nil {
		return err
	}

	renderFiles(
		v,
		files,
		false,
		"",
	)

	return nil
}

func renderFiles(
	v *gocui.View,
	files *gitlab.Files,
	focused bool,
	selectedPath string,
) {
	v.Clear()

	if files == nil {
		return
	}

	v.Highlight = focused
	v.SelBgColor = gocui.NewRGBColor(35, 35, 35)
	v.SelFgColor = gocui.ColorWhite

	tree := filetree.Build(files.Items)

	nodes := filetree.VisibleNodes(
		tree,
		filesCollapsedPaths,
	)

	selectedLine := -1

	for index, node := range nodes {
		if node.Path == selectedPath {
			selectedLine = index
			break
		}
	}

	for _, node := range nodes {
		renderFileNode(
			v,
			node,
		)
	}

	if !focused || selectedLine < 0 {
		return
	}

	scrollFilesView(
		v,
		selectedLine,
	)

	cursorY := selectedLine - v.OriginY()

	if cursorY < 0 {
		cursorY = 0
	}

	v.SetCursor(0, cursorY)
}

func renderFileNode(
	v *gocui.View,
	node *filetree.Node,
) {
	depth := treeDepth(node)

	indentation := strings.Repeat(
		"  ",
		depth,
	)

	if node.IsDirectory() {
		icon := fileTreeExpandedIcon

		if filesCollapsedPaths.IsCollapsed(node.Path) {
			icon = fileTreeCollapsedIcon
		}

		fmt.Fprintf(
			v,
			"%s%s %s %s\n",
			indentation,
			icon,
			fileTreeDirectoryIcon,
			node.Name,
		)

		return
	}

	fmt.Fprintf(
		v,
		"%s%s %s\n",
		indentation,
		fileTreeFileIcon,
		node.Name,
	)
}

func treeDepth(node *filetree.Node) int {
	if node == nil || node.Path == "" {
		return 0
	}

	return strings.Count(
		strings.Trim(node.Path, "/"),
		"/",
	)
}

func scrollFilesView(
	v *gocui.View,
	selectedLine int,
) {
	_, height := v.Size()

	if height <= 0 {
		return
	}

	contentHeight := height - 2

	if contentHeight <= 0 {
		contentHeight = 1
	}

	originY := v.OriginY()

	if selectedLine < originY {
		originY = selectedLine
	}

	if selectedLine >= originY+contentHeight {
		originY = selectedLine - contentHeight + 1
	}

	if originY < 0 {
		originY = 0
	}

	v.SetOriginY(originY)
}


func (gui *GUI) refreshFiles() error {
	v, err := gui.g.View(filesView)
	if err != nil {
		return fmt.Errorf(
			"get files view: %w",
			err,
		)
	}

	focused := false

	currentView := gui.g.CurrentView()

	if currentView != nil {
		focused = currentView.Name() == filesView
	}

	if focused {
		gui.initializeFileSelection()
	}

	renderFiles(
		v,
		gui.State.GetFiles(),
		focused,
		gui.State.GetSelectedFilePath(),
	)

	return nil
}



func (gui *GUI) selectPreviousFile() error {
	if gui.State == nil || gui.State.GetFiles() == nil {
		return nil
	}

	tree := filetree.Build(gui.State.GetFiles().Items)
	nodes := filetree.VisibleNodes(tree, filesCollapsedPaths)

	if len(nodes) == 0 {
		return nil
	}

	currentPath := gui.State.GetSelectedFilePath()

	currentIndex := -1

	for index, node := range nodes {
		if node.Path == currentPath {
			currentIndex = index
			break
		}
	}

	if currentIndex == -1 {
		file := firstVisibleFile(nodes)

		if file == nil {
			return nil
		}
		//gui.State.SetSelectedFilePath(nodes[0].Path)
		gui.State.SetSelectedFilePath(file.Path)
		return gui.Refresh()
	}

	if currentIndex == 0 {
		return nil
	}

	gui.State.SetSelectedFilePath(
		nodes[currentIndex-1].Path,
	)

	return gui.Refresh()
}

func (gui *GUI) selectNextFile() error {
	if gui.State == nil || gui.State.GetFiles() == nil {
		return nil
	}

	tree := filetree.Build(gui.State.GetFiles().Items)
	nodes := filetree.VisibleNodes(tree, filesCollapsedPaths)

	if len(nodes) == 0 {
		return nil
	}

	currentPath := gui.State.GetSelectedFilePath()

	currentIndex := -1

	for index, node := range nodes {
		if node.Path == currentPath {
			currentIndex = index
			break
		}
	}

	if currentIndex == -1 {
		gui.State.SetSelectedFilePath(nodes[0].Path)
		return gui.Refresh()
	}

	if currentIndex >= len(nodes)-1 {
		return nil
	}

	gui.State.SetSelectedFilePath(
		nodes[currentIndex+1].Path,
	)

	return gui.Refresh()
}


func (gui *GUI) toggleSelectedFileTreeNode() error {
	if gui.State == nil || gui.State.GetFiles() == nil {
		return nil
	}

	tree := filetree.Build(gui.State.GetFiles().Items)
	nodes := filetree.VisibleNodes(tree, filesCollapsedPaths)

	selectedPath := gui.State.GetSelectedFilePath()

	for _, node := range nodes {
		if node.Path != selectedPath {
			continue
		}

		if node.IsFile() {
			return nil
		}

		filesCollapsedPaths.Toggle(node.Path)

		return gui.Refresh()
	}

	return nil
}


func (gui *GUI) initializeFileSelection() {
	if gui.State == nil || gui.State.GetFiles() == nil {
		return
	}

	if gui.State.GetSelectedFilePath() != "" {
		return
	}

	tree := filetree.Build(gui.State.GetFiles().Items)

	nodes := filetree.VisibleNodes(
		tree,
		filesCollapsedPaths,
	)

	if len(nodes) == 0 {
		return
	}
  file := firstVisibleFile(nodes)

	if file == nil {
		return
	}
	gui.State.SetSelectedFilePath(file.Path)
}

func firstVisibleFile(
	nodes []*filetree.Node,
) *filetree.Node {
	for _, node := range nodes {
		if node.IsFile() {
			return node
		}
	}

	return nil
}

