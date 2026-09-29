package keybindings

import "github.com/jesseduffield/gocui"

func InitializeMergeRequestConfirmation(
	g *gocui.Gui,
	confirmationView string,
	confirmAction func() error,
	cancelAction func() error,
) {
	g.SetKeybinding(
		confirmationView,
		gocui.NewKeyRune('y'),
		func(g *gocui.Gui, v *gocui.View) error {
			return confirmAction()
		},
	)

	g.SetKeybinding(
		confirmationView,
		gocui.NewKeyRune('n'),
		func(g *gocui.Gui, v *gocui.View) error {
			return cancelAction()
		},
	)

	g.SetKeybinding(
		confirmationView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return cancelAction()
		},
	)
}

