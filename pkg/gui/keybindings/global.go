package keybindings

import (
	"github.com/jesseduffield/gocui"
)

func InitializeGlobal(
	g *gocui.Gui,
	focusableViews []string,
	focusHandler func(string) func(*gocui.Gui, *gocui.View) error,
	nextFocus func(*gocui.Gui, *gocui.View) error,
	createMergeRequest func() error,
	openHelp func() error,
	quit func() error,
	refresh func() error,
) {
	for index, viewName := range focusableViews {
		key := rune('1' + index)

		g.SetKeybinding(
			"",
			gocui.NewKeyRune(key),
			focusHandler(viewName),
		)
	}

	g.SetKeybinding(
		"",
		gocui.NewKeyName(gocui.KeyTab),
		nextFocus,
	)

	g.SetKeybinding(
		"",
		gocui.NewKeyRune('c'),
		func(g *gocui.Gui, v *gocui.View) error {
			return createMergeRequest()
		},
	)

	g.SetKeybinding(
		"",
		gocui.NewKeyRune('q'),
		func(g *gocui.Gui, v *gocui.View) error {
			return quit()
		},
	)

	g.SetKeybinding(
		"",
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return quit()
		},
	)
	g.SetKeybinding(
		"", 
		gocui.NewKeyRune('?'), 
		func(g *gocui.Gui, v *gocui.View) error { 
			return openHelp() 
		}, 
	)
  g.SetKeybinding(
  	"",
  	gocui.NewKeyRune('R'),
  	func(g *gocui.Gui, v *gocui.View) error {
    	return refresh()
	  },
  )
}

