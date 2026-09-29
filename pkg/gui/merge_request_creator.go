package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/gitlab"
)

const (
	mergeRequestCreatorView       = "create_merge_request"
	mergeRequestCreatorFooterView = "create_merge_request_footer"

	mergeRequestCreatorSelectorView       = "create_merge_request_selector"
	mergeRequestCreatorSelectorSearchView = "create_merge_request_selector_search"
	mergeRequestCreatorEditorView         = "create_merge_request_editor"
)

type mergeRequestCreatorField int

const (
	mergeRequestCreatorSourceBranchField mergeRequestCreatorField = iota
	mergeRequestCreatorTargetBranchField
	mergeRequestCreatorAssigneeField
	mergeRequestCreatorReviewerField
	mergeRequestCreatorTitleField
	mergeRequestCreatorDescriptionField
	mergeRequestCreatorDraftField
	mergeRequestCreatorDeleteField
	mergeRequestCreatorSquashField
)

type mergeRequestCreator struct {
	gui *GUI

	selectedField mergeRequestCreatorField

	sourceBranch string
	targetBranch string

	assignee   string
	assigneeID int

	reviewer   string
	reviewerID int

	titleValue       string
	descriptionValue string

	draft              bool
	deleteSourceBranch bool
	squash             bool

	selectorOpen     bool
	selectorSelected int
	selectorItems    []string
	selectorField    mergeRequestCreatorField

	editorOpen  bool
	editorField mergeRequestCreatorField

	selectorSearchOpen bool

	previousView string

	active bool
}

func newMergeRequestCreator(gui *GUI) *mergeRequestCreator {
	return &mergeRequestCreator{
		gui: gui,
	}
}

func (creator *mergeRequestCreator) open() error {
	if creator.active {
		return nil
	}

	currentView := creator.gui.g.CurrentView()

	creator.previousView = ""

	if currentView != nil {
		creator.previousView = currentView.Name()
	}

	creator.reset()

	if err := creator.loadBranches(); err != nil {
		return err
	}

	creator.active = true

	return creator.openForm()
}

func (creator *mergeRequestCreator) reset() {
	creator.selectedField =
		mergeRequestCreatorSourceBranchField

	creator.sourceBranch = ""
	creator.targetBranch = ""

	creator.assignee = ""
	creator.assigneeID = 0

	creator.reviewer = ""
	creator.reviewerID = 0

	creator.titleValue = ""
	creator.descriptionValue = ""

	creator.draft = false
	creator.deleteSourceBranch = false
	creator.squash = false

	creator.selectorOpen = false
	creator.selectorSelected = 0
	creator.selectorItems = nil
	creator.selectorField =
		mergeRequestCreatorSourceBranchField

	creator.editorOpen = false
	creator.editorField =
		mergeRequestCreatorTitleField

	creator.selectorSearchOpen = false
}

func (creator *mergeRequestCreator) close() error {
	creator.closeSelectorSearch()
	creator.closeSelector()
	creator.closeEditor()

	if err := creator.closeViews(); err != nil {
		return err
	}

	previousView := creator.previousView

	creator.previousView = ""
	creator.active = false

	if previousView != "" {
		creator.gui.setFocusByName(previousView)
		return nil
	}

	creator.gui.setFocus(focusRepositories)

	return nil
}


func (creator *mergeRequestCreator) closeViews() error {
	for _, name := range []string{
		mergeRequestCreatorView,
		mergeRequestCreatorFooterView,
		mergeRequestCreatorSelectorView,
		mergeRequestCreatorSelectorSearchView,
		mergeRequestCreatorEditorView,
	} {
		if _, err := creator.gui.g.View(name); err == nil {
			if err := creator.gui.g.DeleteView(name); err != nil {
				return fmt.Errorf(
					"delete merge request creator view %q: %w",
					name,
					err,
				)
			}
		}
	}

	return nil
}

func (creator *mergeRequestCreator) loadBranches() error {
	if creator.gui.State == nil {
		return fmt.Errorf(
			"cannot load branches: UI state is unavailable",
		)
	}

	branches := creator.gui.State.GetBranches()

	if branches == nil {
		mergeRequests := creator.gui.State.GetMergeRequests()

		if mergeRequests == nil ||
			mergeRequests.Project == "" {
			return fmt.Errorf(
				"cannot load branches: GitLab project is unavailable",
			)
		}

		if creator.gui.GitLab == nil {
			return fmt.Errorf(
				"cannot load branches: GitLab client is unavailable",
			)
		}

		var err error

		branches, err = gitlab.NewBranches(
			creator.gui.GitLab,
			mergeRequests.Project,
		)
		if err != nil {
			return fmt.Errorf(
				"load GitLab branches: %w",
				err,
			)
		}

		creator.gui.State.Branches = branches
	}

	if len(branches.Items) == 0 {
		return fmt.Errorf(
			"branches state is empty for project %q",
			branches.Project,
		)
	}

	return nil
}

func (creator *mergeRequestCreator) loadMembers() ([]gitlab.Member, error) {
	if creator.gui.State == nil {
		return nil, fmt.Errorf(
			"cannot load members: UI state is unavailable",
		)
	}

	mergeRequests := creator.gui.State.GetMergeRequests()

	if mergeRequests == nil ||
		mergeRequests.Project == "" {
		return nil, fmt.Errorf(
			"cannot load members: GitLab project is unavailable",
		)
	}

	if creator.gui.GitLab == nil {
		return nil, fmt.Errorf(
			"cannot load members: GitLab client is unavailable",
		)
	}

	members, err := gitlab.NewMembers(
		creator.gui.GitLab,
		mergeRequests.Project,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"load GitLab project members: %w",
			err,
		)
	}

	return members.Items, nil
}

func (creator *mergeRequestCreator) openForm() error {
	if err := creator.closeViews(); err != nil {
		return err
	}

	width, height := creator.gui.g.Size()

	modalWidth := width * 4 / 5

	if modalWidth < 80 {
		modalWidth = 80
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

	modal, err := creator.gui.Popup.Create(
		mergeRequestCreatorView,
		"Create Merge Request",
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf(
			"create merge request creator: %w",
			err,
		)
	}

	modal.Highlight = false

	footerHeight := 4

	contentY0 := y0 + 1
	contentY1 := y1 - footerHeight - 1

	footerY0 := contentY1 + 1
	footerY1 := y1 - 1

	content, err := creator.gui.g.SetView(
		mergeRequestCreatorView,
		x0+1,
		contentY0,
		x1-1,
		contentY1,
		0,
	)
	if content == nil {
		creator.gui.Popup.Close(mergeRequestCreatorView)

		return fmt.Errorf(
			"create merge request creator content: %w",
			err,
		)
	}

	footer, err := creator.gui.g.SetView(
		mergeRequestCreatorFooterView,
		x0+1,
		footerY0,
		x1-1,
		footerY1,
		0,
	)
	if footer == nil {
		creator.gui.Popup.Close(mergeRequestCreatorView)

		return fmt.Errorf(
			"create merge request creator footer: %w",
			err,
		)
	}

	content.Title = "FORM"

	content.FrameColor = gocui.ColorBlue
	content.TitleColor = gocui.ColorBlue

	footer.FrameColor = gocui.ColorBlue
	footer.TitleColor = gocui.ColorBlue

	content.Highlight = false
	footer.Highlight = false

	creator.renderForm(content)
	creator.renderFooter(footer)

	if _, err := creator.gui.g.SetViewOnTop(
		mergeRequestCreatorView,
	); err != nil {
		creator.close()

		return fmt.Errorf(
			"stack merge request creator: %w",
			err,
		)
	}

	if _, err := creator.gui.g.SetViewOnTop(
		mergeRequestCreatorFooterView,
	); err != nil {
		creator.close()

		return fmt.Errorf(
			"stack merge request creator footer: %w",
			err,
		)
	}

	if err := creator.gui.Popup.Focus(
		mergeRequestCreatorView,
	); err != nil {
		creator.close()

		return fmt.Errorf(
			"focus merge request creator: %w",
			err,
		)
	}

	creator.setFieldCursor(content)

	return nil
}

func (creator *mergeRequestCreator) renderForm(
	v *gocui.View,
) {
	v.Clear()

	fmt.Fprintln(v)
	fmt.Fprintln(
		v,
		"  Create a new GitLab Merge Request",
	)
	fmt.Fprintln(v)

	fields := []struct {
		field mergeRequestCreatorField
		label string
		value string
	}{
		{
			field: mergeRequestCreatorSourceBranchField,
			label: "Source Branch",
			value: creator.valueForField(
				mergeRequestCreatorSourceBranchField,
			),
		},
		{
			field: mergeRequestCreatorTargetBranchField,
			label: "Target Branch",
			value: creator.valueForField(
				mergeRequestCreatorTargetBranchField,
			),
		},
		{
			field: mergeRequestCreatorAssigneeField,
			label: "Assignee",
			value: creator.valueForField(
				mergeRequestCreatorAssigneeField,
			),
		},
		{
			field: mergeRequestCreatorReviewerField,
			label: "Reviewer",
			value: creator.valueForField(
				mergeRequestCreatorReviewerField,
			),
		},
		{
			field: mergeRequestCreatorTitleField,
			label: "Title",
			value: creator.valueForField(
				mergeRequestCreatorTitleField,
			),
		},
		{
			field: mergeRequestCreatorDescriptionField,
			label: "Description",
			value: creator.valueForField(
				mergeRequestCreatorDescriptionField,
			),
		},
	}

	for _, item := range fields {
		prefix := "  "

		if item.field == creator.selectedField {
			prefix = "> "
		}

		fmt.Fprintf(
			v,
			"%s%-18s [ %s ]\n",
			prefix,
			item.label,
			item.value,
		)
	}

	fmt.Fprintln(v)

	draftPrefix := "  "
	if creator.selectedField ==
		mergeRequestCreatorDraftField {
		draftPrefix = "> "
	}

	deletePrefix := "  "
	if creator.selectedField ==
		mergeRequestCreatorDeleteField {
		deletePrefix = "> "
	}

	squashPrefix := "  "
	if creator.selectedField ==
		mergeRequestCreatorSquashField {
		squashPrefix = "> "
	}

	draft := "[ ]"
	if creator.draft {
		draft = "[x]"
	}

	deleteSource := "[ ]"
	if creator.deleteSourceBranch {
		deleteSource = "[x]"
	}

	squash := "[ ]"
	if creator.squash {
		squash = "[x]"
	}

	fmt.Fprintf(
		v,
		"%s%s Draft\n",
		draftPrefix,
		draft,
	)

	fmt.Fprintf(
		v,
		"%s%s Delete source branch\n",
		deletePrefix,
		deleteSource,
	)

	fmt.Fprintf(
		v,
		"%s%s Squash commits\n",
		squashPrefix,
		squash,
	)
}

func (creator *mergeRequestCreator) renderFooter(
	v *gocui.View,
) {
	v.Clear()

	//fmt.Fprintln(v)
	fmt.Fprintln(
		v,
		"  [↑/↓] Navigate   [Enter] Edit   [Space] Toggle    [c] Create    [Esc] Cancel",
	)
	fmt.Fprintln(
		v,
		"  [Ctrl + s] Save Description   [Ctrl + F] Search Func",
	)

}

func (creator *mergeRequestCreator) valueForField(
	field mergeRequestCreatorField,
) string {
	switch field {
	case mergeRequestCreatorSourceBranchField:
		if creator.sourceBranch == "" {
			return "-"
		}

		return creator.sourceBranch

	case mergeRequestCreatorTargetBranchField:
		if creator.targetBranch == "" {
			return "-"
		}

		return creator.targetBranch

	case mergeRequestCreatorAssigneeField:
		if creator.assignee == "" {
			return "-"
		}

		return creator.assignee

	case mergeRequestCreatorReviewerField:
		if creator.reviewer == "" {
			return "-"
		}

		return creator.reviewer

	case mergeRequestCreatorTitleField:
		if creator.titleValue == "" {
			return ""
		}

		return creator.titleValue

	case mergeRequestCreatorDescriptionField:
		if creator.descriptionValue == "" {
			return ""
		}
		return creator.descriptionPreview()

	}

	return ""
}

func (creator *mergeRequestCreator) refreshForm() error {
	v, err := creator.gui.g.View(
		mergeRequestCreatorView,
	)
	if err != nil || v == nil {
		return nil
	}

	creator.renderForm(v)
	creator.setFieldCursor(v)

	return nil
}

func (creator *mergeRequestCreator) next() error {
	if creator.selectorOpen ||
		creator.editorOpen {
		return nil
	}

	if creator.selectedField >=
		mergeRequestCreatorSquashField {
		creator.selectedField =
			mergeRequestCreatorSourceBranchField
	} else {
		creator.selectedField++
	}

	return creator.refreshForm()
}

func (creator *mergeRequestCreator) previous() error {
	if creator.selectorOpen ||
		creator.editorOpen {
		return nil
	}

	if creator.selectedField <=
		mergeRequestCreatorSourceBranchField {
		creator.selectedField =
			mergeRequestCreatorSquashField
	} else {
		creator.selectedField--
	}

	return creator.refreshForm()
}

func (creator *mergeRequestCreator) setFieldCursor(
	v *gocui.View,
) {
	if v == nil {
		return
	}

	line := int(creator.selectedField) + 3

	v.SetCursor(0, line)
}

func (creator *mergeRequestCreator) edit() error {
	if creator.selectorOpen {
		return creator.selectItem()
	}

	if creator.editorOpen {
		return creator.saveEditor()
	}

	switch creator.selectedField {
	case mergeRequestCreatorSourceBranchField:
		return creator.openSelector(
			mergeRequestCreatorSourceBranchField,
			creator.branchItems(),
		)

	case mergeRequestCreatorTargetBranchField:
		return creator.openSelector(
			mergeRequestCreatorTargetBranchField,
			creator.branchItems(),
		)

	case mergeRequestCreatorAssigneeField:
		return creator.openMemberSelector(
			mergeRequestCreatorAssigneeField,
		)

	case mergeRequestCreatorReviewerField:
		return creator.openMemberSelector(
			mergeRequestCreatorReviewerField,
		)

	case mergeRequestCreatorTitleField,
		mergeRequestCreatorDescriptionField:
		return creator.openEditor(
			creator.selectedField,
		)

	case mergeRequestCreatorDraftField:
		creator.draft = !creator.draft

	case mergeRequestCreatorDeleteField:
		creator.deleteSourceBranch =
			!creator.deleteSourceBranch

	case mergeRequestCreatorSquashField:
		creator.squash = !creator.squash
	}

	return creator.refreshForm()
}

func (creator *mergeRequestCreator) toggle() error {
	if creator.selectorOpen ||
		creator.editorOpen {
		return nil
	}

	switch creator.selectedField {
	case mergeRequestCreatorDraftField:
		creator.draft = !creator.draft

	case mergeRequestCreatorDeleteField:
		creator.deleteSourceBranch =
			!creator.deleteSourceBranch

	case mergeRequestCreatorSquashField:
		creator.squash = !creator.squash

	default:
		return nil
	}

	return creator.refreshForm()
}

func (creator *mergeRequestCreator) branchItems() []string {
	if creator.gui.State == nil {
		return nil
	}

	branches := creator.gui.State.GetBranches()

	if branches == nil {
		return nil
	}

	items := make([]string, 0, len(branches.Items))

	for _, branch := range branches.Items {
		items = append(items, branch.Name)
	}

	return items
}

func (creator *mergeRequestCreator) openMemberSelector(
	field mergeRequestCreatorField,
) error {
	members, err := creator.loadMembers()
	if err != nil {
		return err
	}

	items := []string{"-"}

	for _, member := range members {
		if member.Name == "" {
			continue
		}

		items = append(
			items,
			member.Name,
		)
	}

	return creator.openSelector(
		field,
		items,
	)
}

func (creator *mergeRequestCreator) openSelector(
	field mergeRequestCreatorField,
	items []string,
) error {
	if len(items) == 0 {
		return nil
	}

	creator.selectorField = field
	creator.selectorItems = items
	creator.selectorSelected = 0
	creator.selectorOpen = true

	currentValue := creator.valueForField(field)

	for index, item := range items {
		if item == currentValue {
			creator.selectorSelected = index
			break
		}
	}

	width, height := creator.gui.g.Size()

	selectorWidth := width * 3 / 5

	if selectorWidth < 50 {
		selectorWidth = 50
	}

	if selectorWidth > width-6 {
		selectorWidth = width - 6
	}

	selectorHeight := len(items) + 2

	if selectorHeight > height-6 {
		selectorHeight = height - 6
	}

	if selectorHeight < 5 {
		selectorHeight = 5
	}

	x0 := (width - selectorWidth) / 2
	y0 := (height - selectorHeight) / 2
	x1 := x0 + selectorWidth
	y1 := y0 + selectorHeight

	selector, err := creator.gui.Popup.Create(
		mergeRequestCreatorSelectorView,
		"Select",
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		creator.selectorOpen = false

		return fmt.Errorf(
			"create merge request selector: %w",
			err,
		)
	}

	selector.Highlight = true
	selector.SelBgColor = gocui.NewRGBColor(
		35,
		35,
		35,
	)
	selector.SelFgColor = gocui.ColorWhite

	creator.renderSelector(selector)

	if _, err := creator.gui.g.SetViewOnTop(
		mergeRequestCreatorSelectorView,
	); err != nil {
		creator.closeSelector()

		return err
	}

	return creator.gui.Popup.Focus(
		mergeRequestCreatorSelectorView,
	)
}

func (creator *mergeRequestCreator) renderSelector(
	v *gocui.View,
) {
	v.Clear()

	for _, item := range creator.selectorItems {
		fmt.Fprintln(v, item)
	}

	creator.setSelectorCursor(v)
}

func (creator *mergeRequestCreator) setSelectorCursor(
	v *gocui.View,
) {
	if v == nil || len(creator.selectorItems) == 0 {
		return
	}

	selected := creator.selectorSelected

	if selected < 0 {
		selected = 0
	}

	if selected >= len(creator.selectorItems) {
		selected = len(creator.selectorItems) - 1
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
		originY =
			selected - visibleRows + 1
	}

	maxOriginY :=
		len(creator.selectorItems) - visibleRows

	if maxOriginY < 0 {
		maxOriginY = 0
	}

	if originY > maxOriginY {
		originY = maxOriginY
	}

	cursorY := selected - originY

	v.SetOrigin(
		0,
		originY,
	)

	v.SetCursor(
		0,
		cursorY,
	)
}

func (creator *mergeRequestCreator) selectorNext() error {
	if !creator.selectorOpen {
		return nil
	}

	if len(creator.selectorItems) == 0 {
		return nil
	}

	creator.selectorSelected++

	if creator.selectorSelected >=
		len(creator.selectorItems) {
		creator.selectorSelected = 0
	}

	v, err := creator.gui.g.View(
		mergeRequestCreatorSelectorView,
	)
	if err != nil || v == nil {
		return nil
	}

	creator.setSelectorCursor(v)

	return nil
}

func (creator *mergeRequestCreator) selectorPrevious() error {
	if !creator.selectorOpen {
		return nil
	}

	if len(creator.selectorItems) == 0 {
		return nil
	}

	creator.selectorSelected--

	if creator.selectorSelected < 0 {
		creator.selectorSelected =
			len(creator.selectorItems) - 1
	}

	v, err := creator.gui.g.View(
		mergeRequestCreatorSelectorView,
	)
	if err != nil || v == nil {
		return nil
	}

	creator.setSelectorCursor(v)

	return nil
}

func (creator *mergeRequestCreator) selectItem() error {
	if !creator.selectorOpen ||
		len(creator.selectorItems) == 0 {
		return nil
	}

	selected := creator.selectorItems[
		creator.selectorSelected,
	]

	switch creator.selectorField {
	case mergeRequestCreatorSourceBranchField:
		creator.sourceBranch = selected

	case mergeRequestCreatorTargetBranchField:
		creator.targetBranch = selected

	case mergeRequestCreatorAssigneeField:
		if selected == "-" {
			creator.assignee = ""
			creator.assigneeID = 0
		} else {
			creator.assignee = selected

			if err := creator.setMemberID(
				selected,
				true,
			); err != nil {
				return err
			}
		}

	case mergeRequestCreatorReviewerField:
		if selected == "-" {
			creator.reviewer = ""
			creator.reviewerID = 0
		} else {
			creator.reviewer = selected

			if err := creator.setMemberID(
				selected,
				false,
			); err != nil {
				return err
			}
		}
	}

	return creator.closeSelector()
}

func (creator *mergeRequestCreator) setMemberID(
	name string,
	assignee bool,
) error {
	members, err := creator.loadMembers()
	if err != nil {
		return err
	}

	for _, member := range members {
		if member.Name != name {
			continue
		}

		if assignee {
			creator.assigneeID = int(member.ID)
		} else {
			creator.reviewerID = int(member.ID)
		}

		return nil
	}

	return fmt.Errorf(
		"GitLab member %q was not found",
		name,
	)
}

func (creator *mergeRequestCreator) closeSelector() error {
	if !creator.selectorOpen {
		return nil
	}

	creator.closeSelectorSearch()

	creator.gui.Popup.Close(
		mergeRequestCreatorSelectorView,
	)

	creator.selectorOpen = false
	creator.selectorSelected = 0
	creator.selectorItems = nil

	if err := creator.gui.Popup.Focus(
		mergeRequestCreatorView,
	); err != nil {
		return err
	}

	return creator.refreshForm()
}

func (creator *mergeRequestCreator) openSelectorSearch() error {
	if !creator.selectorOpen ||
		creator.selectorSearchOpen {
		return nil
	}

	width, height := creator.gui.g.Size()

	modalWidth := width / 2

	if modalWidth < 40 {
		modalWidth = 40
	}

	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	modalHeight := 3

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2
	x1 := x0 + modalWidth
	y1 := y0 + modalHeight

	search, err := creator.gui.Popup.Create(
		mergeRequestCreatorSelectorSearchView,
		"Search",
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		return fmt.Errorf(
			"create selector search: %w",
			err,
		)
	}

	search.Highlight = true
	search.Editable = true
	search.Editor = gocui.DefaultEditor

	creator.selectorSearchOpen = true

	if _, err := creator.gui.g.SetViewOnTop(
		mergeRequestCreatorSelectorSearchView,
	); err != nil {
		creator.closeSelectorSearch()

		return err
	}

	return creator.gui.Popup.Focus(
		mergeRequestCreatorSelectorSearchView,
	)
}

func (creator *mergeRequestCreator) executeSelectorSearch() error {
	if !creator.selectorSearchOpen {
		return nil
	}

	search, err := creator.gui.g.View(
		mergeRequestCreatorSelectorSearchView,
	)
	if err != nil || search == nil {
		return nil
	}

	query := strings.TrimSpace(
		search.Buffer(),
	)

	if query == "" {
		return nil
	}

	index := searchItem(
		creator.selectorItems,
		query,
	)

	if index < 0 {
		return nil
	}

	creator.selectorSelected = index

	selectorView, err := creator.gui.g.View(
		mergeRequestCreatorSelectorView,
	)
	if err == nil && selectorView != nil {
		creator.setSelectorCursor(selectorView)
	}

	return creator.closeSelectorSearch()
}

func (creator *mergeRequestCreator) closeSelectorSearch() error {
	if !creator.selectorSearchOpen {
		return nil
	}

	creator.gui.Popup.Close(
		mergeRequestCreatorSelectorSearchView,
	)

	creator.selectorSearchOpen = false

	return creator.gui.Popup.Focus(
		mergeRequestCreatorSelectorView,
	)
}

func searchItem(
	items []string,
	query string,
) int {
	if len(items) == 0 {
		return -1
	}

	if !strings.ContainsAny(query, "*?") {
		for index, item := range items {
			if strings.EqualFold(item, query) {
				return index
			}
		}

		queryLower := strings.ToLower(query)

		for index, item := range items {
			if strings.Contains(
				strings.ToLower(item),
				queryLower,
			) {
				return index
			}
		}

		return -1
	}

	re, err := searchRegexp(query)
	if err != nil {
		return -1
	}

	for index, item := range items {
		if re.MatchString(item) {
			return index
		}
	}

	return -1
}

func searchRegexp(
	query string,
) (*regexp.Regexp, error) {
	var pattern strings.Builder

	for i := 0; i < len(query); i++ {
		switch query[i] {
		case '*':
			pattern.WriteString(".*")

		case '?':
			pattern.WriteByte('.')

		default:
			pattern.WriteByte(query[i])
		}
	}

	return regexp.Compile(
		"(?i)" + pattern.String(),
	)
}

func (creator *mergeRequestCreator) openEditor(
	field mergeRequestCreatorField,
) error {
	if creator.editorOpen {
		return nil
	}

	creator.editorField = field
	creator.editorOpen = true

	value := creator.valueForField(field)

	width, height := creator.gui.g.Size()

	modalWidth := width * 3 / 4

	if modalWidth < 60 {
		modalWidth = 60
	}

	if modalWidth > width-6 {
		modalWidth = width - 6
	}

	modalHeight := 5

	if field == mergeRequestCreatorDescriptionField {
		modalHeight = height * 3 / 5

		if modalHeight < 10 {
			modalHeight = 10
		}

		if modalHeight > height-6 {
			modalHeight = height - 6
		}
	}

	x0 := (width - modalWidth) / 2
	y0 := (height - modalHeight) / 2
	x1 := x0 + modalWidth
	y1 := y0 + modalHeight

	title := "Edit"

	if field == mergeRequestCreatorTitleField {
		title = "Edit Title"
	}

	if field == mergeRequestCreatorDescriptionField {
		title = "Edit Description"
	}

	editor, err := creator.gui.Popup.Create(
		mergeRequestCreatorEditorView,
		title,
		x0,
		y0,
		x1,
		y1,
	)
	if err != nil {
		creator.editorOpen = false

		return fmt.Errorf(
			"create merge request editor: %w",
			err,
		)
	}

	editor.Highlight = true
	editor.Editable = true
	editor.Editor = gocui.DefaultEditor

	editor.SelBgColor = gocui.NewRGBColor(
		35,
		35,
		35,
	)
	editor.SelFgColor = gocui.ColorWhite

	editor.TextArea.Clear()
	editor.TextArea.AutoWrap =
		field == mergeRequestCreatorDescriptionField
	editor.TextArea.AutoWrapWidth = editor.InnerWidth()

	editor.TextArea.TypeString(value)
	editor.RenderTextArea()

	if _, err := creator.gui.g.SetViewOnTop(
		mergeRequestCreatorEditorView,
	); err != nil {
		creator.closeEditor()

		return err
	}

	return creator.gui.Popup.Focus(
		mergeRequestCreatorEditorView,
	)
}

func (creator *mergeRequestCreator) saveEditor() error {
	if !creator.editorOpen {
		return nil
	}

	editor, err := creator.gui.g.View(
		mergeRequestCreatorEditorView,
	)
	if err != nil || editor == nil {
		return nil
	}

	value := strings.TrimSpace(
		editor.TextArea.GetUnwrappedContent(),
	)

	switch creator.editorField {
	case mergeRequestCreatorTitleField:
		creator.titleValue = value

	case mergeRequestCreatorDescriptionField:
		creator.descriptionValue = value
	}

	return creator.closeEditor()
}

func (creator *mergeRequestCreator) closeEditor() error {
	if !creator.editorOpen {
		return nil
	}

	creator.gui.Popup.Close(
		mergeRequestCreatorEditorView,
	)

	creator.editorOpen = false

	if err := creator.gui.Popup.Focus(
		mergeRequestCreatorView,
	); err != nil {
		return err
	}

	return creator.refreshForm()
}

func (creator *mergeRequestCreator) create() error {
	title := strings.TrimSpace(
		creator.titleValue,
	)

	if title == "" {
		return nil
	}

	if creator.sourceBranch == "" ||
		creator.targetBranch == "" {
		return nil
	}

	if creator.gui.State == nil {
		return nil
	}

	mergeRequests :=
		creator.gui.State.GetMergeRequests()

	if mergeRequests == nil ||
		creator.gui.GitLab == nil {
		return nil
	}

	creator.gui.State.GetCommandLogs().Add(
		fmt.Sprintf(
			"Creating merge request %s → %s...",
			creator.sourceBranch,
			creator.targetBranch,
		),
		false,
	)

	createdMergeRequest, err := mergeRequests.Create(
		creator.gui.GitLab,
		creator.sourceBranch,
		creator.targetBranch,
		title,
		strings.TrimSpace(
			creator.descriptionValue,
		),
		creator.draft,
		creator.deleteSourceBranch,
		creator.squash,
	)
	if err != nil {
		creator.gui.State.GetCommandLogs().Add(
			fmt.Sprintf(
				"Failed to create merge request: %s",
				err,
			),
			true,
		)

		_ = creator.gui.refreshCommandLogs()

		return creator.close()
	}

	if createdMergeRequest == nil {
		creator.gui.State.GetCommandLogs().Add(
			"Failed to create merge request: GitLab returned no merge request",
			true,
		)

		_ = creator.gui.refreshCommandLogs()

		return creator.close()
	}

	creator.gui.State.GetCommandLogs().Add(
		fmt.Sprintf(
			"Merge request !%d created successfully",
			createdMergeRequest.IID,
		),
		false,
	)

	if err := mergeRequests.Refresh(
		creator.gui.GitLab,
	); err != nil {
		creator.gui.State.GetCommandLogs().Add(
			fmt.Sprintf(
				"Failed to refresh merge requests after creating !%d: %s",
				createdMergeRequest.IID,
				err,
			),
			true,
		)

		_ = creator.gui.refreshCommandLogs()

		return creator.close()
	}

	selectedIndex := -1

	for index := range mergeRequests.Items {
		if mergeRequests.Items[index].IID ==
			createdMergeRequest.IID {
			selectedIndex = index
			break
		}
	}

	if selectedIndex == -1 {
		creator.gui.State.GetCommandLogs().Add(
			fmt.Sprintf(
				"Created merge request !%d was not found after refresh",
				createdMergeRequest.IID,
			),
			true,
		)

		_ = creator.gui.refreshCommandLogs()

		return creator.close()
	}

	if creator.gui.selectMergeRequest != nil {
		if err := creator.gui.selectMergeRequest(
			selectedIndex,
		); err != nil {
			creator.gui.State.GetCommandLogs().Add(
				fmt.Sprintf(
					"Failed to select created merge request !%d: %s",
					createdMergeRequest.IID,
					err,
				),
				true,
			)

			_ = creator.gui.refreshCommandLogs()

			return creator.close()
		}
	}

	return creator.close()
}

func (creator *mergeRequestCreator) title() string {
	return strings.TrimSpace(creator.titleValue)
}

func (creator *mergeRequestCreator) description() string {
	return strings.TrimSpace(creator.descriptionValue)
}

func (creator *mergeRequestCreator) discard() {
	creator.reset()
}



func (creator *mergeRequestCreator) editorEnter() error {
	if !creator.editorOpen {
		return nil
	}

	if creator.editorField ==
		mergeRequestCreatorDescriptionField {
		editor, err := creator.gui.g.View(
			mergeRequestCreatorEditorView,
		)
		if err != nil || editor == nil {
			return nil
		}

		editor.TextArea.TypeCharacter("\n")
		editor.RenderTextArea()

		return nil
	}

	return creator.saveEditor()
}


func (creator *mergeRequestCreator) descriptionPreview() string {
	if creator.descriptionValue == "" {
		return "-"
	}

	const maxLength = 80

	preview := strings.Join(
		strings.Fields(creator.descriptionValue),
		" ",
	)

	runes := []rune(preview)

	if len(runes) > maxLength {
		return string(runes[:maxLength]) + "..."
	}

	return preview
}
