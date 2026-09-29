package gitlab

import (
	"fmt"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

type Branch struct {
	Name string
}

type Branches struct {
	Project string
	Items   []Branch
}

func NewBranches(
	client *Client,
	project string,
) (*Branches, error) {
	branches := &Branches{
		Project: project,
		Items:   make([]Branch, 0),
	}

	if err := branches.Refresh(client); err != nil {
		return nil, err
	}

	return branches, nil
}

func (b *Branches) Refresh(client *Client) error {
	if b.Project == "" {
		return fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	page := int64(1)
	perPage := int64(100)

	items := make([]Branch, 0)

	for {
		options := &gitlabapi.ListBranchesOptions{
			ListOptions: gitlabapi.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		branches, response, err :=
			client.api.Branches.ListBranches(
				b.Project,
				options,
			)

		if err != nil {
			return fmt.Errorf(
				"list GitLab branches for %q: %w",
				b.Project,
				err,
			)
		}

		for _, branch := range branches {
			if branch == nil {
				continue
			}

			items = append(
				items,
				Branch{
					Name: branch.Name,
				},
			)
		}

		if response == nil ||
			response.CurrentPage >= response.TotalPages {
			break
		}

		page++
	}

	b.Items = items

	return nil
}

