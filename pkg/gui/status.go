package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

func createStatusView(
	g *gocui.Gui,
	mergeRequests *gitlab.MergeRequests,
) error {
	v, err := createView(
		g,
		statusView,
		"[5]-Status",
	)
	if v == nil {
		return err
	}

	renderStatus(v, mergeRequests)
	return nil
}

func renderStatus(
	v *gocui.View,
	mergeRequests *gitlab.MergeRequests,
) {
	v.Clear()

	if mergeRequests == nil || len(mergeRequests.Items) == 0 {
		return
	}

	mergeRequest := &mergeRequests.Items[0]

	if mergeRequest.Draft {
		fmt.Fprintln(v, "Draft")
		return
	}

	if mergeRequest.MergeError != "" {
		fmt.Fprintln(v, mergeRequest.MergeError)
		return
	}

	if mergeRequest.HasConflicts {
		fmt.Fprintln(v, "Cannot merge")
		fmt.Fprintln(v, "Conflicts")
		return
	}

	if !mergeRequest.BlockingDiscussionsResolved {
		fmt.Fprintln(v, "Cannot merge")
		fmt.Fprintln(v, "Unresolved discussions")
		return
	}

	switch mergeRequest.DetailedMergeStatus {
	case "mergeable":
		fmt.Fprintln(v, "Ready to merge")

	case "conflict":
		fmt.Fprintln(v, "Cannot merge")
		fmt.Fprintln(v, "Conflicts")

	case "blocked":
		fmt.Fprintln(v, "Cannot merge")

	case "checking":
		fmt.Fprintln(v, "Checking")

	case "ci_must_pass":
		fmt.Fprintln(v, "Waiting for pipeline")

	case "ci_still_running":
		fmt.Fprintln(v, "Pipeline running")

	case "not_approved":
		fmt.Fprintln(v, "Waiting for approval")

	case "draft":
		fmt.Fprintln(v, "Draft")

	default:
		if mergeRequest.DetailedMergeStatus != "" {
			fmt.Fprintln(v, mergeRequest.DetailedMergeStatus)
		} else {
			fmt.Fprintln(v, "Unknown")
		}
	}
}


func (gui *GUI) refreshStatus() error {
	v, err := gui.g.View(statusView)
	if err != nil {
		return fmt.Errorf("get status view: %w",err,)
	}
	renderStatus(v,gui.State.GetMergeRequests(),)
	return nil
}
