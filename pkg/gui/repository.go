package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"

	"lazymr/pkg/gitlab"
)

const (
	recentRepositoriesView = "recent_repositories"
)

func createRepositoryView(
	g *gocui.Gui,
	repository *gitlab.Repository,
) (*gocui.View, error) {
	v, err := createView(
		g,
		repositoriesView,
		"[1]-Repository",
	)
	if v == nil {
		return nil, err
	}

	renderRepository(v, repository)

	return v, nil
}

func renderRepository(
	v *gocui.View,
	repository *gitlab.Repository,
) {
	v.Clear()

	if repository == nil {
		return
	}

	fmt.Fprintln(v, repository.Name)
	fmt.Fprintln(v, repository.PathWithNamespace)
}

func (gui *GUI) openRepositoryPopup() error {
	if gui.ConfigState == nil {
		return nil
	}

	if gui.Popup.Exists(recentRepositoriesView) {
		return nil
	}

	width, height := gui.g.Size()

	popupWidth := width * 2 / 3

	if popupWidth < 40 {
		popupWidth = 40
	}

	if popupWidth > width-4 {
		popupWidth = width - 4
	}

	popupHeight := len(gui.ConfigState.RecentRepos) + 4

	if popupHeight > height-2 {
		popupHeight = height - 2
	}

	x0 := (width - popupWidth) / 2
	y0 := (height - popupHeight) / 2
	x1 := x0 + popupWidth
	y1 := y0 + popupHeight

	v, err := gui.Popup.Create(
		recentRepositoriesView,
		"[1]-Recent Projects",
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf(
			"create recent repositories popup: %w",
			err,
		)
	}

	gui.renderRepositoryPopup(v)
	gui.initializeRepositoryPopupKeybindings()

	if err := gui.Popup.Focus(recentRepositoriesView); err != nil {
		return fmt.Errorf(
			"focus recent repositories popup: %w",
			err,
		)
	}

	return nil
}

func (gui *GUI) renderRepositoryPopup(
	v *gocui.View,
) {
	v.Clear()

	repositories := gui.ConfigState.RecentRepos

	if len(repositories) == 0 {
		fmt.Fprintln(v)
		fmt.Fprintln(v, "No recent projects.")

		return
	}

	for _, repository := range repositories {
		fmt.Fprintln(v, repository)
	}

	v.Footer = fmt.Sprintf(
		"%d repositories",
		len(repositories),
	)

	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
}

/* TO-DO: Move this out of here */
func (gui *GUI) initializeRepositoryPopupKeybindings() {
	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyArrowUp),
		func(g *gocui.Gui, v *gocui.View) error {
			gui.moveRepositoryCursor(v, -1)

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyArrowDown),
		func(g *gocui.Gui, v *gocui.View) error {
			gui.moveRepositoryCursor(v, 1)

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyPgup),
		func(g *gocui.Gui, v *gocui.View) error {
			gui.moveRepositoryCursor(
				v,
				-gui.repositoryPopupPageSize(v),
			)

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyPgdn),
		func(g *gocui.Gui, v *gocui.View) error {
			gui.moveRepositoryCursor(
				v,
				gui.repositoryPopupPageSize(v),
			)

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyHome),
		func(g *gocui.Gui, v *gocui.View) error {
			gui.setRepositoryCursor(v, 0)

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyEnd),
		func(g *gocui.Gui, v *gocui.View) error {
			repositories := gui.ConfigState.RecentRepos

			if len(repositories) == 0 {
				return nil
			}

			gui.setRepositoryCursor(
				v,
				len(repositories)-1,
			)

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyEsc),
		func(g *gocui.Gui, v *gocui.View) error {
			gui.closeRepositoryPopup()

			return nil
		},
	)

	gui.g.SetKeybinding(
		recentRepositoriesView,
		gocui.NewKeyName(gocui.KeyEnter),
		func(g *gocui.Gui, v *gocui.View) error {
			return gui.selectRepository(v)
		},
	)
}

func (gui *GUI) moveRepositoryCursor(
	v *gocui.View,
	delta int,
) {
	repositories := gui.ConfigState.RecentRepos

	if len(repositories) == 0 {
		return
	}

	_, cursorY := v.Cursor()
	originY := v.OriginY()

	selected := originY + cursorY
	selected += delta

	if selected < 0 {
		selected = 0
	}

	if selected >= len(repositories) {
		selected = len(repositories) - 1
	}

	gui.setRepositoryCursor(v, selected)
}

func (gui *GUI) setRepositoryCursor(
	v *gocui.View,
	selected int,
) {
	repositories := gui.ConfigState.RecentRepos

	if len(repositories) == 0 {
		return
	}

	if selected < 0 {
		selected = 0
	}

	if selected >= len(repositories) {
		selected = len(repositories) - 1
	}

	_, cursorY := v.Cursor()
	originY := v.OriginY()

	_, height := v.Size()

	if height <= 0 {
		return
	}

	visibleRows := height - 1

	if visibleRows < 1 {
		visibleRows = 1
	}

	if selected < originY {
		originY = selected
	}

	if selected >= originY+visibleRows {
		originY = selected - visibleRows + 1
	}

	if originY < 0 {
		originY = 0
	}

	maxOriginY := len(repositories) - visibleRows

	if maxOriginY < 0 {
		maxOriginY = 0
	}

	if originY > maxOriginY {
		originY = maxOriginY
	}

	cursorY = selected - originY

	if cursorY < 0 {
		cursorY = 0
	}

	if cursorY >= visibleRows {
		cursorY = visibleRows - 1
	}

	v.SetOrigin(0, originY)
	v.SetCursor(0, cursorY)
}

func (gui *GUI) repositoryPopupPageSize(
	v *gocui.View,
) int {
	_, height := v.Size()

	pageSize := height - 1

	if pageSize < 1 {
		pageSize = 1
	}

	return pageSize
}

func (gui *GUI) closeRepositoryPopup() {
	gui.Popup.Close(recentRepositoriesView)

	gui.setFocus(focusRepositories)
}

func (gui *GUI) selectRepository(
	v *gocui.View,
) error {
	repositories := gui.ConfigState.RecentRepos

	if len(repositories) == 0 {
		return nil
	}

	_, cursorY := v.Cursor()
	originY := v.OriginY()

	selected := originY + cursorY

	if selected < 0 || selected >= len(repositories) {
		return nil
	}

	selectedRepository := repositories[selected]

	if gui.ChangeRepository == nil {
		return fmt.Errorf(
			"repository change handler is not configured",
		)
	}

	gui.closeRepositoryPopup()

	return gui.ChangeRepository(selectedRepository)
}

func (gui *GUI) refreshRepository() error {
	if gui.Views.Repository == nil {
		return fmt.Errorf(
			"repository view is not initialized",
		)
	}

	renderRepository(
		gui.Views.Repository,
		gui.State.GetRepository(),
	)

	return nil
}

