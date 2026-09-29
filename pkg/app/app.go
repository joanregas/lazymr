package app

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"lazymr/pkg/commands/oscommands"
	"lazymr/pkg/config"
	"lazymr/pkg/debug"
	"lazymr/pkg/git"
	"lazymr/pkg/gitlab"
	"lazymr/pkg/gui"
)

type App struct {
	BuildInfo *BuildInfo
	CLIArgs   *cliArgs

	OSCommand *oscommands.OSCommand
	GitLab    *gitlab.Client
	GUI       *gui.GUI

	GitlabState *gitlab.GitlabState
	ConfigState *config.State

	CurrentDir string
	Project    string
}

func (app *App) Run() error {
	debug.Log("Starting App: %s ", app.Project)
	return app.GUI.Run()
}

func (app *App) verifyRepository() error {
	gitPath := filepath.Join(app.CurrentDir, ".git")

	fileInfo, err := os.Stat(gitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"not a git repository: %s",
				app.CurrentDir,
			)
		}

		return fmt.Errorf(
			"checking git repository: %w",
			err,
		)
	}

	// .git can either be a directory or a file.
	// A file is used by Git worktrees and some other Git layouts.
	if !fileInfo.IsDir() && !fileInfo.Mode().IsRegular() {
		return fmt.Errorf(
			"invalid git repository: %s",
			app.CurrentDir,
		)
	}

	return nil
}

func (app *App) verifyGitLabProject() error {
	_, err := app.GitLab.GetProject(app.Project)
	if err != nil {
		return fmt.Errorf(
			"no access to GitLab project %q: %w",
			app.Project,
			err,
		)
	}

	return nil
}

func (app *App) getGitLabProject() (string, error) {
	gitConfigPath := filepath.Join(
		app.CurrentDir,
		".git",
		"config",
	)

	data, err := os.ReadFile(gitConfigPath)
	if err != nil {
		return "", fmt.Errorf(
			"read git config %s: %w",
			gitConfigPath,
			err,
		)
	}

	remoteURL, err := git.ParseOriginURL(string(data))
	if err != nil {
		return "", err
	}

	project, err := gitLabProjectPath(remoteURL)
	if err != nil {
		return "", err
	}

	return project, nil
}

func gitLabProjectPath(remoteURL string) (string, error) {
	remoteURL = strings.TrimSpace(remoteURL)

	var path string

	switch {
	case strings.HasPrefix(remoteURL, "git@"):
		parts := strings.SplitN(remoteURL, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf(
				"invalid GitLab SSH remote URL: %q",
				remoteURL,
			)
		}

		path = parts[1]

	case strings.HasPrefix(remoteURL, "ssh://"):
		u, err := url.Parse(remoteURL)
		if err != nil {
			return "", fmt.Errorf(
				"parse GitLab SSH URL: %w",
				err,
			)
		}

		path = u.Path

	case strings.HasPrefix(remoteURL, "http://"),
		strings.HasPrefix(remoteURL, "https://"):
		u, err := url.Parse(remoteURL)
		if err != nil {
			return "", fmt.Errorf(
				"parse GitLab HTTP URL: %w",
				err,
			)
		}

		path = u.Path

	default:
		return "", fmt.Errorf(
			"unsupported Git remote URL: %q",
			remoteURL,
		)
	}

	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")

	if path == "" {
		return "", fmt.Errorf(
		"empty GitLab project path from remote URL %q",
		remoteURL,
		)
	}

	return path, nil
}

func NewApp(buildInfo *BuildInfo, cliArgs *cliArgs) (*App, error) {
	start := time.Now()

	osCommand := oscommands.NewOSCommand()

	app := &App{
		BuildInfo: buildInfo,
		CLIArgs:   cliArgs,
		OSCommand: osCommand,
	}

	logStep(start, "oscommands.NewOSCommand")

	cfg, _, err := config.Load()
	logStep(start, "config.Load")

	if err != nil {
		if errors.Is(err, config.ErrConfigurationRequired) {
			if guiErr := gui.ShowConfigurationRequired(); guiErr != nil {
				return nil, guiErr
			}

			return nil, err
		}

		return nil, err
	}

	currentDir, err := osCommand.GetCurrentDir()
	logStep(start, "GetCurrentDir")

	if err != nil {
		return nil, err
	}

	app.CurrentDir = currentDir

	configState, _, err := config.LoadState()
	logStep(start, "config.LoadState")

	if err != nil {
		return nil, err
	}

	app.ConfigState = configState

	showRecentRepos := false

	if err := app.verifyRepository(); err != nil {
		logStep(start, "verifyRepository failed")

		showRecentRepos = true

		if _, err := app.getRecentRepository(); err != nil {
			return nil, err
		}

		logStep(start, "getRecentRepository")
	} else {
		logStep(start, "verifyRepository")
	}

	gitlabClient, err := gitlab.NewClient(cfg)
	logStep(start, "gitlab.NewClient")

	if err != nil {
		return nil, err
	}

	app.GitLab = gitlabClient

	project, err := app.getGitLabProject()
	logStep(start, "getGitLabProject")

	if err != nil {
		return nil, err
	}

	app.Project = project

	app.GitlabState, err = app.loadGitlabState(0)
	logStep(start, "loadGitlabState")

	if err != nil {
		return nil, err
	}

	guiInstance, err := gui.New(
		app.GitlabState,
		app.ConfigState,
		app.GitLab,
		app.ChangeRepository,
		showRecentRepos,
		app.OpenURL,
		app.refreshGitlabState,
	)
	logStep(start, "gui.New")

	if err != nil {
		return nil, err
	}

	app.GUI = guiInstance

	logStep(start, "NewApp complete")

	return app, nil
}

func (app *App) ChangeRepository(path string) error {
	path = filepath.Clean(path)

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf(
			"check repository directory %q: %w",
			path,
			err,
		)
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"repository path is not a directory: %s",
			path,
		)
	}

	if err := app.OSCommand.ChangeDir(path); err != nil {
		return fmt.Errorf(
			"change directory to %q: %w",
			path,
			err,
		)
	}

	app.CurrentDir = path

	if err := app.verifyRepository(); err != nil {
		return err
	}

	project, err := app.getGitLabProject()
	if err != nil {
		return err
	}

	app.Project = project

	if err := app.verifyGitLabProject(); err != nil {
		return err
	}

	uiState, err := app.loadGitlabState(0)
	if err != nil {
		return err
	}

	app.GitlabState = uiState
	app.GUI.State = uiState

	app.ConfigState.AddRecentRepo(path)

	if err := config.SaveState(app.ConfigState); err != nil {
		return fmt.Errorf(
			"save recent repositories: %w",
			err,
		)
	}

	return app.GUI.Refresh()
}

func (app *App) refreshGitlabState(
	selectedMergeRequestIID int,
) (*gitlab.GitlabState, error) {
	return app.loadGitlabState(selectedMergeRequestIID)
}

func (app *App) loadGitlabState(
	selectedMergeRequestIID int,
) (*gitlab.GitlabState, error) {
	start := time.Now()

	var (
		repository    *gitlab.Repository
		mergeRequests *gitlab.MergeRequests
		members       *gitlab.Members

		repositoryErr    error
		mergeRequestsErr error
		membersErr       error

		wg sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		defer wg.Done()

		repository, repositoryErr = gitlab.NewRepository(
			app.GitLab,
			app.Project,
		)
	}()

	go func() {
		defer wg.Done()

		mergeRequests, mergeRequestsErr = gitlab.NewMergeRequests(
			app.GitLab,
			app.Project,
		)
	}()

	go func() {
		defer wg.Done()

		members, membersErr = gitlab.NewMembers(
			app.GitLab,
			app.Project,
		)
	}()

	wg.Wait()

	logStep(
		start,
		"NewRepository + NewMergeRequests + NewMembers",
	)

	if repositoryErr != nil {
		return nil, fmt.Errorf(
			"load GitLab repository: %w",
			repositoryErr,
		)
	}

	if mergeRequestsErr != nil {
		return nil, fmt.Errorf(
			"load GitLab merge requests: %w",
			mergeRequestsErr,
		)
	}

	if membersErr != nil {
		return nil, fmt.Errorf(
			"load GitLab project members: %w",
			membersErr,
		)
	}

	selectedIndex := 0

	if selectedMergeRequestIID > 0 {
		for index, mergeRequest := range mergeRequests.Items {
			if mergeRequest.IID == selectedMergeRequestIID {
				selectedIndex = index
				break
			}
		}
	}

	var (
		files     *gitlab.Files
		pipelines *gitlab.Pipelines
	)

	if len(mergeRequests.Items) > 0 {
		selectedMergeRequest := mergeRequests.Items[selectedIndex]

		var (
			filesErr     error
			pipelinesErr error
		)

		wg.Add(2)

		go func() {
			defer wg.Done()

			files, filesErr = gitlab.NewFiles(
				app.GitLab,
				app.Project,
				selectedMergeRequest.IID,
			)
		}()

		go func() {
			defer wg.Done()

			pipelines, pipelinesErr = gitlab.NewPipelines(
				app.GitLab,
				app.Project,
				selectedMergeRequest.IID,
			)
		}()

		wg.Wait()

		logStep(start, "NewFiles + NewPipelines")

		if filesErr != nil {
			return nil, fmt.Errorf(
				"load GitLab merge request files: %w",
				filesErr,
			)
		}

		if pipelinesErr != nil {
			return nil, fmt.Errorf(
				"load GitLab merge request pipelines: %w",
				pipelinesErr,
			)
		}
	}

	state := gitlab.NewGitlabState(
		repository,
		nil,
		mergeRequests,
		members,
		files,
		pipelines,
	)

	state.SelectedMergeRequest = selectedIndex

	return state, nil
}

func (app *App) getRecentRepository() (string, error) {
	repository, ok := app.ConfigState.GetRecentRepo()
	if !ok {
		return "", fmt.Errorf(
			"not in a git repository: %s: no valid recent repositories",
			app.CurrentDir,
		)
	}

	if err := app.OSCommand.ChangeDir(repository); err != nil {
		return "", fmt.Errorf(
			"change directory to recent repository %q: %w",
			repository,
			err,
		)
	}

	app.CurrentDir = repository

	if err := app.verifyRepository(); err != nil {
		return "", fmt.Errorf(
			"recent repository %q is not valid: %w",
			repository,
			err,
		)
	}

	return repository, nil
}

func logStep(start time.Time, name string) {
	debug.Log(
		"[TIMING] %-35s %v",
		name,
		time.Since(start),
	)
}

func (app *App) OpenURL(url string) error {
	return app.OSCommand.OpenURL(url)
}

