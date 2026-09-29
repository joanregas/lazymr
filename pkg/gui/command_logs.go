package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

func createCommandLogsView(g *gocui.Gui) error {
	v, err := createView(
		g,
		commandLogsView,
		"[7]-Command Logs",
	)

	if v == nil {
		return err
	}
	v.Wrap = true

	renderCommandLogs(v, nil)

	return nil
}

func renderCommandLogs(
	v *gocui.View,
	logs *gitlab.CommandLogs,
) {
	v.Clear()

	if logs == nil || len(logs.Items) == 0 {
		renderCommandLogsWelcome(v)
		return
	}

	for _, entry := range logs.Items {
		fmt.Fprintf(
			v,
			"%s %s\n",
			entry.Time.Format("15:04:05"),
			entry.Message,
		)
	}
}

func renderCommandLogsWelcome(v *gocui.View) {
	fmt.Fprintln(
		v,
		"You can hide/focus this panel by pressing '@'",
	)
	fmt.Fprintln(v)
	fmt.Fprintln(
		v,
		"Waiting for GitLab operations...",
	)
}

func (gui *GUI) refreshCommandLogs() error {
	v, err := gui.g.View(commandLogsView)
	if err != nil {
		return fmt.Errorf(
			"get command logs view: %w",
			err,
		)
	}

	renderCommandLogs(
		v,
		gui.State.GetCommandLogs(),
	)

	return nil
}

func (gui *GUI) toggleCommandLogs() error { 
	gui.commandLogsVisible = !gui.commandLogsVisible 
	return gui.layout(gui.g) 
}

