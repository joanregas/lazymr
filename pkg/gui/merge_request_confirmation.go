package gui

import "fmt"

const (
	mergeRequestConfirmationView = "merge_request_confirmation"
)

type mergeRequestConfirmation struct {
	gui *GUI

	confirm func() error
}

func newMergeRequestConfirmation(gui *GUI) *mergeRequestConfirmation {
	return &mergeRequestConfirmation{
		gui: gui,
	}
}

func (confirmation *mergeRequestConfirmation) open(
	title string,
	message string,
	confirm func() error,
) error {
	confirmation.confirm = confirm

	if confirmation.gui.Popup.Exists(
		mergeRequestConfirmationView,
	) {
		confirmation.gui.Popup.Close(
			mergeRequestConfirmationView,
		)
	}

	width, height := confirmation.gui.g.Size()

	modalWidth := width * 2 / 3

	if modalWidth < 60 {
		modalWidth = 60
	}

	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	modalHeight := 8

	if modalHeight > height-2 {
		modalHeight = height - 2
	}

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2
	x1 := x0 + modalWidth
	y1 := y0 + modalHeight

	v, err := confirmation.gui.Popup.Create(
		mergeRequestConfirmationView,
		title,
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf(
			"create confirmation popup: %w",
			err,
		)
	}

	v.Clear()

	fmt.Fprintln(v)
	fmt.Fprintln(v, message)
	fmt.Fprintln(v)
	fmt.Fprintln(v, "[y] Yes    [n] No")

	v.FrameColor = colorWhite
	v.TitleColor = colorWhite

	_, err = confirmation.gui.g.SetViewOnTop(
		mergeRequestConfirmationView,
	)
	if err != nil {
		return fmt.Errorf(
			"raise confirmation popup: %w",
			err,
		)
	}

	return confirmation.gui.Popup.Focus(
		mergeRequestConfirmationView,
	)
}

func (confirmation *mergeRequestConfirmation) close() error {
	confirmation.gui.Popup.Close(
		mergeRequestConfirmationView,
	)

	confirmation.confirm = nil

	return confirmation.restorePreviousView()
}

func (confirmation *mergeRequestConfirmation) confirmAction() error {
	if confirmation.confirm == nil {
		return confirmation.close()
	}

	confirm := confirmation.confirm

	confirmation.gui.Popup.Close(
		mergeRequestConfirmationView,
	)

	confirmation.confirm = nil
	confirmation.gui.mergeRequestConfirmationPreviousView = ""

	return confirm()
}

func (confirmation *mergeRequestConfirmation) cancelAction() error {
	return confirmation.close()
}

func (confirmation *mergeRequestConfirmation) restorePreviousView() error {
	viewName := confirmation.gui.mergeRequestConfirmationPreviousView

	confirmation.gui.mergeRequestConfirmationPreviousView = ""

	if viewName == "" {
		viewName = repositoriesView
	}

	if err := confirmation.gui.focusView(
		viewName,
		focusableViews,
	); err != nil {
		return fmt.Errorf(
			"restore previous view %q: %w",
			viewName,
			err,
		)
	}

	footer, err := confirmation.gui.g.View(footerView)
	if err == nil {
		renderFooter(footer, viewName)
	}

	return nil
}

