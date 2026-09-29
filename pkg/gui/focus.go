package gui

import (
	"github.com/jesseduffield/gocui"
)

const (
	focusRepositories = iota
	focusMergeRequests
	focusFiles
	focusPipelines
	focusStatus
	focusOverview
)

var focusableViews = []string{
	repositoriesView,
	mergeRequestView,
	filesView,
	pipelinesView,
	statusView,
	overviewView,
}

func (gui *GUI) focusHandler(
	viewName string,
) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		gui.setFocusByName(viewName)
		return nil
	}
}

func (gui *GUI) nextFocus(
	g *gocui.Gui,
	v *gocui.View,
) error {
	currentView := g.CurrentView()

	if currentView == nil {
		gui.setFocus(focusRepositories)
		return nil
	}

	currentIndex := gui.focusIndex(currentView.Name())

	if currentIndex == -1 {
		gui.setFocus(focusRepositories)
		return nil
	}

	nextIndex := (currentIndex + 1) % len(focusableViews)

	gui.setFocus(nextIndex)

	return nil
}

func (gui *GUI) setFocus(index int) {
	if index < 0 || index >= len(focusableViews) {
		return
	}

	gui.setFocusByName(focusableViews[index])
}

func (gui *GUI) setFocusByName(viewName string) {
	if err := gui.focusView(
		viewName,
		focusableViews,
	); err != nil {
		return
	}

	footer, err := gui.g.View(footerView)
	if err == nil {
		renderFooter(footer, viewName)
	}

	_ = gui.Refresh()
}

func (gui *GUI) focusView(
	viewName string,
	views []string,
) error {
	_, err := gui.g.SetCurrentView(viewName)
	if err != nil {
		return err
	}

	for _, name := range views {
		view, err := gui.g.View(name)
		if err != nil || view == nil {
			continue
		}

		if name == viewName {
			view.FrameColor = gocui.ColorWhite
			view.TitleColor = gocui.ColorWhite
		} else {
			view.FrameColor = colorBlue
			view.TitleColor = colorBlue
		}
	}

	return nil
}

func (gui *GUI) focusIndex(viewName string) int {
	for index, name := range focusableViews {
		if name == viewName {
			return index
		}
	}

	return -1
}

