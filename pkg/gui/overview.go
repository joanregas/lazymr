package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

const (
	diffReset  = "\x1b[0m"
	diffRed    = "\x1b[31m"
	diffGreen  = "\x1b[32m"
	diffCyan   = "\x1b[36m"
	diffYellow = "\x1b[33m"
)

type OverviewData struct {
	MergeRequest *gitlab.MergeRequest
	File         *gitlab.File
}

func createOverviewView(g *gocui.Gui) error {
	v, err := createView(
		g,
		overviewView,
		"[6]-Overview",
	)
	if v == nil {
		return err
	}

	v.Wrap = true

	renderOverview(v, OverviewData{})

	return nil
}


func renderOverview(
	v *gocui.View,
	data OverviewData,
) {
	v.Clear()

	switch {
	case data.MergeRequest != nil:
		renderMergeRequestOverview(
			v,
			data.MergeRequest,
		)

	case data.File != nil:
		renderFileOverview(
			v,
			data.File,
		)

	default:
		renderOverviewSplash(v)
	}
}

func renderMergeRequestOverview(
	v *gocui.View,
	mergeRequest *gitlab.MergeRequest,
) {
	fmt.Fprintf(
		v,
		"%s!%d%s %s\n\n",
		diffCyan,
		mergeRequest.IID,
		diffReset,
		mergeRequest.Title,
	)

	renderOverviewField(
		v,
		"State",
		mergeRequest.State,
	)

	renderOverviewField(
		v,
		"Author",
		mergeRequest.Author,
	)

	renderOverviewField(
		v,
		"Source",
		mergeRequest.SourceBranch,
	)

	renderOverviewField(
		v,
		"Target",
		mergeRequest.TargetBranch,
	)

	renderOverviewField(
		v,
		"Draft",
		fmt.Sprintf("%t", mergeRequest.Draft),
	)

	renderOverviewField(
		v,
		"Merge status",
		mergeRequest.DetailedMergeStatus,
	)

	renderOverviewField(
		v,
		"Conflicts",
		fmt.Sprintf("%t", mergeRequest.HasConflicts),
	)

	renderOverviewField(
		v,
		"Discussions resolved",
		fmt.Sprintf("%t", mergeRequest.BlockingDiscussionsResolved),
	)

	renderOverviewField(
		v,
		"Merge when pipeline succeeds",
		fmt.Sprintf("%t", mergeRequest.MergeWhenPipelineSucceeds),
	)

	if mergeRequest.MergeError != "" {
		renderOverviewField(
			v,
			"Merge error",
			mergeRequest.MergeError,
		)
	}

	renderOverviewField(
		v,
		"Assignees",
		formatStringList(mergeRequest.Assignees),
	)

	renderOverviewField(
		v,
		"Reviewers",
		formatStringList(mergeRequest.Reviewers),
	)

	renderOverviewField(
		v,
		"Labels",
		formatStringList(mergeRequest.Labels),
	)

	renderOverviewField(
		v,
		"Created",
		mergeRequest.CreatedAt,
	)

	renderOverviewField(
		v,
		"Updated",
		mergeRequest.UpdatedAt,
	)

	renderOverviewField(
		v,
		"Notes",
		fmt.Sprintf("%d", mergeRequest.UserNotesCount),
	)

	if mergeRequest.Description != "" {
		fmt.Fprintln(v)
		fmt.Fprintf(
			v,
			"%sDescription%s\n",
			diffCyan,
			diffReset,
		)
		fmt.Fprintln(v)

		for _, line := range strings.Split(
			mergeRequest.Description,
			"\n",
		) {
			fmt.Fprintln(v, line)
		}
	}
}

func renderOverviewField(
	v *gocui.View,
	label string,
	value string,
) {
	if value == "" {
		value = "-"
	}

	fmt.Fprintf(
		v,
		"%-28s %s\n",
		label+":",
		value,
	)
}

func formatStringList(values []string) string {
	if len(values) == 0 {
		return "-"
	}

	return strings.Join(values, ", ")
}

func renderFileOverview(
	v *gocui.View,
	file *gitlab.File,
) {
	fmt.Fprintf(
		v,
		"%s%s%s\n\n",
		diffCyan,
		file.NewPath,
		diffReset,
	)

	renderDiff(v, file.Diff)
}

func renderOverviewSplash(v *gocui.View) {
	fmt.Fprintln(v)
	fmt.Fprintln(v)
	fmt.Fprintln(v, " _                     __  __ ____")
	fmt.Fprintln(v, "| |    __ _ _____   _ |  \\/  |  _ \\")
	fmt.Fprintln(v, "| |   / _` |_  / | | || |\\/| | |_) |")
	fmt.Fprintln(v, "| |__| (_| |/ /| |_| || |  | |  _ <")
	fmt.Fprintln(v, "|_____\\__,_/___|\\__, ||_|  |_|_| \\_\\")
	fmt.Fprintln(v, "                 __/ |")
	fmt.Fprintln(v, "                |___/")
	fmt.Fprintln(v)
	fmt.Fprintln(v, "GitLab Merge Request TUI")
	fmt.Fprintln(v)
	fmt.Fprintln(v, "Copyright 2026 lazymr")
	fmt.Fprintln(v)
	fmt.Fprintln(v, "A terminal UI for working with GitLab Merge Requests.")
	fmt.Fprintln(v)
	fmt.Fprintln(v, "Select a merge request or file to view its details.")
	fmt.Fprintln(v)
}

func renderDiff(
	v *gocui.View,
	diff string,
) {
	if diff == "" {
		fmt.Fprintln(v, "No diff available.")
		return
	}

	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+++") ||
			strings.HasPrefix(line, "---"):
			fmt.Fprintf(
				v,
				"%s%s%s\n",
				diffYellow,
				line,
				diffReset,
			)

		case strings.HasPrefix(line, "@@"):
			fmt.Fprintf(
				v,
				"%s%s%s\n",
				diffCyan,
				line,
				diffReset,
			)

		case strings.HasPrefix(line, "+"):
			fmt.Fprintf(
				v,
				"%s%s%s\n",
				diffGreen,
				line,
				diffReset,
			)

		case strings.HasPrefix(line, "-"):
			fmt.Fprintf(
				v,
				"%s%s%s\n",
				diffRed,
				line,
				diffReset,
			)

		default:
			fmt.Fprintln(v, line)
		}
	}
}

func (gui *GUI) refreshOverview() error {
	v, err := gui.g.View(overviewView)
	if err != nil {
		return fmt.Errorf(
			"get overview view: %w",
			err,
		)
	}

	data := OverviewData{}

	if gui.State != nil {
		data.File = gui.State.GetSelectedFile()

		if data.File == nil {
			mergeRequests := gui.State.GetMergeRequests()

			if mergeRequests != nil {
				selectedIndex := gui.State.GetSelectedMergeRequest()

				if selectedIndex >= 0 &&
					selectedIndex < len(mergeRequests.Items) {
					data.MergeRequest =
						&mergeRequests.Items[selectedIndex]
				}
			}
		}
	}

	renderOverview(v, data)

	return nil
}


const overviewScrollStep = 1

func (gui *GUI) scrollOverview(delta int) {
	v, err := gui.g.View(overviewView)
	if err != nil || v == nil {
		return
	}

	originY := v.OriginY()
	originY += delta

	if originY < 0 {
		originY = 0
	}

	_, height := v.Size()

	if height <= 0 {
		return
	}

	lines := strings.Count(
		v.Buffer(),
		"\n",
	)

	maxOriginY := lines - height + 1

	if maxOriginY < 0 {
		maxOriginY = 0
	}

	if originY > maxOriginY {
		originY = maxOriginY
	}

	v.SetOriginY(originY)
}

func (gui *GUI) overviewPageSize() int {
	v, err := gui.g.View(overviewView)
	if err != nil || v == nil {
		return 1
	}

	_, height := v.Size()

	if height <= 1 {
		return 1
	}

	return height - 1
}

func (gui *GUI) overviewHome() {
	v, err := gui.g.View(overviewView)
	if err != nil || v == nil {
		return
	}

	v.SetOriginY(0)
}

func (gui *GUI) overviewEnd() {
	v, err := gui.g.View(overviewView)
	if err != nil || v == nil {
		return
	}

	_, height := v.Size()

	if height <= 0 {
		return
	}

	lines := strings.Count(
		v.Buffer(),
		"\n",
	)

	maxOriginY := lines - height + 1

	if maxOriginY < 0 {
		maxOriginY = 0
	}

	v.SetOriginY(maxOriginY)
}

