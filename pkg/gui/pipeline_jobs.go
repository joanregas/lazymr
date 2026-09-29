package gui

import (
	"fmt"
	"time"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

const (
	pipelineModalView        = "pipeline_modal"
	pipelineModalHeaderView  = "pipeline_modal_header"
	pipelineModalTabsView    = "pipeline_modal_tabs"
	pipelineModalContentView = "pipeline_modal_content"
	pipelineModalFooterView  = "pipeline_modal_footer"
)

type pipelineModalTab int

const (
	pipelineJobsTab pipelineModalTab = iota
	pipelineDownstreamTab
)

type pipelineModalState struct {
	pipeline    gitlab.Pipeline
	projectPath string
	tab         pipelineModalTab
	selected    int
}

type pipelineModal struct {
	pipeline    *gitlab.Pipeline
	projectPath string
	tab         pipelineModalTab
	selected    int

	history []pipelineModalState
}

func (gui *GUI) openPipelineModal() error {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		len(gui.State.Pipelines.Items) == 0 {
		return nil
	}

	if gui.pipelineSelector == nil {
		gui.pipelineSelector = &pipelineSelector{}
	}

	selected := gui.pipelineSelector.selected

	if selected < 0 ||
		selected >= len(gui.State.Pipelines.Items) {
		return nil
	}

	if gui.Popup.Exists(pipelineModalView) {
		return nil
	}

	pipeline := gui.State.Pipelines.Items[selected]
	projectPath := gui.State.Pipelines.Project

	if err := gui.State.Pipelines.LoadPipelineDetails(
		gui.GitLab,
		projectPath,
		pipeline.ID,
	); err != nil {
		return fmt.Errorf(
			"load pipeline #%d details: %w",
			pipeline.IID,
			err,
		)
	}

	width, height := gui.g.Size()

	modalWidth := width * 4 / 5

	if modalWidth < 70 {
		modalWidth = 70
	}

	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	modalHeight := height * 4 / 5

	if modalHeight < 20 {
		modalHeight = 20
	}

	if modalHeight > height-2 {
		modalHeight = height - 2
	}

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2
	x1 := x0 + modalWidth
	y1 := y0 + modalHeight

	modal, err := gui.Popup.Create(
		pipelineModalView,
		fmt.Sprintf("[Pipeline #%d]", pipeline.IID),
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf(
			"create pipeline popup: %w",
			err,
		)
	}

	modal.Highlight = false

	innerX0 := x0 + 1
	innerX1 := x1 - 1

	headerHeight := 10
	tabsHeight := 3
	footerHeight := 4

	headerY0 := y0 + 1
	headerY1 := headerY0 + headerHeight - 1

	tabsY0 := headerY1 + 1
	tabsY1 := tabsY0 + tabsHeight - 1

	footerY1 := y1 - 1
	footerY0 := footerY1 - footerHeight + 1

	contentY0 := tabsY1 + 1
	contentY1 := footerY0 - 1

	header, err := gui.g.SetView(
		pipelineModalHeaderView,
		innerX0,
		headerY0,
		innerX1,
		headerY1,
		0,
	)
	if header == nil {
		gui.Popup.Close(pipelineModalView)

		return fmt.Errorf(
			"create pipeline modal header: %w",
			err,
		)
	}

	tabs, err := gui.g.SetView(
		pipelineModalTabsView,
		innerX0,
		tabsY0,
		innerX1,
		tabsY1,
		0,
	)
	if tabs == nil {
		gui.Popup.Close(pipelineModalHeaderView)
		gui.Popup.Close(pipelineModalView)

		return fmt.Errorf(
			"create pipeline modal tabs: %w",
			err,
		)
	}

	content, err := gui.g.SetView(
		pipelineModalContentView,
		innerX0,
		contentY0,
		innerX1,
		contentY1,
		0,
	)
	if content == nil {
		gui.Popup.Close(pipelineModalTabsView)
		gui.Popup.Close(pipelineModalHeaderView)
		gui.Popup.Close(pipelineModalView)

		return fmt.Errorf(
			"create pipeline modal content: %w",
			err,
		)
	}

	footer, err := gui.g.SetView(
		pipelineModalFooterView,
		innerX0,
		footerY0,
		innerX1,
		footerY1,
		0,
	)
	if footer == nil {
		gui.Popup.Close(pipelineModalContentView)
		gui.Popup.Close(pipelineModalTabsView)
		gui.Popup.Close(pipelineModalHeaderView)
		gui.Popup.Close(pipelineModalView)

		return fmt.Errorf(
			"create pipeline modal footer: %w",
			err,
		)
	}

	header.Title = "HEADER"
	//tabs.Title = "TABS"
	content.Title = "CONTENT"
	//footer.Title = "FOOTER"

	header.FrameColor = colorBlue 
	header.TitleColor = colorBlue

  tabs.FrameColor = colorBlue 
	content.FrameColor = colorBlue 
	content.TitleColor = colorBlue


	header.Highlight = false
	tabs.Highlight = false

	content.Highlight = true
	content.SelBgColor = gocui.NewRGBColor(35, 35, 35)
	content.SelFgColor = gocui.ColorWhite

	footer.Highlight = false

	gui.pipelineModal = &pipelineModal{
		pipeline:    &pipeline,
		projectPath: projectPath,
		tab:         pipelineJobsTab,
		selected:    0,
	}

	gui.renderPipelineModalHeader(header)
	gui.renderPipelineModalTabs(tabs)
	gui.renderPipelineModalContent(content)
	gui.renderPipelineModalFooter(footer)

	if _, err := gui.g.SetViewOnTop(pipelineModalHeaderView); err != nil {
		gui.closePipelineModal()

		return fmt.Errorf(
			"stack pipeline modal header: %w",
			err,
		)
	}

	if _, err := gui.g.SetViewOnTop(pipelineModalTabsView); err != nil {
		gui.closePipelineModal()

		return fmt.Errorf(
			"stack pipeline modal tabs: %w",
			err,
		)
	}

	if _, err := gui.g.SetViewOnTop(pipelineModalContentView); err != nil {
		gui.closePipelineModal()

		return fmt.Errorf(
			"stack pipeline modal content: %w",
			err,
		)
	}

	if _, err := gui.g.SetViewOnTop(pipelineModalFooterView); err != nil {
		gui.closePipelineModal()

		return fmt.Errorf(
			"stack pipeline modal footer: %w",
			err,
		)
	}

	if err := gui.Popup.Focus(pipelineModalContentView); err != nil {
		gui.closePipelineModal()

		return fmt.Errorf(
			"focus pipeline modal content: %w",
			err,
		)
	}

	return nil
}

func (gui *GUI) pipelineModalNextTab() {
	if gui.pipelineModal == nil {
		return
	}

	gui.pipelineModal.tab = (gui.pipelineModal.tab + 1) % 2
	gui.pipelineModal.selected = 0

	gui.renderPipelineModalView()
}

func (gui *GUI) pipelineModalNext() {
	if gui.pipelineModal == nil {
		return
	}

	count := gui.pipelineModalItemCount()

	if count == 0 {
		gui.pipelineModal.selected = 0
		return
	}

	v, err := gui.g.View(pipelineModalContentView)

	if err != nil || v == nil {
		return
	}

	_, cursorY := v.Cursor()
	originY := v.OriginY()

	selected := originY + cursorY
	selected++

	if selected >= count {
		selected = 0
	}

	gui.pipelineModal.selected = selected

	gui.renderPipelineModalView()
}

func (gui *GUI) pipelineModalPrevious() {
	if gui.pipelineModal == nil {
		return
	}

	count := gui.pipelineModalItemCount()

	if count == 0 {
		gui.pipelineModal.selected = 0
		return
	}

	v, err := gui.g.View(pipelineModalContentView)

	if err != nil || v == nil {
		return
	}

	_, cursorY := v.Cursor()
	originY := v.OriginY()

	selected := originY + cursorY
	selected--

	if selected < 0 {
		selected = count - 1
	}

	gui.pipelineModal.selected = selected

	gui.renderPipelineModalView()
}

func (gui *GUI) pipelineModalItemCount() int {
	if gui.pipelineModal == nil ||
		gui.State == nil ||
		gui.State.Pipelines == nil {
		return 0
	}

	switch gui.pipelineModal.tab {
	case pipelineJobsTab:
		return len(gui.selectedPipelineJobs())

	case pipelineDownstreamTab:
		return len(gui.selectedPipelineBridges())

	default:
		return 0
	}
}

func (gui *GUI) openSelectedPipelineItem() error {
	if gui.pipelineModal == nil {
		return nil
	}

	switch gui.pipelineModal.tab {
	case pipelineJobsTab:
		jobs := gui.selectedPipelineJobs()

		if len(jobs) == 0 {
			return nil
		}

		selected := gui.pipelineModal.selected

		if selected < 0 || selected >= len(jobs) {
			return nil
		}

		return gui.openJobLogModal(
			gui.pipelineModal.projectPath,
			jobs[selected],
		)

	case pipelineDownstreamTab:
		return gui.openSelectedDownstreamPipeline()

	default:
		return nil
	}
}

func (gui *GUI) openSelectedDownstreamPipeline() error {
	bridges := gui.selectedPipelineBridges()

	if len(bridges) == 0 {
		return nil
	}

	selected := gui.pipelineModal.selected

	if selected < 0 || selected >= len(bridges) {
		return nil
	}

	bridge := bridges[selected]

	if bridge.DownstreamPipeline == nil {
		return nil
	}

	downstream := bridge.DownstreamPipeline

	for {
		if downstream.ProjectPath == "" {
			return fmt.Errorf(
				"downstream pipeline #%d has no project path",
				downstream.ID,
			)
		}

		if err := gui.State.Pipelines.LoadPipelineDetails(
			gui.GitLab,
			downstream.ProjectPath,
			downstream.ID,
		); err != nil {
			return fmt.Errorf(
				"load downstream pipeline #%d details: %w",
				downstream.ID,
				err,
			)
		}

		hasJobs := false

		for _, job := range gui.State.Pipelines.Jobs {
			if job.PipelineID == downstream.ID {
				hasJobs = true
				break
			}
		}

		downstreamBridges := make([]gitlab.Bridge, 0)

		for _, bridge := range gui.State.Pipelines.Downstream {
			if bridge.PipelineID != downstream.ID {
				continue
			}

			downstreamBridges = append(
				downstreamBridges,
				bridge,
			)
		}

		if hasJobs || len(downstreamBridges) != 1 {
			break
		}

		next := downstreamBridges[0].DownstreamPipeline

		if next == nil {
			break
		}

		downstream = next
	}

	gui.savePipelineModalState()

	gui.pipelineModal.pipeline = &gitlab.Pipeline{
		ID:     downstream.ID,
		Status: downstream.Status,
		Ref:    downstream.Ref,
		SHA:    downstream.SHA,
		WebURL: downstream.WebURL,
	}

	gui.pipelineModal.projectPath = downstream.ProjectPath
	gui.pipelineModal.selected = 0

	if len(gui.selectedPipelineJobs()) > 0 {
		gui.pipelineModal.tab = pipelineJobsTab
	} else {
		gui.pipelineModal.tab = pipelineDownstreamTab
	}

	v, err := gui.g.View(pipelineModalView)
	if err != nil || v == nil {
		return nil
	}

	v.Title = fmt.Sprintf(
		"[Pipeline #%d]",
		downstream.ID,
	)

	gui.renderPipelineModalView()

	return nil
}

func (gui *GUI) selectedPipelineJobs() []gitlab.PipelineJobs {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return nil
	}

	pipelineID := gui.pipelineModal.pipeline.ID

	jobs := make([]gitlab.PipelineJobs, 0)

	for _, job := range gui.State.Pipelines.Jobs {
		if job.PipelineID == pipelineID {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

func (gui *GUI) selectedPipelineBridges() []gitlab.Bridge {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return nil
	}

	pipelineID := gui.pipelineModal.pipeline.ID

	bridges := make([]gitlab.Bridge, 0)

	for _, bridge := range gui.State.Pipelines.Downstream {
		if bridge.PipelineID == pipelineID {
			bridges = append(bridges, bridge)
		}
	}

	return bridges
}

func (gui *GUI) renderPipelineModalView() {
	if gui.pipelineModal == nil {
		return
	}

	header, err := gui.g.View(pipelineModalHeaderView)
	if err == nil && header != nil {
		gui.renderPipelineModalHeader(header)
	}

	tabs, err := gui.g.View(pipelineModalTabsView)
	if err == nil && tabs != nil {
		gui.renderPipelineModalTabs(tabs)
	}

	content, err := gui.g.View(pipelineModalContentView)
	if err == nil && content != nil {
		gui.renderPipelineModalContent(content)
	}

	footer, err := gui.g.View(pipelineModalFooterView)
	if err == nil && footer != nil {
		gui.renderPipelineModalFooter(footer)
	}
}

func (gui *GUI) renderPipelineModalHeader(v *gocui.View) {
	v.Clear()

	if gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return
	}

	pipeline := gui.pipelineModal.pipeline

	fmt.Fprintln(v)

	fmt.Fprintf(
		v,
		"Project: %s\n",
		gui.pipelineModal.projectPath,
	)

	fmt.Fprintf(
		v,
		"Status:  %s\n",
		pipeline.Status,
	)

	fmt.Fprintf(
		v,
		"Ref:     %s\n",
		pipeline.Ref,
	)

	fmt.Fprintf(
		v,
		"SHA:     %s\n",
		pipeline.SHA,
	)

	fmt.Fprintf(
		v,
		"Created: %s\n",
		formatPipelineTime(pipeline.CreatedAt),
	)

	fmt.Fprintf(
		v,
		"Updated: %s\n",
		formatPipelineTime(pipeline.UpdatedAt),
	)

	fmt.Fprintln(v)
	fmt.Fprintln(
		v,
		"────────────────────────────────────────────────────────",
	)
}

func (gui *GUI) renderPipelineModalTabs(v *gocui.View) {
	v.Clear()

	if gui.pipelineModal == nil {
		return
	}

	jobs := "Jobs"
	downstream := "Downstream"

	if gui.pipelineModal.tab == pipelineJobsTab {
		jobs = "[Jobs]"
	} else {
		downstream = "[Downstream]"
	}

	fmt.Fprintf(
		v,
		"%s  %s\n",
		jobs,
		downstream,
	)
}

func (gui *GUI) renderPipelineModalContent(v *gocui.View) {
	v.Clear()

	if gui.pipelineModal == nil {
		return
	}

	switch gui.pipelineModal.tab {
	case pipelineJobsTab:
		gui.renderPipelineJobs(v)

	case pipelineDownstreamTab:
		gui.renderPipelineDownstream(v)
	}

	gui.setPipelineModalCursor(
		v,
		gui.pipelineModal.selected,
	)
}

func (gui *GUI) renderPipelineModalFooter(v *gocui.View) {
	v.Clear()

	//fmt.Fprintln(v)

	fmt.Fprintln(
		v,
		"  [↑/↓] Select   [Tab] Switch downstream/jobs   [Enter] Open",
	)

	fmt.Fprintln(
		v,
		"  [r] Retry     [R] Refresh     [p] Run/Play    [Esc] Close",
	)
}

func (gui *GUI) renderPipelineJobs(v *gocui.View) {
	jobs := gui.selectedPipelineJobs()

	if len(jobs) == 0 {
		fmt.Fprintln(v, "  No jobs found.")
		return
	}

	for index, job := range jobs {
		fmt.Fprintf(
			v,
			"%3d. %s %s (ID: %d)\n",
			index+1,
			pipelineJobStatusSymbol(job.Status),
			job.Name,
			job.ID,
		)
	}
}

func (gui *GUI) renderPipelineDownstream(v *gocui.View) {
	bridges := gui.selectedPipelineBridges()

	if len(bridges) == 0 {
		fmt.Fprintln(v, "  No downstream pipelines found.")
		return
	}

	for index, bridge := range bridges {
		if bridge.DownstreamPipeline == nil {
			fmt.Fprintf(
				v,
				"%3d. %s (ID: %d) — no downstream pipeline\n",
				index+1,
				bridge.Name,
				bridge.ID,
			)

			continue
		}

		downstream := bridge.DownstreamPipeline

		fmt.Fprintf(
			v,
			"%3d. #%d %s %s\n",
			index+1,
			downstream.ID,
			downstream.Status,
			downstream.Ref,
		)
	}
}

func (gui *GUI) setPipelineModalCursor(
	v *gocui.View,
	selected int,
) {
	if v == nil {
		return
	}

	count := gui.pipelineModalItemCount()

	if count == 0 {
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)

		return
	}

	if selected < 0 {
		selected = 0
	}

	if selected >= count {
		selected = count - 1
	}

	visibleRows := v.InnerHeight()

	if visibleRows <= 0 {
		return
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

	maxOriginY := count - visibleRows

	if maxOriginY < 0 {
		maxOriginY = 0
	}

	if originY > maxOriginY {
		originY = maxOriginY
	}

	cursorY := selected - originY

	if cursorY < 0 {
		cursorY = 0
	}

	if cursorY >= visibleRows {
		cursorY = visibleRows - 1
	}

	v.SetOrigin(0, originY)
	v.SetCursor(0, cursorY)
}

func (gui *GUI) closePipelineModal() {
	gui.Popup.Close(pipelineModalFooterView)
	gui.Popup.Close(pipelineModalContentView)
	gui.Popup.Close(pipelineModalTabsView)
	gui.Popup.Close(pipelineModalHeaderView)
	gui.Popup.Close(pipelineModalView)

	gui.pipelineModal = nil

	gui.setFocus(focusPipelines)
}

func formatPipelineTime(t *time.Time) string {
	if t == nil {
		return "-"
	}

	return t.Format("2006-01-02 15:04")
}

func pipelineJobStatusSymbol(status string) string {
	switch status {
	case "success":
		return "✅"

	case "failed":
		return "❌"

	case "running":
		return "🔄"

	case "pending":
		return "⏳"

	case "created":
		return "🆕"

	case "canceled":
		return "⏹️"

	case "skipped":
		return "⏭️"

	case "manual":
		return "⏸️"

	case "scheduled":
		return "🕐"

	default:
		return "•"
	}
}

func (gui *GUI) retrySelectedPipelineJob() error {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return nil
	}

	jobs := gui.selectedPipelineJobs()

	if len(jobs) == 0 {
		return nil
	}

	selected := gui.pipelineModal.selected

	if selected < 0 || selected >= len(jobs) {
		return nil
	}

	job := jobs[selected]

	if _, err := gui.GitLab.RetryJob(
		gui.pipelineModal.projectPath,
		job.ID,
	); err != nil {
		return fmt.Errorf(
			"retry pipeline job %d: %w",
			job.ID,
			err,
		)
	}

	if err := gui.State.Pipelines.RefreshPipelineDetails(
		gui.GitLab,
		gui.pipelineModal.projectPath,
		gui.pipelineModal.pipeline.ID,
	); err != nil {
		return fmt.Errorf(
			"refresh pipeline #%d after retrying job %d: %w",
			gui.pipelineModal.pipeline.ID,
			job.ID,
			err,
		)
	}

	gui.pipelineModal.selected = 0
	gui.renderPipelineModalView()

	return nil
}

func (gui *GUI) refreshSelectedPipeline() error {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return nil
	}

	pipelineID := gui.pipelineModal.pipeline.ID

	if err := gui.State.Pipelines.RefreshPipelineDetails(
		gui.GitLab,
		gui.pipelineModal.projectPath,
		gui.pipelineModal.pipeline.ID,
	); err != nil {
		return fmt.Errorf(
			"refresh pipeline #%d: %w",
			pipelineID,
			err,
		)
	}

	gui.renderPipelineModalView()

	return nil
}

func (gui *GUI) pipelineModalEscape() error {
	if gui.pipelineModal == nil {
		return nil
	}

	/*
		If we are looking at the downstream tab, Esc first
		returns to the Jobs tab of the current pipeline.
	*/
	if gui.pipelineModal.tab == pipelineDownstreamTab {
		gui.pipelineModal.tab = pipelineJobsTab

		count := gui.pipelineModalItemCount()

		if count == 0 {
			gui.pipelineModal.selected = 0
		} else if gui.pipelineModal.selected >= count {
			gui.pipelineModal.selected = count - 1
		}

		gui.renderPipelineModalView()

		return nil
	}

	if gui.restorePreviousPipelineModalState() {
		gui.updatePipelineModalTitle()
		gui.renderPipelineModalView()

		return nil
	}

	gui.closePipelineModal()

	return nil
}

func (gui *GUI) updatePipelineModalTitle() {
	v, err := gui.g.View(pipelineModalView)

	if err != nil ||
		v == nil ||
		gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return
	}

	v.Title = fmt.Sprintf(
		"[Pipeline #%d]",
		gui.pipelineModal.pipeline.IID,
	)
}

func (gui *GUI) savePipelineModalState() {
	if gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return
	}

	gui.pipelineModal.history = append(
		gui.pipelineModal.history,
		pipelineModalState{
			pipeline:    *gui.pipelineModal.pipeline,
			projectPath: gui.pipelineModal.projectPath,
			tab:         gui.pipelineModal.tab,
			selected:    gui.pipelineModal.selected,
		},
	)
}

func (gui *GUI) restorePreviousPipelineModalState() bool {
	if gui.pipelineModal == nil ||
		len(gui.pipelineModal.history) == 0 {
		return false
	}

	last := len(gui.pipelineModal.history) - 1
	previous := gui.pipelineModal.history[last]

	gui.pipelineModal.history = gui.pipelineModal.history[:last]

	pipeline := previous.pipeline

	gui.pipelineModal.pipeline = &pipeline
	gui.pipelineModal.projectPath = previous.projectPath
	gui.pipelineModal.tab = previous.tab
	gui.pipelineModal.selected = previous.selected

	return true
}

func (gui *GUI) openSelectedJobURL() error {
	if gui.pipelineModal == nil {
		return nil
	}

	jobs := gui.selectedPipelineJobs()

	selected := gui.pipelineModal.selected

	if selected < 0 || selected >= len(jobs) {
		return nil
	}

	job := jobs[selected]

	if job.WebURL == "" {
		return fmt.Errorf(
			"job %d has no web URL",
			job.ID,
		)
	}

	if err := gui.OpenUrl(job.WebURL); err != nil {
		return fmt.Errorf(
			"open job %d: %w",
			job.ID,
			err,
		)
	}

	return nil
}

func (gui *GUI) openSelectedDownstreamPipelineURL() error {
	if gui.pipelineModal == nil {
		return nil
	}

	bridges := gui.selectedPipelineBridges()

	selected := gui.pipelineModal.selected

	if selected < 0 || selected >= len(bridges) {
		return nil
	}

	bridge := bridges[selected]

	if bridge.DownstreamPipeline == nil {
		return fmt.Errorf(
			"bridge %d has no downstream pipeline",
			bridge.ID,
		)
	}

	if bridge.DownstreamPipeline.WebURL == "" {
		return fmt.Errorf(
			"downstream pipeline %d has no web URL",
			bridge.DownstreamPipeline.ID,
		)
	}

	if err := gui.OpenUrl(bridge.DownstreamPipeline.WebURL); err != nil {
		return fmt.Errorf(
			"open downstream pipeline %d: %w",
			bridge.DownstreamPipeline.ID,
			err,
		)
	}

	return nil
}

func (gui *GUI) openSelectedPipelineURL() error {
	if gui.pipelineModal == nil {
		return nil
	}

	switch gui.pipelineModal.tab {
	case pipelineJobsTab:
		return gui.openSelectedJobURL()

	case pipelineDownstreamTab:
		return gui.openSelectedDownstreamPipelineURL()

	default:
		return nil
	}
}

func (gui *GUI) playSelectedPipelineJob() error {
	if gui.State == nil ||
		gui.State.Pipelines == nil ||
		gui.pipelineModal == nil ||
		gui.pipelineModal.pipeline == nil {
		return nil
	}

	jobs := gui.selectedPipelineJobs()

	if len(jobs) == 0 {
		return nil
	}

	selected := gui.pipelineModal.selected

	if selected < 0 || selected >= len(jobs) {
		return nil
	}

	job := jobs[selected]

	if _, err := gui.GitLab.PlayJob(
		gui.pipelineModal.projectPath,
		job.ID,
	); err != nil {
		return fmt.Errorf(
			"play pipeline job %d: %w",
			job.ID,
			err,
		)
	}

	if err := gui.State.Pipelines.RefreshPipelineDetails(
		gui.GitLab,
		gui.pipelineModal.projectPath,
		gui.pipelineModal.pipeline.ID,
	); err != nil {
		return fmt.Errorf(
			"refresh pipeline #%d after playing job %d: %w",
			gui.pipelineModal.pipeline.ID,
			job.ID,
			err,
		)
	}

	gui.pipelineModal.selected = 0
	gui.renderPipelineModalView()

	return nil
}

