package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

/* since im the only one that knows how the hell this is built better throw an error for users so they don't need to read the f** code */
const configurationRequiredMessage = `No GitLab authentication configuration was found.

Lazy MR looked for:

  ~/.config/glab-cli/config.yml
  ~/.config/lazymr/config.yml

A GitLab Personal Access Token is required.

Configure it in or configure glab-cli properly:

  ~/.config/lazymr/config.yml

Required permissions:

  api`

func ShowConfigurationRequired() error {
	g, err := gocui.NewGui(
		gocui.NewGuiOpts{
			OutputMode: gocui.OutputTrue,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"initialize configuration error GUI: %w",
			err,
		)
	}

	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		width := 76
		height := 22

		if width > maxX-2 {
			width = maxX - 2
		}

		if height > maxY-2 {
			height = maxY - 2
		}

		if width < 40 || height < 10 {
			return nil
		}

		x0 := (maxX - width) / 2
		y0 := (maxY - height) / 2

   	v, _ := g.SetView(
			"error",
			x0,
			y0,
			x0+width,
			y0+height,
			0,
		)

		if v == nil {
			return fmt.Errorf(
				"create configuration error view",
			)
		}

		v.FrameRunes = []rune{
			'─',
			'│',
			'╭',
			'╮',
			'╰',
			'╯',
		}

		v.FrameColor = colorRed
		v.Title = " Lazy MR "
		v.TitleColor = gocui.ColorWhite
		v.Wrap = true

		v.Clear()

		fmt.Fprintln(v)
		fmt.Fprintln(
			v,
			"  Lazy MR configuration required",
		)
		fmt.Fprintln(v)

		for _, line := range splitLines(configurationRequiredMessage) {
			fmt.Fprintln(v, "  "+line)
		}

		fmt.Fprintln(v)
		fmt.Fprintln(v, "  [q] Quit")

		return nil
	})

	g.SetKeybinding(
		"",
		gocui.NewKeyRune('q'),
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		},
	)

	g.SetKeybinding(
		"",
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		},
	)

	if err := g.MainLoop(); err != nil &&
		err != gocui.ErrQuit {
		return fmt.Errorf(
			"configuration error GUI: %w",
			err,
		)
	}

	return nil
}

func splitLines(text string) []string {
	var lines []string

	current := ""

	for _, character := range text {
		if character == '\n' {
			lines = append(lines, current)
			current = ""
			continue
		}

		current += string(character)
	}

	lines = append(lines, current)

	return lines
}

