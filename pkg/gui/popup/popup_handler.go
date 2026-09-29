package popup

import "github.com/jesseduffield/gocui"

type PopupHandler struct {
	gui *gocui.Gui
}

func NewPopupHandler(gui *gocui.Gui) *PopupHandler {
	return &PopupHandler{
		gui: gui,
	}
}

func (p *PopupHandler) View(name string) *gocui.View {
	v, err := p.gui.View(name)
	if err != nil {
		return nil
	}

	return v
}

func (p *PopupHandler) Exists(name string) bool {
	_, err := p.gui.View(name)

	return err == nil
}

func (p *PopupHandler) Create(
	name string,
	title string,
	x0 int,
	y0 int,
	x1 int,
	y1 int,
) (*gocui.View, error) {
	v, err := p.gui.SetView(
		name,
		x0,
		y0,
		x1,
		y1,
		0,
	)

	if v == nil {
		return nil, err
	}

	v.FrameRunes = []rune{
		'─',
		'│',
		'╭',
		'╮',
		'╰',
		'╯',
	}

	v.FrameColor = gocui.ColorBlue
	v.Title = title
	v.TitleColor = gocui.ColorBlue

	v.Highlight = true
	v.SelBgColor = gocui.NewRGBColor(35, 35, 35)
	v.SelFgColor = gocui.ColorWhite

	return v, nil
}

func (p *PopupHandler) Focus(name string) error {
	_, err := p.gui.SetCurrentView(name)

	return err
}

func (p *PopupHandler) Close(name string) {
	_ = p.gui.DeleteView(name)
}

