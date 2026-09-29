package gui

import (
	"fmt"
	"strings"
	"time"

	"lazymr/pkg/debug"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

const (
	jobLogModalView = "job_log_modal"
)

type jobLogModal struct {
	job     gitlab.PipelineJobs
	content string
}

func (gui *GUI) openJobLogModal(
	projectPath string,
	job gitlab.PipelineJobs,
) error {
	if gui.GitLab == nil ||
		gui.State == nil ||
		gui.State.Pipelines == nil {
		return fmt.Errorf("GitLab state is not initialized")
	}

	if projectPath == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	if gui.Popup.Exists(jobLogModalView) {
		return nil
	}

	content, err := gui.GitLab.GetJobTrace(
		projectPath,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf(
			"load job %q logs: %w",
			job.Name,
			err,
		)
	}


	width, height := gui.g.Size()

	modalWidth := width - 4

	if modalWidth < 70 {
		modalWidth = width - 2
	}

	modalHeight := height - 4

	if modalHeight < 10 {
		modalHeight = height - 2
	}

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2
	x1 := x0 + modalWidth
	y1 := y0 + modalHeight

	v, err := gui.Popup.Create(
		jobLogModalView,
		fmt.Sprintf(
			"[Job: %s #%d]",
			job.Name,
			job.ID,
		),
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf(
			"create job log popup: %w",
			err,
		)
	}

	v.Wrap = true
	v.Autoscroll = false
	v.Highlight = false

	gui.jobLogModal = &jobLogModal{
		job:     job,
		content: content,
	}

debug.Log(
	"Opening job log: ID=%d name=%q status=%q running=%t finishedAt=%v",
	job.ID,
	job.Name,
	job.Status,
	job.IsRunning(),
	job.FinishedAt,
)


	v.Clear()
	//v.SetOriginY(0)

	fmt.Fprint(v, content)

 /* _, height = v.Size() 
	bufferLines := len(v.BufferLines()) 

	if height > 0 && bufferLines > height { 
		v.SetOriginY(bufferLines - height) 
	}*/

	if err := gui.Popup.Focus(jobLogModalView); err != nil {
		gui.jobLogModal = nil
		gui.Popup.Close(jobLogModalView)

		return fmt.Errorf(
			"focus job log popup: %w",
			err,
		)
	}
  if job.IsRunning() { 
		gui.startJobLogRefresh(projectPath, job) 
	}

	return nil
}



const jobLogScrollStep = 1

func (gui *GUI) scrollJobLog(delta int) {
	v, err := gui.g.View(jobLogModalView)
	if err != nil || v == nil {
		return
	}

	originY := v.OriginY()

	gui.setJobLogOrigin(
		v,
		originY+delta,
	)
}

func (gui *GUI) jobLogPageSize() int {
	v, err := gui.g.View(jobLogModalView)
	if err != nil || v == nil {
		return 1
	}

	_, height := v.Size()

	if height <= 1 {
		return 1
	}

	return height - 1
}

func (gui *GUI) setJobLogOrigin(
	v *gocui.View,
	originY int,
) {
	maxOriginY := gui.jobLogMaxOrigin(v)

	if originY < 0 {
		originY = 0
	}

	if originY > maxOriginY {
		originY = maxOriginY
	}

	v.SetOriginY(originY)
}



func (gui *GUI) jobLogMaxOrigin(
	v *gocui.View,
) int {
	_, height := v.InnerSize()

	if height <= 0 {
		return 0
	}

	lineCount := v.ViewLinesHeight()

	maxOriginY := lineCount - height

	if maxOriginY < 0 {
		return 0
	}

	return maxOriginY
}

func (gui *GUI) jobLogLineCount() int {
	if gui.jobLogModal == nil ||
		gui.jobLogModal.content == "" {
		return 0
	}

	return strings.Count(
		gui.jobLogModal.content,
		"\n",
	) + 1
}

func (gui *GUI) jobLogHome() {
	v, err := gui.g.View(jobLogModalView)
	if err != nil || v == nil {
		return
	}

	v.SetOriginY(0)
	v.SetCursor(0, 0)
}



func (gui *GUI) jobLogEnd() {
	v, err := gui.g.View(jobLogModalView)
	if err != nil || v == nil {
		return
	}

	maxOriginY := gui.jobLogMaxOrigin(v)

	v.SetOriginY(maxOriginY)

	_, height := v.Size()

	if height <= 0 {
		return
	}

	v.SetCursor(0, height-1)
}


func (gui *GUI) closeJobLogModal() {
	gui.Popup.Close(jobLogModalView)

	gui.jobLogModal = nil

	if gui.Popup.Exists(pipelineModalContentView) {
		_ = gui.Popup.Focus(pipelineModalContentView)
		return
	}

	gui.setFocus(focusPipelines)
}



func (gui *GUI) startJobLogRefresh(
	projectPath string,
	job gitlab.PipelineJobs,
) {
	go func() {
		for {
			time.Sleep(3 * time.Second)

			if !gui.refreshJobLog(
				projectPath,
				job,
			) {
				return
			}
		}
	}()
}

func (gui *GUI) refreshJobLog(
	projectPath string,
	job gitlab.PipelineJobs,
) bool {
	latestJob, err := gui.GitLab.GetJob(
		projectPath,
		job.ID,
	)
	if err != nil {
		debug.Log(
			"Failed to refresh job %d status: %v",
			job.ID,
			err,
		)

		return true
	}

	content, err := gui.GitLab.GetJobTrace(
		projectPath,
		job.ID,
	)
	if err != nil {
		debug.Log(
			"Failed to refresh job %d trace: %v",
			job.ID,
			err,
		)

		return true
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if gui.jobLogModal == nil {
			return nil
		}

		if content == gui.jobLogModal.content {
			return nil
		}

		gui.jobLogModal.content = content

		v, err := g.View(jobLogModalView)
		if err != nil || v == nil {
			return nil
		}

		v.Clear()
		fmt.Fprint(v, content)

		if latestJob.IsRunning() {
			gui.setJobLogOrigin(
				v,
				gui.jobLogMaxOrigin(v),
			)
		}

		return nil
	})

	return latestJob.IsRunning()
}

