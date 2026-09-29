package keybindings

import "github.com/jesseduffield/gocui"

func InitializePipelineJobLogs(
	g *gocui.Gui,
	view string,
	scroll func(int) error,
	pageSize func() int,
	home func(),
	end func(),
	close func(*gocui.Gui, *gocui.View) error,
) {
	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			return scroll(-1)
		},
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			return scroll(1)
		},
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyPgup),
		func(g *gocui.Gui, v *gocui.View) error {
			return scroll(-pageSize())
		},
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyPgdn),
		func(g *gocui.Gui, v *gocui.View) error {
			return scroll(pageSize())
		},
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyHome),
		func(g *gocui.Gui, v *gocui.View) error {
			home()
			return nil
		},
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyEnd),
		func(g *gocui.Gui, v *gocui.View) error {
			end()
			return nil
		},
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyRune('q'),
		close,
	)

	g.SetKeybinding(
		view,
		gocui.NewKeyName(gocui.KeyEsc),
		close,
	)
}

