package gui

import (
	"fmt"
)

const loadingView = "loading"

func (gui *GUI) startLoadingScreen() error {
	currentView := gui.g.CurrentView()

	if currentView != nil {
		gui.loadingPreviousView = currentView.Name()
	}

	if gui.Popup.Exists(loadingView) {
		gui.Popup.Close(loadingView)
	}

	width, height := gui.g.Size()

	modalWidth := 50
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	modalHeight := 6
	if modalHeight > height-2 {
		modalHeight = height - 2
	}

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2
	x1 := x0 + modalWidth
	y1 := y0 + modalHeight

	v, err := gui.Popup.Create(
		loadingView,
		"Please wait",
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf("create loading popup: %w", err)
	}

	v.Clear()

	fmt.Fprintln(v)
	fmt.Fprintln(v, "  Refreshing GitLab data...")
	fmt.Fprintln(v)
	fmt.Fprintln(v, "  Please wait...")

	v.Highlight = false
	v.FrameColor = colorWhite
	v.TitleColor = colorWhite

	_, err = gui.g.SetViewOnTop(loadingView)
	if err != nil {
		return fmt.Errorf("raise loading popup: %w", err)
	}

	return gui.Popup.Focus(loadingView)
}

func (gui *GUI) stopLoadingScreen() {
	gui.Popup.Close(loadingView)

	if gui.loadingPreviousView != "" {
		_ = gui.focusView(
			gui.loadingPreviousView,
			focusableViews,
		)

		gui.loadingPreviousView = ""
	}
}
