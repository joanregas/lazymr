package keybindings

import (
	"github.com/jesseduffield/gocui"
)

func InitializeFiles(
	g *gocui.Gui,
	viewName string,
	toggleTreeNode func() error,
	next func() error,
	previous func() error,
) {
	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return toggleTreeNode()
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			return next()
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			return previous()
		},
	)
}

