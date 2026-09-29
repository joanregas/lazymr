package gitlab

import (
	"fmt"
	"strings"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
	"lazymr/pkg/debug"
)

type MergeRequest struct {
	ID           int64
	IID          int
	Title        string
	State        string
	Author       string
	SourceBranch string
	TargetBranch string
	WebURL        string

	Description string

	Draft                       bool
	DetailedMergeStatus         string
	HasConflicts                bool
	BlockingDiscussionsResolved bool
	MergeWhenPipelineSucceeds   bool
	MergeError                  string

	Assignees []string
	Reviewers []string
	Labels    []string

	CreatedAt      string
	UpdatedAt      string
	UserNotesCount int
}

type MergeRequests struct {
	Project string
	Items   []MergeRequest
}

func NewMergeRequests(client *Client, project string) (*MergeRequests, error) {
	if project == "" {
		return nil, fmt.Errorf("GitLab project path cannot be empty")
	}

	mergeRequests := &MergeRequests{
		Project: project,
	}

	if err := mergeRequests.Refresh(client); err != nil {
		return nil, err
	}

	return mergeRequests, nil
}

// Refresh retrieves the current merge requests for the project.
func (m *MergeRequests) Refresh(client *Client) error {
	if m.Project == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	options := &gitlabapi.ListProjectMergeRequestsOptions{
		State: gitlabapi.Ptr("opened"),
	}

	mergeRequests, _, err := client.api.MergeRequests.ListProjectMergeRequests(
		m.Project,
		options,
	)
	if err != nil {
		return fmt.Errorf(
			"get merge requests for GitLab project %q: %w",
			m.Project,
			err,
		)
	}

	m.Items = make([]MergeRequest, 0, len(mergeRequests))

	for _, mergeRequest := range mergeRequests {
		m.Items = append(m.Items, MergeRequest{
			ID:           mergeRequest.ID,
			IID:          int(mergeRequest.IID),
			Title:        mergeRequest.Title,
			State:        mergeRequest.State,
			Author:       mergeRequest.Author.Name,
			SourceBranch: mergeRequest.SourceBranch,
			TargetBranch: mergeRequest.TargetBranch,
			WebURL:        mergeRequest.WebURL,
			Draft:                       mergeRequest.Draft,
			DetailedMergeStatus:         mergeRequest.DetailedMergeStatus,
			HasConflicts:                mergeRequest.HasConflicts,
			BlockingDiscussionsResolved: mergeRequest.BlockingDiscussionsResolved,
			MergeWhenPipelineSucceeds:   mergeRequest.MergeWhenPipelineSucceeds,
		})
	}

	return nil
}

// RefreshDetails retrieves the complete merge request from GitLab.
func (m *MergeRequest) RefreshDetails(
	client *Client,
	project string,
) error {
	if project == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	if m.IID <= 0 {
		return fmt.Errorf("GitLab merge request IID must be greater than zero")
	}

	debug.Log(
		"RefreshDetails: requesting project=%q iid=%d",
		project,
		m.IID,
	)

	mergeRequest, _, err := client.api.MergeRequests.GetMergeRequest(
		project,
		int64(m.IID),
		nil,
	)
	if err != nil {
		debug.Log(
			"RefreshDetails: GetMergeRequest failed: %v",
			err,
		)

		return fmt.Errorf(
			"get GitLab merge request !%d: %w",
			m.IID,
			err,
		)
	}

	debug.Log(
		"RefreshDetails: received MR !%d title=%q",
		mergeRequest.IID,
		mergeRequest.Title,
	)

	m.populateDetails(mergeRequest)

	return nil
}

func (m *MergeRequest) populateDetails(
	mergeRequest *gitlabapi.MergeRequest,
) {
	debug.Log("populateDetails: BEGIN")
	debug.Log(
		"API ID=%d IID=%d",
		mergeRequest.ID,
		mergeRequest.IID,
	)
	debug.Log(
		"API Title=%q",
		mergeRequest.Title,
	)
	debug.Log(
		"API State=%q Draft=%v",
		mergeRequest.State,
		mergeRequest.Draft,
	)
	debug.Log(
		"API Author=%+v",
		mergeRequest.Author,
	)
	debug.Log(
		"API SourceBranch=%q",
		mergeRequest.SourceBranch,
	)
	debug.Log(
		"API TargetBranch=%q",
		mergeRequest.TargetBranch,
	)
	debug.Log(
		"API WebURL=%q",
		mergeRequest.WebURL,
	)
	debug.Log(
		"API Description=%q",
		mergeRequest.Description,
	)
	debug.Log(
		"API DetailedMergeStatus=%q",
		mergeRequest.DetailedMergeStatus,
	)
	debug.Log(
		"API HasConflicts=%v",
		mergeRequest.HasConflicts,
	)
	debug.Log(
		"API BlockingDiscussionsResolved=%v",
		mergeRequest.BlockingDiscussionsResolved,
	)
	debug.Log(
		"API MergeWhenPipelineSucceeds=%v",
		mergeRequest.MergeWhenPipelineSucceeds,
	)
	debug.Log(
		"API MergeError=%q",
		mergeRequest.MergeError,
	)
	debug.Log(
		"API Assignees=%+v",
		mergeRequest.Assignees,
	)
	debug.Log(
		"API Reviewers=%+v",
		mergeRequest.Reviewers,
	)
	debug.Log(
		"API Labels=%+v",
		mergeRequest.Labels,
	)
	debug.Log(
		"API CreatedAt=%v",
		mergeRequest.CreatedAt,
	)
	debug.Log(
		"API UpdatedAt=%v",
		mergeRequest.UpdatedAt,
	)
	debug.Log(
		"API UserNotesCount=%d",
		mergeRequest.UserNotesCount,
	)

	m.ID = mergeRequest.ID
	m.IID = int(mergeRequest.IID)
	m.Title = mergeRequest.Title
	m.State = mergeRequest.State

	if mergeRequest.Author != nil {
		m.Author = mergeRequest.Author.Name
	}

	m.SourceBranch = mergeRequest.SourceBranch
	m.TargetBranch = mergeRequest.TargetBranch
	m.WebURL = mergeRequest.WebURL

	m.Description = mergeRequest.Description

	m.Draft = mergeRequest.Draft
	m.DetailedMergeStatus = mergeRequest.DetailedMergeStatus
	m.HasConflicts = mergeRequest.HasConflicts
	m.BlockingDiscussionsResolved =
		mergeRequest.BlockingDiscussionsResolved
	m.MergeWhenPipelineSucceeds =
		mergeRequest.MergeWhenPipelineSucceeds
	m.MergeError = mergeRequest.MergeError

	m.Assignees = make([]string, 0, len(mergeRequest.Assignees))

	for _, assignee := range mergeRequest.Assignees {
		if assignee == nil {
			continue
		}

		m.Assignees = append(
			m.Assignees,
			assignee.Name,
		)
	}

	m.Reviewers = make([]string, 0, len(mergeRequest.Reviewers))

	for _, reviewer := range mergeRequest.Reviewers {
		if reviewer == nil {
			continue
		}

		m.Reviewers = append(
			m.Reviewers,
			reviewer.Name,
		)
	}

	m.Labels = append(
		[]string(nil),
		mergeRequest.Labels...,
	)

	if mergeRequest.CreatedAt != nil {
		m.CreatedAt = mergeRequest.CreatedAt.String()
	}

	if mergeRequest.UpdatedAt != nil {
		m.UpdatedAt = mergeRequest.UpdatedAt.String()
	}

	m.UserNotesCount = int(mergeRequest.UserNotesCount)

	debug.Log("populateDetails: RESULT")
	debug.Log(
		"MR ID=%d IID=%d Title=%q State=%q",
		m.ID,
		m.IID,
		m.Title,
		m.State,
	)
	debug.Log(
		"MR Author=%q SourceBranch=%q TargetBranch=%q",
		m.Author,
		m.SourceBranch,
		m.TargetBranch,
	)
	debug.Log(
		"MR Description=%q",
		m.Description,
	)
	debug.Log(
		"MR Assignees=%+v",
		m.Assignees,
	)
	debug.Log(
		"MR Reviewers=%+v",
		m.Reviewers,
	)
	debug.Log(
		"MR Labels=%+v",
		m.Labels,
	)
	debug.Log(
		"MR CreatedAt=%q UpdatedAt=%q Notes=%d",
		m.CreatedAt,
		m.UpdatedAt,
		m.UserNotesCount,
	)
	debug.Log("populateDetails: END")
}

func (m *MergeRequests) Create(
	client *Client,
	sourceBranch string,
	targetBranch string,
	title string,
	description string,
	draft bool,
	removeSourceBranch bool,
	squash bool,
) (*MergeRequest, error) {
	if m.Project == "" {
		return nil, fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	if sourceBranch == "" {
		return nil, fmt.Errorf(
			"GitLab merge request source branch cannot be empty",
		)
	}

	if targetBranch == "" {
		return nil, fmt.Errorf(
			"GitLab merge request target branch cannot be empty",
		)
	}

	if title == "" {
		return nil, fmt.Errorf(
			"GitLab merge request title cannot be empty",
		)
	}

	if client == nil {
		return nil, fmt.Errorf(
			"GitLab client cannot be nil",
		)
	}

	if draft {
		title = ensureDraftTitle(title)
	}

	options := &gitlabapi.CreateMergeRequestOptions{
		Title:              gitlabapi.Ptr(title),
		Description:        gitlabapi.Ptr(description),
		SourceBranch:       gitlabapi.Ptr(sourceBranch),
		TargetBranch:       gitlabapi.Ptr(targetBranch),
		RemoveSourceBranch: gitlabapi.Ptr(removeSourceBranch),
		Squash:             gitlabapi.Ptr(squash),
	}

	debug.Log(
		"CreateMergeRequest: project=%q source=%q target=%q title=%q draft=%v remove_source_branch=%v squash=%v",
		m.Project,
		sourceBranch,
		targetBranch,
		title,
		draft,
		removeSourceBranch,
		squash,
	)

	mergeRequest, _, err :=
		client.api.MergeRequests.CreateMergeRequest(
			m.Project,
			options,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"create GitLab merge request %q -> %q: %w",
			sourceBranch,
			targetBranch,
			err,
		)
	}

	if mergeRequest == nil {
		return nil, fmt.Errorf(
			"GitLab returned an empty merge request",
		)
	}

	result := &MergeRequest{}
	result.populateDetails(mergeRequest)

	debug.Log(
		"CreateMergeRequest: created MR !%d title=%q",
		result.IID,
		result.Title,
	)

	return result, nil
}

func ensureDraftTitle(title string) string {
	lowerTitle := strings.ToLower(strings.TrimSpace(title))

	if strings.HasPrefix(lowerTitle, "draft:") ||
		strings.HasPrefix(lowerTitle, "[draft]") ||
		strings.HasPrefix(lowerTitle, "(draft)") {
		return title
	}

	return "Draft: " + title
}

func (m *MergeRequest) Close(
	client *Client,
	project string,
) error {
	if client == nil {
		return fmt.Errorf(
			"GitLab client cannot be nil",
		)
	}

	if project == "" {
		return fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	if m.IID <= 0 {
		return fmt.Errorf(
			"GitLab merge request IID must be greater than zero",
		)
	}

	debug.Log(
		"CloseMergeRequest: project=%q iid=%d",
		project,
		m.IID,
	)

	options := &gitlabapi.UpdateMergeRequestOptions{
		StateEvent: gitlabapi.Ptr("close"),
	}

	mergeRequest, _, err :=
		client.api.MergeRequests.UpdateMergeRequest(
			project,
			int64(m.IID),
			options,
		)
	if err != nil {
		return fmt.Errorf(
			"close GitLab merge request !%d: %w",
			m.IID,
			err,
		)
	}

	if mergeRequest == nil {
		return fmt.Errorf(
			"GitLab returned an empty merge request after closing !%d",
			m.IID,
		)
	}

	m.populateDetails(mergeRequest)

	debug.Log(
		"CloseMergeRequest: closed MR !%d",
		m.IID,
	)

	return nil
}

func (s *GitlabState) SelectMergeRequest(
	client *Client,
	index int,
) error {
	if s.MergeRequests == nil {
		return fmt.Errorf(
			"merge requests are not loaded",
		)
	}

	if index < 0 || index >= len(s.MergeRequests.Items) {
		return fmt.Errorf(
			"merge request index out of range: %d",
			index,
		)
	}

	selectedMergeRequest := &s.MergeRequests.Items[index]

	files, err := NewFiles(
		client,
		s.MergeRequests.Project,
		selectedMergeRequest.IID,
	)
	if err != nil {
		return fmt.Errorf(
			"load merge request files !%d: %w",
			selectedMergeRequest.IID,
			err,
		)
	}

	pipelines, err := NewPipelines(
		client,
		s.MergeRequests.Project,
		selectedMergeRequest.IID,
	)
	if err != nil {
		return fmt.Errorf(
			"load merge request pipelines !%d: %w",
			selectedMergeRequest.IID,
			err,
		)
	}

	s.Files = files
	s.Pipelines = pipelines
	s.SetSelectedMergeRequest(index)

	return nil
}


func (m *MergeRequest) Merge(
	client *Client,
	project string,
) error {
	if client == nil {
		return fmt.Errorf(
			"GitLab client cannot be nil",
		)
	}

	if project == "" {
		return fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	if m.IID <= 0 {
		return fmt.Errorf(
			"GitLab merge request IID must be greater than zero",
		)
	}

	debug.Log(
		"MergeMergeRequest: project=%q iid=%d",
		project,
		m.IID,
	)

	mergeRequest, _, err :=
		client.api.MergeRequests.AcceptMergeRequest(
			project,
			int64(m.IID),
			&gitlabapi.AcceptMergeRequestOptions{},
		)
	if err != nil {
		return fmt.Errorf(
			"merge GitLab merge request !%d: %w",
			m.IID,
			err,
		)
	}

	if mergeRequest == nil {
		return fmt.Errorf(
			"GitLab returned an empty merge request after merging !%d",
			m.IID,
		)
	}

	m.populateDetails(mergeRequest)

	debug.Log(
		"MergeMergeRequest: merged MR !%d",
		m.IID,
	)

	return nil
}

