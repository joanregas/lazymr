package keybindings

import (
	"github.com/jesseduffield/gocui"
)

func InitializePipelines(
	g *gocui.Gui,
	viewName string,
	previous func(),
	next func(),
	open func() error,
	openURL func() error,
) {
	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			previous()
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			next()
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return open()
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyRune('b'),
		func(g *gocui.Gui, v *gocui.View) error {
			return openURL()
		},
	)
}

