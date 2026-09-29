package keybindings

import "github.com/jesseduffield/gocui"

func InitializeCommandLogs(
	g *gocui.Gui,
	toggleCommandLogs func() error,
) {
	g.SetKeybinding(
		"",
		gocui.NewKeyRune('@'),
		func(g *gocui.Gui, v *gocui.View) error {
			return toggleCommandLogs()
		},
	)
}
