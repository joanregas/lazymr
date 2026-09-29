package keybindings

import (
	"github.com/jesseduffield/gocui"
)

func InitializeMergeRequest(
	g *gocui.Gui,
	viewName string,
	close func() error,
  merge func() error,
	next func() error,
	previous func() error,
) {
	g.SetKeybinding(
		viewName,
		gocui.NewKeyRune('C'),
		func(g *gocui.Gui, v *gocui.View) error {
			return close()
		},
	)
	g.SetKeybinding( 
		viewName, 
		gocui.NewKeyRune('m'), 
		func(g *gocui.Gui, v *gocui.View) error { 
			return merge() 
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

