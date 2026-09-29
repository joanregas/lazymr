package keybindings

import "github.com/jesseduffield/gocui"

func InitializeMergeRequestCreator(
	g *gocui.Gui,
	formView string,
	selectorView string,
	selectorSearchView string,
	editorView string,
	previous func() error,
	next func() error,
	edit func() error,
	toggle func() error,
	create func() error,
	close func() error,
	selectorPrevious func() error,
	selectorNext func() error,
	selectorEnter func() error,
	selectorSearch func() error,
	selectorClose func() error,
	searchEnter func() error,
	searchClose func() error,
	editorEnter func() error,
	editorSave func() error,
	editorClose func() error,
) {

	g.SetKeybinding(
		formView,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			return previous()
		},
	)

	g.SetKeybinding(
		formView,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			return next()
		},
	)

	g.SetKeybinding(
		formView,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return edit()
		},
	)

	g.SetKeybinding(
		formView,
		gocui.NewKeyRune(' '),
		func(g *gocui.Gui, v *gocui.View) error {
			return toggle()
		},
	)

	g.SetKeybinding(
		formView,
		gocui.NewKeyRune('c'),
		func(g *gocui.Gui, v *gocui.View) error {
			return create()
		},
	)

	g.SetKeybinding(
		formView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return close()
		},
	)

	g.SetKeybinding(
		selectorView,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			return selectorPrevious()
		},
	)

	g.SetKeybinding(
		selectorView,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			return selectorNext()
		},
	)

	g.SetKeybinding(
		selectorView,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return selectorEnter()
		},
	)

	g.SetKeybinding(
		selectorView,
		gocui.NewKeyStrMod("f", gocui.ModCtrl),
		func(g *gocui.Gui, v *gocui.View) error {
			return selectorSearch()
		},
	)

	g.SetKeybinding(
		selectorView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return selectorClose()
		},
	)

	g.SetKeybinding(
		selectorSearchView,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return searchEnter()
		},
	)

	g.SetKeybinding(
		selectorSearchView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return searchClose()
		},
	)

	g.SetKeybinding(
		editorView,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return editorEnter()
		},
	)

	g.SetKeybinding(
		editorView,
		gocui.NewKeyStrMod("s", gocui.ModCtrl),
/*		gocui.NewKey(
			gocui.KeyEnter,
			"",
			gocui.ModCtrl,
		),*/
		func(g *gocui.Gui, v *gocui.View) error {
			return editorSave()
		},
	)

	g.SetKeybinding(
		editorView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return editorClose()
		},
	)
}

