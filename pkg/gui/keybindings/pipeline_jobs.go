package keybindings

import (
	"github.com/jesseduffield/gocui"
)

func InitializePipelineJobs(
	g *gocui.Gui,
	viewName string,
	previous func(),
	next func(),
	nextTab func(),
	open func() error,
	openURL func() error,
	play func() error,
	retry func() error,
	refresh func() error,
	close func() error,
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
		gocui.NewKeyName(gocui.KeyTab),
		func(g *gocui.Gui, v *gocui.View) error {
			nextTab()
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

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return close()
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyRune('r'),
		func(g *gocui.Gui, v *gocui.View) error {
			return retry()
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyRune('R'),
		func(g *gocui.Gui, v *gocui.View) error {
			return refresh()
		},
	)

	g.SetKeybinding( 
		viewName, 
		gocui.NewKeyRune('p'), 
		func(g *gocui.Gui, v *gocui.View) error { 
			return play() 
		},
	)
}

