package keybindings

import (
	"github.com/jesseduffield/gocui"
)

func InitializeRepository(
	g *gocui.Gui,
	viewName string,
	open func() error,
) {
	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return open()
		},
	)
}
