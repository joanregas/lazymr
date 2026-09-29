package gui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/jesseduffield/gocui"
)

/* A poor attempt to use Tcell colors */

const (
	//colorWhite  = gocui.Attribute(tcell.ColorWhite)
	colorYellow = gocui.Attribute(tcell.ColorYellow)
	colorCyan   = gocui.Attribute(tcell.ColorAqua)
	colorRed    = gocui.Attribute(tcell.ColorRed)
	colorGreen  = gocui.Attribute(tcell.ColorGreen)
//	colorBlue   = gocui.Attribute(tcell.ColorDeepSkyBlue)
)

var ( 
 colorBlue = gocui.NewRGBColor(0, 0, 220)
 colorWhite = gocui.NewRGBColor(255, 255, 255)
)

