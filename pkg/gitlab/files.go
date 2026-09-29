package gitlab

import (
	"fmt"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

type Files struct {
	Project       string
	MergeRequestIID int
	Items         []File
}

func NewFiles(
	client *Client,
	project string,
	mergeRequestIID int,
) (*Files, error) {
	if project == "" {
		return nil, fmt.Errorf("GitLab project path cannot be empty")
	}

	if mergeRequestIID <= 0 {
		return nil, fmt.Errorf(
			"GitLab merge request IID must be greater than zero",
		)
	}

	files := &Files{
		Project:        project,
		MergeRequestIID: mergeRequestIID,
	}

	if err := files.Refresh(client); err != nil {
		return nil, err
	}

	return files, nil
}

// Refresh retrieves the files changed by the merge request.
func (f *Files) Refresh(client *Client) error {
	if f.Project == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	if f.MergeRequestIID <= 0 {
		return fmt.Errorf(
			"GitLab merge request IID must be greater than zero",
		)
	}

	options := &gitlabapi.ListMergeRequestDiffsOptions{}

	mergeRequestDiffs, _, err := client.api.MergeRequests.ListMergeRequestDiffs(
		f.Project,
		int64(f.MergeRequestIID),
		options,
	)
	if err != nil {
		return fmt.Errorf(
			"get files for GitLab merge request !%d in project %q: %w",
			f.MergeRequestIID,
			f.Project,
			err,
		)
	}

	f.Items = make([]File, 0, len(mergeRequestDiffs))

	for _, diff := range mergeRequestDiffs {
		f.Items = append(f.Items, File{
			OldPath:       diff.OldPath,
			NewPath:       diff.NewPath,
			AMode:         diff.AMode,
			BMode:         diff.BMode,
			Diff:          diff.Diff,
			NewFile:       diff.NewFile,
			RenamedFile:   diff.RenamedFile,
			DeletedFile:   diff.DeletedFile,
			GeneratedFile: diff.GeneratedFile,
		})
	}

	return nil
}
