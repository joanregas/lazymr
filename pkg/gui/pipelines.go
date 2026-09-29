package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

type pipelineSelector struct {
	selected int
}

/* TO-DO: Again... im not good with this tcell colors..... fix this so we can have everything in colors.go T.T */
const (
	pipelineColorReset  = "\x1b[0m"
	pipelineColorGreen  = "\x1b[32m"
	pipelineColorRed    = "\x1b[31m"
	pipelineColorYellow = "\x1b[33m"
	pipelineColorGrey   = "\x1b[90m"
	pipelineColorCyan   = "\x1b[36m"
)

/*TO-DO:  That should go to icons.go */
const (
	pipelineStatusSuccess = "✓"
	pipelineStatusFailed  = "✗"
	pipelineStatusRunning = "●"
	pipelineStatusPending = "◷"
	pipelineStatusCreated = "○"
	pipelineStatusCanceled = "⏸"
	pipelineStatusSkipped  = "⊘"
	pipelineStatusManual   = "▶"
	pipelineStatusUnknown  = "?"
)

func createPipelinesView(
	g *gocui.Gui,
	pipelines *gitlab.Pipelines,
) error {
	v, err := createView(
		g,
		pipelinesView,
		"[4]-Pipelines",
	)
	if v == nil {
		return err
	}

	renderPipelines(
		v,
		pipelines,
		false,
		0,
	)

	return nil
}

func renderPipelines(
	v *gocui.View,
	pipelines *gitlab.Pipelines,
	focused bool,
	selected int,
) {
	v.Clear()

	if pipelines == nil || len(pipelines.Items) == 0 {
		v.Footer = "0 of 0"
		v.Highlight = false
		v.SetOriginY(0)

		return
	}

	if selected < 0 {
		selected = 0
	}

	if selected >= len(pipelines.Items) {
		selected = len(pipelines.Items) - 1
	}

	v.Highlight = focused
	v.SelBgColor = gocui.NewRGBColor(35, 35, 35)
	v.SelFgColor = gocui.ColorWhite

	/*
		Render all pipeline lines first.
		This can cause problems for some users if they have lots of pipelines...
		we will see 
	*/
	for _, pipeline := range pipelines.Items {
		icon, color := pipelineStatusIcon(pipeline.Status)

		fmt.Fprintf(
			v,
			"%s%s%s #%d %s %s\n",
			color,
			icon,
			pipelineColorReset,
			pipeline.IID,
			pipeline.Status,
			pipeline.Ref,
		)
	}

	_, height := v.Size()

	visibleRows := height - 2

	if visibleRows < 1 {
		visibleRows = 1
	}

	originY := v.OriginY()

	if selected < originY {
		originY = selected
	}

	if selected >= originY+visibleRows {
		originY = selected - visibleRows + 1
	}

	if originY < 0 {
		originY = 0
	}

	v.SetOriginY(originY)

	cursorY := selected - originY

	if cursorY < 0 { cursorY = 0 } 
	if cursorY >= visibleRows { cursorY = visibleRows - 1 } 
	v.SetCursor(0, cursorY)

	v.Footer = fmt.Sprintf(
		"%d of %d",
		selected+1,
		len(pipelines.Items),
	)
}

func (gui *GUI) pipelineNext() {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		len(gui.State.Pipelines.Items) == 0 {
		return
	}

	if gui.pipelineSelector == nil {
		gui.pipelineSelector = &pipelineSelector{}
	}

	if gui.pipelineSelector.selected >=
		len(gui.State.Pipelines.Items)-1 {
		return
	}

	gui.pipelineSelector.selected++

	gui.renderPipelines()
}

func (gui *GUI) pipelinePrevious() {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		len(gui.State.Pipelines.Items) == 0 {
		return
	}

	if gui.pipelineSelector == nil {
		gui.pipelineSelector = &pipelineSelector{}
	}

	if gui.pipelineSelector.selected <= 0 {
		return
	}

	gui.pipelineSelector.selected--

	gui.renderPipelines()
}

func (gui *GUI) renderPipelines() {
	v, err := gui.g.View(pipelinesView)
	if err != nil || v == nil {
		return
	}

	if gui.State == nil {
		return
	}

	selected := 0

	if gui.pipelineSelector != nil {
		selected = gui.pipelineSelector.selected
	}

	focused := false

	currentView := gui.g.CurrentView()

	if currentView != nil {
		focused = currentView.Name() == pipelinesView
	}

	renderPipelines(
		v,
		gui.State.Pipelines,
		focused,
		selected,
	)
}

func (gui *GUI) refreshPipelines() error {
	v, err := gui.g.View(pipelinesView)
	if err != nil {
		return fmt.Errorf(
			"get pipelines view: %w",
			err,
		)
	}

	if gui.State == nil {
		renderPipelines(
			v,
			nil,
			false,
			0,
		)

		return nil
	}

	pipelines := gui.State.GetPipelines()

	selected := 0

	if gui.pipelineSelector != nil {
		selected = gui.pipelineSelector.selected
	}

	if pipelines != nil && len(pipelines.Items) > 0 {
		if selected < 0 {
			selected = 0
		}

		if selected >= len(pipelines.Items) {
			selected = len(pipelines.Items) - 1
		}

		if gui.pipelineSelector == nil {
			gui.pipelineSelector = &pipelineSelector{}
		}

		gui.pipelineSelector.selected = selected
	} else {
		selected = 0

		if gui.pipelineSelector != nil {
			gui.pipelineSelector.selected = 0
		}
	}

	focused := false

	currentView := gui.g.CurrentView()

	if currentView != nil {
		focused = currentView.Name() == pipelinesView
	}

	renderPipelines(
		v,
		pipelines,
		focused,
		selected,
	)

	return nil
}

func pipelineStatusIcon(status string) (string, string) {
	switch status {
	case "success":
		return pipelineStatusSuccess, pipelineColorGreen

	case "failed":
		return pipelineStatusFailed, pipelineColorRed

	case "running":
		return pipelineStatusRunning, pipelineColorYellow

	case "pending",
		"waiting_for_resource",
		"preparing",
		"scheduled":
		return pipelineStatusPending, pipelineColorYellow

	case "created":
		return pipelineStatusCreated, pipelineColorGrey

	case "canceled":
		return pipelineStatusCanceled, pipelineColorGrey

	case "skipped":
		return pipelineStatusSkipped, pipelineColorGrey

	case "manual":
		return pipelineStatusManual, pipelineColorCyan

	default:
		return pipelineStatusUnknown, pipelineColorGrey
	}
}

func (gui *GUI) openSelectedPipelineBrowserURL() error {
	if gui.State == nil || gui.State.Pipelines == nil {
		return nil
	}

	selected := gui.pipelineSelector.selected

	if selected < 0 || selected >= len(gui.State.Pipelines.Items) {
		return nil
	}

	pipeline := gui.State.Pipelines.Items[selected]

	if pipeline.WebURL == "" {
		return fmt.Errorf(
			"pipeline #%d has no web URL",
			pipeline.IID,
		)
	}

	if err := gui.OpenUrl(pipeline.WebURL); err != nil {
		return fmt.Errorf(
			"open pipeline #%d: %w",
			pipeline.IID,
			err,
		)
	}

	return nil
}

