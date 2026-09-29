package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

const helpView = "help"

func (gui *GUI) openHelp() error {
	if _, err := gui.g.View(helpView); err == nil {
		return nil
	}

	currentView := gui.g.CurrentView()

	if currentView != nil { 
		gui.helpPreviousView = currentView.Name() 
	}

	v, err := gui.g.SetView(
		helpView,
		0,
		0,
		1,
		1,
		0,
	)

	if v == nil {
		return err
	}

	v.FrameRunes = []rune{
		'─',
		'│',
		'╭',
		'╮',
		'╰',
		'╯',
	}

	v.FrameColor = colorBlue
	v.Title = "Keybindings"
	v.TitleColor = colorBlue
	v.Wrap = true
	v.Autoscroll = false

	renderHelp(v)

	gui.setHelpDimensions()

	_, err = gui.g.SetCurrentView(helpView)

	return err
}

func (gui *GUI) closeHelp() error {
	if _, err := gui.g.View(helpView); err != nil {
		return nil
	}

	if err := gui.g.DeleteView(helpView); err != nil {
		return err
	}

	//gui.restoreFocusedView()
	previousView := gui.helpPreviousView 
	gui.helpPreviousView = "" 
	if previousView != "" { 
		gui.setFocusByName(previousView) 
	} else { 
		gui.setFocus(focusRepositories) 
	}

	return nil
}

func (gui *GUI) restoreFocusedView() {
	currentView := gui.g.CurrentView()

	if currentView == nil {
		gui.setFocus(focusRepositories)
		return
	}

	if currentView.Name() == helpView {
		gui.setFocus(focusRepositories)
		return
	}

	gui.setFocusByName(currentView.Name())
}

func renderHelp(v *gocui.View) {
	v.Clear()

	fmt.Fprintln(v, "Global")
	fmt.Fprintln(v, "  ?        Show this help")
	fmt.Fprintln(v, "  <tab>    Next window")
	fmt.Fprintln(v, "  1-6      Focus window")
	fmt.Fprintln(v, "  c        Create merge request")
	fmt.Fprintln(v, "  @        Toggle command logs")
	fmt.Fprintln(v, "  q / esc  Quit")
	fmt.Fprintln(v, "  R        Refresh Data")
	fmt.Fprintln(v)

	fmt.Fprintln(v, "Repositories")
	fmt.Fprintln(v, "  <enter>  Switch to a recent repository")
	fmt.Fprintln(v, "  e        Edit config file")
	fmt.Fprintln(v)

	fmt.Fprintln(v, "Merge Requests")
	fmt.Fprintln(v, "  ↑ / ↓    Select merge request")
	fmt.Fprintln(v, "  C        Close merge request")
	fmt.Fprintln(v, "  m        Merge merge request")
	fmt.Fprintln(v)

	fmt.Fprintln(v, "Files")
	fmt.Fprintln(v, "  ↑ / ↓    Select file")
	fmt.Fprintln(v, "  <enter>  Expand / collapse directory")
	fmt.Fprintln(v)

	fmt.Fprintln(v, "Pipelines")
	fmt.Fprintln(v, "  ↑ / ↓    Select pipeline")
	fmt.Fprintln(v, "  b        Open Pipeline in Browser")
	fmt.Fprintln(v, "  <enter>  Open pipeline")
	fmt.Fprintln(v)

	fmt.Fprintln(v, "Overview")
	fmt.Fprintln(v, "  ↑ / ↓    Scroll")
	fmt.Fprintln(v, "  pgup     Page up")
	fmt.Fprintln(v, "  pgdn     Page down")
	fmt.Fprintln(v, "  home     Go to beginning")
	fmt.Fprintln(v, "  end      Go to end")
	fmt.Fprintln(v)

	fmt.Fprintln(v, "Press <esc> to close")
}


func (gui *GUI) setHelpDimensions() {
	width, height := gui.g.Size()

	modalWidth := 62
	modalHeight := 30

	if modalWidth > width {
		modalWidth = width
	}

	if modalHeight > height {
		modalHeight = height
	}

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2

	_, err := gui.g.SetView(
		helpView,
		x0,
		y0,
		x0+modalWidth,
		y0+modalHeight,
		0,
	)

	if err != nil && err != gocui.ErrUnknownView {
		return
	}
}

