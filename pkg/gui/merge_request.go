package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

/*TO-DO: Rethink this... because i couldnt make it work with tcell colors */
const (
	mergeRequestColorReset  = "\x1b[0m"
	mergeRequestColorRed    = "\x1b[31m"
	mergeRequestColorGreen  = "\x1b[32m"
	mergeRequestColorYellow = "\x1b[33m"
	mergeRequestColorCyan   = "\x1b[36m"
	mergeRequestColorGrey   = "\x1b[90m"
)

func createMergeRequestView(
	g *gocui.Gui,
	mergeRequests *gitlab.MergeRequests,
) error {
	v, err := createView(
		g,
		mergeRequestView,
		"[2]-Merge Requests",
	)
	if v == nil {
		return err
	}

	renderMergeRequests(
		v,
		mergeRequests,
		false,
		0,
	)

	return nil
}

func renderMergeRequests(
	v *gocui.View,
	mergeRequests *gitlab.MergeRequests,
	focused bool,
	selected int,
) {
	v.Clear()

	if mergeRequests == nil || len(mergeRequests.Items) == 0 {
		v.Footer = "0 of 0"
		v.Highlight = false
		v.SetOriginY(0)

		return
	}

	if selected < 0 {
		selected = 0
	}

	if selected >= len(mergeRequests.Items) {
		selected = len(mergeRequests.Items) - 1
	}

	v.Highlight = focused
	v.SelBgColor = gocui.NewRGBColor(35, 35, 35)
	v.SelFgColor = gocui.ColorWhite

	/*
		Render all merge requests first.

		This is important because OriginY operates on the
		content stored in the view.
	*/
	for index := range mergeRequests.Items {
		mergeRequest := &mergeRequests.Items[index]

		statusColor := mergeRequestStatusColor(
			mergeRequest,
		)

		fmt.Fprintf(
			v,
			"%s!%d%s %s\n",
			statusColor,
			mergeRequest.IID,
			mergeRequestColorReset,
			mergeRequest.Title,
		)
	}

	/*
		The first and last rows are occupied by the frame.
		The footer is rendered on the bottom frame, so it
		does not consume a content row or thats what i belive.
	*/
	_, height := v.Size()

	visibleRows := height - 2

	if visibleRows < 1 {
		visibleRows = 1
	}

	/*
		Keep the selected merge request visible.
	*/
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

	if cursorY < 0 {
		cursorY = 0
	}

	if cursorY >= visibleRows {
		cursorY = visibleRows - 1
	}

	v.SetCursor(0, cursorY)

	v.Footer = fmt.Sprintf(
		"%d of %d",
		selected+1,
		len(mergeRequests.Items),
	)
}

func (gui *GUI) closeSelectedMergeRequest() error {
	if gui.State == nil {
		return nil
	}

	mergeRequests := gui.State.GetMergeRequests()

	if mergeRequests == nil ||
		len(mergeRequests.Items) == 0 {
		return nil
	}

	selectedIndex := gui.State.GetSelectedMergeRequest()

	if selectedIndex < 0 ||
		selectedIndex >= len(mergeRequests.Items) {
		return nil
	}

	selected := &mergeRequests.Items[selectedIndex]

	return gui.mergeRequestConfirmation.open(
		fmt.Sprintf(
			"Close merge request !%d?",
			selected.IID,
		),
		"Are you sure you want to close this merge request?",
		func() error {
			return gui.closeSelectedMergeRequestConfirmed()
		},
	)
}


func (gui *GUI) refreshMergeRequests() error {
	v, err := gui.g.View(mergeRequestView)
	if err != nil {
		return fmt.Errorf(
			"get merge request view: %w",
			err,
		)
	}

	focused := false

	currentView := gui.g.CurrentView()

	if currentView != nil {
		focused = currentView.Name() == mergeRequestView
	}

	renderMergeRequests(
		v,
		gui.State.GetMergeRequests(),
		focused,
		gui.State.GetSelectedMergeRequest(),
	)

	return nil
}

func mergeRequestStatusColor(
	mergeRequest *gitlab.MergeRequest,
) string {
	if mergeRequest == nil {
		return mergeRequestColorGrey
	}

	if mergeRequest.Draft ||
		mergeRequest.DetailedMergeStatus == "draft" {
		return mergeRequestColorCyan
	}

	if mergeRequest.MergeError != "" {
		return mergeRequestColorRed
	}

	if mergeRequest.HasConflicts ||
		!mergeRequest.BlockingDiscussionsResolved {
		return mergeRequestColorRed
	}

	switch mergeRequest.DetailedMergeStatus {
	case "mergeable":
		return mergeRequestColorGreen

	case "conflict":
		return mergeRequestColorYellow

	case "blocked":
		return mergeRequestColorRed

	case "checking":
		return mergeRequestColorYellow

	case "ci_must_pass":
		return mergeRequestColorYellow

	case "ci_still_running":
		return mergeRequestColorYellow

	case "not_approved":
		return mergeRequestColorRed

	default:
		return mergeRequestColorGrey
	}
}


func (gui *GUI) refreshMergeRequestSelection() { 
	v, err := gui.g.View(mergeRequestView) 
	if err != nil || v == nil { 
		return 
	} 
	if gui.State == nil { 
		return 
	} 
	currentView := gui.g.CurrentView() 
	focused := currentView != nil && 
	  currentView.Name() == mergeRequestView 
	renderMergeRequests( v, gui.State.GetMergeRequests(), focused, gui.State.GetSelectedMergeRequest(), ) 
}

func (gui *GUI) selectNextMergeRequest() error {
	if gui.State == nil || gui.State.GetMergeRequests() == nil {
		return nil
	}

	items := gui.State.GetMergeRequests().Items

	if len(items) == 0 {
		return nil
	}

	selected := gui.State.GetSelectedMergeRequest()

	if selected >= len(items)-1 {
		selected = 0
	} else {
		selected++
	}

	return gui.selectMergeRequest(selected)
}

func (gui *GUI) selectPreviousMergeRequest() error {
	if gui.State == nil || gui.State.GetMergeRequests() == nil {
		return nil
	}

	items := gui.State.GetMergeRequests().Items

	if len(items) == 0 {
		return nil
	}

	selected := gui.State.GetSelectedMergeRequest()

	if selected <= 0 {
		selected = len(items) - 1
	} else {
		selected--
	}

	return gui.selectMergeRequest(selected)
}


func (gui *GUI) selectMergeRequest(index int) error {
	if err := gui.State.SelectMergeRequest(
		gui.GitLab,
		index,
	); err != nil {
		return err
	}

	return gui.Refresh()
}



func (gui *GUI) closeSelectedMergeRequestConfirmed() error {
	if gui.State == nil {
		return nil
	}

	commandLogs := gui.State.GetCommandLogs()
	mergeRequests := gui.State.GetMergeRequests()

	if mergeRequests == nil ||
		len(mergeRequests.Items) == 0 {
		return nil
	}

	selectedIndex := gui.State.GetSelectedMergeRequest()

	if selectedIndex < 0 ||
		selectedIndex >= len(mergeRequests.Items) {
		return nil
	}

	selected := &mergeRequests.Items[selectedIndex]

	commandLogs.Add(
		fmt.Sprintf(
			"Closing merge request !%d...",
			selected.IID,
		),
		false,
	)

	if err := selected.Close(
		gui.GitLab,
		mergeRequests.Project,
	); err != nil {
		commandLogs.Add(
			fmt.Sprintf(
				"Failed to close merge request !%d: %s",
				selected.IID,
				err,
			),
			true,
		)

		_ = gui.refreshCommandLogs()

		return nil
	}

	commandLogs.Add(
		fmt.Sprintf(
			"Merge request !%d closed successfully",
			selected.IID,
		),
		false,
	)

	if err := mergeRequests.Refresh(
		gui.GitLab,
	); err != nil {
		commandLogs.Add(
			fmt.Sprintf(
				"Failed to refresh merge requests after closing !%d: %s",
				selected.IID,
				err,
			),
			true,
		)

		_ = gui.refreshCommandLogs()

		return nil
	}

	if len(mergeRequests.Items) == 0 {
		gui.State.Files = nil
		gui.State.Pipelines = nil
		gui.State.SetSelectedMergeRequest(0)

		return gui.Refresh()
	}

	if selectedIndex >= len(mergeRequests.Items) {
		selectedIndex = len(mergeRequests.Items) - 1
	}

	/* do i really need this? */ 
	if gui.selectMergeRequest != nil {
		return gui.selectMergeRequest(selectedIndex)
	}

	return gui.Refresh()
}


func (gui *GUI) mergeSelectedMergeRequest() error {
	if gui.State == nil {
		return nil
	}

	mergeRequests := gui.State.GetMergeRequests()

	if mergeRequests == nil ||
		len(mergeRequests.Items) == 0 {
		return nil
	}

	selectedIndex := gui.State.GetSelectedMergeRequest()

	if selectedIndex < 0 ||
		selectedIndex >= len(mergeRequests.Items) {
		return nil
	}

	selected := &mergeRequests.Items[selectedIndex]

	return gui.mergeRequestConfirmation.open(
		fmt.Sprintf(
			"Merge merge request !%d?",
			selected.IID,
		),
		"Are you sure you want to merge this merge request?",
		func() error {
			return gui.mergeSelectedMergeRequestConfirmed()
		},
	)
}

func (gui *GUI) mergeSelectedMergeRequestConfirmed() error {
	if gui.State == nil {
		return nil
	}

	commandLogs := gui.State.GetCommandLogs()
	mergeRequests := gui.State.GetMergeRequests()

	if mergeRequests == nil ||
		len(mergeRequests.Items) == 0 {
		return nil
	}

	selectedIndex := gui.State.GetSelectedMergeRequest()

	if selectedIndex < 0 ||
		selectedIndex >= len(mergeRequests.Items) {
		return nil
	}

	selected := &mergeRequests.Items[selectedIndex]

	commandLogs.Add(
		fmt.Sprintf(
			"Merging merge request !%d...",
			selected.IID,
		),
		false,
	)

	if err := selected.Merge(
		gui.GitLab,
		mergeRequests.Project,
	); err != nil {
		commandLogs.Add(
			fmt.Sprintf(
				"Failed to merge merge request !%d: %s",
				selected.IID,
				err,
			),
			true,
		)

		_ = gui.refreshCommandLogs()

		return nil
	}

	commandLogs.Add(
		fmt.Sprintf(
			"Merge request !%d merged successfully",
			selected.IID,
		),
		false,
	)

	if err := mergeRequests.Refresh(
		gui.GitLab,
	); err != nil {
		commandLogs.Add(
			fmt.Sprintf(
				"Failed to refresh merge requests after merging !%d: %s",
				selected.IID,
				err,
			),
			true,
		)

		_ = gui.refreshCommandLogs()

		return nil
	}

	if len(mergeRequests.Items) == 0 {
		gui.State.Files = nil
		gui.State.Pipelines = nil
		gui.State.SetSelectedMergeRequest(0)

		return gui.Refresh()
	}

	if selectedIndex >= len(mergeRequests.Items) {
		selectedIndex = len(mergeRequests.Items) - 1
	}

	return gui.selectMergeRequest(selectedIndex)
}

