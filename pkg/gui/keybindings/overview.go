package keybindings

import "github.com/jesseduffield/gocui"

func InitializeOverview(
	g *gocui.Gui,
	viewName string,
	scroll func(int),
	pageSize func() int,
	home func(),
	end func(),
) {
	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			scroll(-1)
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			scroll(1)
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyPgup),
		func(g *gocui.Gui, v *gocui.View) error {
			scroll(-pageSize())
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyPgdn),
		func(g *gocui.Gui, v *gocui.View) error {
			scroll(pageSize())
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyHome),
		func(g *gocui.Gui, v *gocui.View) error {
			home()
			return nil
		},
	)

	g.SetKeybinding(
		viewName,
		gocui.NewKeyName(gocui.KeyEnd),
		func(g *gocui.Gui, v *gocui.View) error {
			end()
			return nil
		},
	)
}

