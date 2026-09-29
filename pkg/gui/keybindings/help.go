package keybindings

import "github.com/jesseduffield/gocui"

func InitializeHelp(
	g *gocui.Gui,
	helpView string,
	closeHelp func() error,
) {
	g.SetKeybinding(
		helpView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return closeHelp()
		},
	)

	g.SetKeybinding(
		helpView,
		gocui.NewKeyRune('?'),
		func(g *gocui.Gui, v *gocui.View) error {
			return closeHelp()
		},
	)

	g.SetKeybinding(
		helpView,
		gocui.NewKeyRune('q'),
		func(g *gocui.Gui, v *gocui.View) error {
			return closeHelp()
		},
	)
}

