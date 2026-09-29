package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func createFooterView(g *gocui.Gui) error {
	v, err := g.SetView(
		footerView,
		0,
		0,
		1,
		1,
		0,
	)

	if v == nil {
		return err
	}

	v.Frame = false

	v.FgColor = colorYellow

	renderFooter(v, repositoriesView)

	return nil
}

func renderFooter(
	v *gocui.View,
	focusedView string,
) {
	v.Clear()

	footer := footerText(focusedView)

	fmt.Fprint(v, footer)
}


func footerText(focusedView string) string {
	switch focusedView {
	case repositoriesView:
		return "Edit config file: e | " +
			"Create MR: c | " +
			"Switch to a recent repo: <enter> | " +
			"Refresh Data: R | " +
			"Keybindings: ?"

	case mergeRequestView:
		return "Close MR: C | " +
			"Merge MR: m | " +
			"Create MR: c | " +
			"Refresh Data: R | " +
			"Keybindings: ?"

	case filesView:
		return "Create MR: c | " +
			"Refresh Data: R | " +
			"Keybindings: ?"

	case pipelinesView:
		return "Open pipeline <enter> | " +
			"Create MR: c | " +
			"Open Browser: b | " +
			"Refresh Data: R | " +
			"Keybindings: ?"

	case statusView:
		return "Create MR: c | " +
			"Refresh Data: R | " +
			"Keybindings: ?"

	case overviewView:
		return "Create MR: c | " +
			"Refresh Data: R | " +
			"Keybindings: ?"

	default:
		return "Create MR: c | " +
			"Refresh Data: R | " +
			"Keybindings: ?"
	}
}

