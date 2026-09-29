package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/config"
	"lazymr/pkg/debug"
	"lazymr/pkg/gitlab"
	"lazymr/pkg/gui/keybindings"
	"lazymr/pkg/gui/popup"
)

/* Strings for windows  we will use them for keybindings */
const (
	repositoriesView = "repositories"
	mergeRequestView = "merge_request"
	filesView        = "files"
	pipelinesView   = "pipelines"
	statusView       = "status"
	overviewView     = "overview"
	commandLogsView  = "command_logs"
	footerView       = "footer"
)

/* TO-DO: Work in progress ... */
type Views struct {
	Repository    *gocui.View
	MergeRequests *gocui.View
	Files         *gocui.View
	Pipelines     *gocui.View
	Status        *gocui.View
	Overview      *gocui.View
	CommandLogs   *gocui.View
	Footer        *gocui.View
}

/* TO-DO: This is starting to feel....half-baked... take a look and see if we can move things inside a general struct or something 
          Most of it was constructed without a second tought on what i already built so im pretty sure we can reuse things 
*/
type GUI struct {
	g *gocui.Gui

	ConfigState *config.State
	State       *gitlab.GitlabState
	GitLab      *gitlab.Client

	/* TO-DO: Move this to gitlabstates was harder than anticipated, for now it stays as its */
	ChangeRepository func(string) error
	OpenUrl          func(string) error
	refreshGitlabState func(int) (*gitlab.GitlabState, error)

	commandLogsVisible bool

	/* This ones are still not migrated to popups , old approach */
	pipelineSelector *pipelineSelector
	pipelineModal    *pipelineModal
	jobLogModal      *jobLogModal
	showRecentRepos  bool

	mergeRequestCreator     *mergeRequestCreator
	mergeRequestConfirmation *mergeRequestConfirmation

	windowArrangement *WindowArrangement

	helpPreviousView               string
	mergeRequestConfirmationPreviousView string

	/* TO-DO: Work in progress for views and popups */
	Views          *Views
	Popup          *popup.PopupHandler

	loadingPreviousView string
}

func New(
	state *gitlab.GitlabState,
	configState *config.State,
	gitlabClient *gitlab.Client,
	changeRepository func(string) error,
	showRecentRepos bool,
	openUrl func(string) error,
	refreshGitlabState func(int) (*gitlab.GitlabState, error),
) (*GUI, error) {
	g, err := gocui.NewGui(
		gocui.NewGuiOpts{
			OutputMode: gocui.OutputTrue,
		},
	)
	g.ShowListFooter = true

	if err != nil {
		return nil, fmt.Errorf(
			"initialize GUI: %w",
			err,
		)
	}

	gui := &GUI{
		g:                  g,
		windowArrangement:  NewWindowArrangement(),
		State:              state,
		ConfigState:        configState,
		GitLab:             gitlabClient,
		ChangeRepository:   changeRepository,
		pipelineSelector:   &pipelineSelector{},
		showRecentRepos:    showRecentRepos,
		commandLogsVisible: true,
		Views:              &Views{},
		OpenUrl:             openUrl,
		refreshGitlabState: refreshGitlabState,
	}

	gui.mergeRequestCreator = newMergeRequestCreator(gui)
	gui.Popup = popup.NewPopupHandler(gui.g)
	gui.mergeRequestConfirmation = newMergeRequestConfirmation(gui)

	g.SetManagerFunc(gui.layout)

	if err := gui.initializeViews(); err != nil {
		g.Close()
		return nil, err
	}

	gui.setFocus(focusRepositories)
	gui.initializeKeybindings()

	return gui, nil
}

func (gui *GUI) Run() error {
	defer gui.g.Close()

	if err := gui.StartupOpenRecentRepositories(); err != nil {
		return err
	}

	if err := gui.g.MainLoop(); err != nil &&
		err != gocui.ErrQuit {
		return fmt.Errorf(
			"GUI main loop: %w",
			err,
		)
	}

	return nil
}

func (gui *GUI) Close() {
	gui.g.Close()
}

/* Init the ones that will appear on the main UI */
func (gui *GUI) initializeViews() error {
	repositoryView, err := createRepositoryView(
		gui.g,
		gui.State.GetRepository(),
	)
	if err != nil {
		return err
	}

	gui.Views.Repository = repositoryView

	if err := createMergeRequestView(
		gui.g,
		gui.State.GetMergeRequests(),
	); err != nil {
		return err
	}

	if err := createFilesView(
		gui.g,
		gui.State.GetFiles(),
	); err != nil {
		return err
	}

	if err := createPipelinesView(
		gui.g,
		gui.State.GetPipelines(),
	); err != nil {
		return err
	}

	if err := createStatusView(
		gui.g,
		gui.State.GetMergeRequests(),
	); err != nil {
		return err
	}

	if err := createOverviewView(gui.g); err != nil {
		return err
	}

	if err := createCommandLogsView(gui.g); err != nil {
		return err
	}

	if err := createFooterView(gui.g); err != nil {
		return err
	}

	return nil
}

/* TO-DO: This is a work in progress and we are moving to gui/keybindings */
func (gui *GUI) initializeKeybindings() {
	keybindings.InitializeGlobal(
		gui.g,
		focusableViews,
		gui.focusHandler,
		gui.nextFocus,
		gui.mergeRequestCreator.open,
		gui.openHelp,
		func() error { return gocui.ErrQuit },
		gui.RefreshData,
	)
	keybindings.InitializeOverview(
		gui.g,
		overviewView,
		gui.scrollOverview,
		gui.overviewPageSize,
		gui.overviewHome,
		gui.overviewEnd,
	)
	keybindings.InitializeMergeRequest(
		gui.g,
		mergeRequestView,
		gui.closeSelectedMergeRequest,
		gui.mergeSelectedMergeRequest,
		gui.selectNextMergeRequest,
		gui.selectPreviousMergeRequest,
	)
	keybindings.InitializeFiles(
		gui.g,
		filesView,
		gui.toggleSelectedFileTreeNode,
		gui.selectNextFile,
		gui.selectPreviousFile,
	)
	keybindings.InitializeRepository(
		gui.g,
		repositoriesView,
		gui.openRepositoryPopup,
	)
	keybindings.InitializePipelines(
		gui.g,
		pipelinesView,
		gui.pipelinePrevious,
		gui.pipelineNext,
		gui.openPipelineModal,
		gui.openSelectedPipelineBrowserURL,
	)
	keybindings.InitializeCommandLogs(
		gui.g,
		gui.toggleCommandLogs,
	)
	keybindings.InitializeHelp(
		gui.g,
		helpView,
		gui.closeHelp,
	)
	keybindings.InitializeMergeRequestConfirmation(
		gui.g,
		mergeRequestConfirmationView,
		gui.mergeRequestConfirmation.confirmAction,
		gui.mergeRequestConfirmation.cancelAction,
	)

	/*TO-DO: Ok i think i overdid this one.... becuase of focus and toggle ... focus was fixed...
	        take a look later and try something similar we did with focus */
/*	keybindings.InitializeMergeRequestCreator(
		gui.g,
		mergeRequestCreatorInputViews,
		gui.mergeRequestCreator.nextFocus,
		gui.mergeRequestCreator.toggleDraft,
		gui.mergeRequestCreator.toggleDeleteSourceBranch,
		gui.mergeRequestCreator.toggleSquash,
		gui.mergeRequestCreator.confirmCreate,
		gui.mergeRequestCreator.cancel,
	)*/

	keybindings.InitializePipelineJobs(
		gui.g,
		pipelineModalContentView,
		gui.pipelineModalPrevious,
		gui.pipelineModalNext,
		gui.pipelineModalNextTab,
		gui.openSelectedPipelineItem,
		gui.openSelectedPipelineURL,
		gui.playSelectedPipelineJob,
		gui.retrySelectedPipelineJob,
		gui.refreshSelectedPipeline,
		gui.pipelineModalEscape,
	)

	keybindings.InitializeMergeRequestCreator(
		gui.g,
		mergeRequestCreatorView,
		mergeRequestCreatorSelectorView,
		mergeRequestCreatorSelectorSearchView,
		mergeRequestCreatorEditorView,
		gui.mergeRequestCreator.previous,
		gui.mergeRequestCreator.next,
		gui.mergeRequestCreator.edit,
		gui.mergeRequestCreator.toggle,
		gui.mergeRequestCreator.create,
		gui.mergeRequestCreator.close,
		gui.mergeRequestCreator.selectorPrevious,
		gui.mergeRequestCreator.selectorNext,
		gui.mergeRequestCreator.selectItem,
		gui.mergeRequestCreator.openSelectorSearch,
		gui.mergeRequestCreator.closeSelector,
		gui.mergeRequestCreator.executeSelectorSearch,
		gui.mergeRequestCreator.closeSelectorSearch,
		gui.mergeRequestCreator.editorEnter,
		gui.mergeRequestCreator.saveEditor,
		gui.mergeRequestCreator.closeEditor,
	)
	
	keybindings.InitializePipelineJobLogs(
		gui.g,
		jobLogModalView,
		func(delta int) error {
			gui.scrollJobLog(delta)
			return nil
		},
		gui.jobLogPageSize,
		gui.jobLogHome,
		gui.jobLogEnd,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.closeJobLogModal()
			return nil
		},
	)
}

func createView(
	g *gocui.Gui,
	name string,
	title string,
) (*gocui.View, error) {
	v, err := g.SetView(
		name,
		0,
		0,
		1,
		1,
		0,
	)

	// The local gocui fork can return ErrUnknownView
	// while still returning the newly created view.
	if v == nil {
		return nil, fmt.Errorf(
			"create view %q: %w",
			name,
			err,
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

	v.FrameColor = colorBlue
	v.Title = title
	v.TitleColor = colorBlue

	return v, nil
}

/* Ui Refresh just in case */
func (gui *GUI) Refresh() error {
	if err := gui.refreshMergeRequests(); err != nil {
		return err
	}

	if err := gui.refreshRepository(); err != nil {
		return err
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}

	if err := gui.refreshPipelines(); err != nil {
		return err
	}

	if err := gui.refreshStatus(); err != nil {
		return err
	}

	if err := gui.refreshOverview(); err != nil {
		return err
	}

	if err := gui.refreshCommandLogs(); err != nil {
		return err
	}

	return nil
}

/* if no repo is found on the current dir, and because we arnt lazygit 
   and shouldnt be `git initing` things, 
	 we set the first window to the select recent repos popup */
func (gui *GUI) StartupOpenRecentRepositories() error {
	if !gui.showRecentRepos {
		return nil
	}

	gui.showRecentRepos = false

	return gui.openRepositoryPopup()
}

/* Refresh GitLab data without restarting the application. */
func (gui *GUI) RefreshData() error {
	debug.Log("Refreshing Data")

	if gui.refreshGitlabState == nil {
		return fmt.Errorf("GitLab refresh function is not initialized")
	}

	selectedIID := 0

	if gui.State != nil &&
		gui.State.MergeRequests != nil {
		selectedIndex := gui.State.SelectedMergeRequest

		if selectedIndex >= 0 &&
			selectedIndex < len(gui.State.MergeRequests.Items) {
			selectedIID =
				gui.State.MergeRequests.Items[selectedIndex].IID
		}
	}

	if err := gui.startLoadingScreen(); err != nil {
		return err
	}

	go func() {
		state, err := gui.refreshGitlabState(selectedIID)

		gui.g.Update(func(g *gocui.Gui) error {
			gui.stopLoadingScreen()

			if err != nil {
				if gui.State != nil &&
					gui.State.GetCommandLogs() != nil {
					gui.State.GetCommandLogs().Add(
						fmt.Sprintf(
							"Failed to refresh GitLab data: %s",
							err,
						),
						true,
					)

					_ = gui.refreshCommandLogs()
				}

				debug.Log(
					"Refresh Data failed: %v",
					err,
				)

				return nil
			}

			gui.State = state

			if err := gui.Refresh(); err != nil {
				debug.Log(
					"Refresh GUI failed: %v",
					err,
				)
			}

			return nil
		})
	}()

	return nil
}

