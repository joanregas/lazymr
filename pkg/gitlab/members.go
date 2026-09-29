package gitlab

import (
	"fmt"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
	"lazymr/pkg/debug"
)

type Member struct {
	ID          int64
	Username    string
	Name        string
	AccessLevel int
}

type Members struct {
	Project string
	Items   []Member
}

func NewMembers(
	client *Client,
	project string,
) (*Members, error) {
	if project == "" {
		return nil, fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	members := &Members{
		Project: project,
	}

	if err := members.Refresh(client); err != nil {
		return nil, err
	}

	return members, nil
}

func (m *Members) Refresh(client *Client) error {
	if m.Project == "" {
		return fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	if client == nil {
		return fmt.Errorf(
			"GitLab client cannot be nil",
		)
	}

	debug.Log(
		"Refreshing project members: project=%q",
		m.Project,
	)

	m.Items = make([]Member, 0)

	page := 1

	for {
		options := &gitlabapi.ListProjectMembersOptions{
			ListOptions: gitlabapi.ListOptions{
				Page:    int64(page),
				PerPage: 100,
			},
		}

		members, response, err :=
			client.api.ProjectMembers.ListAllProjectMembers(
				m.Project,
				options,
			)
		if err != nil {
			return fmt.Errorf(
				"get members for GitLab project %q: %w",
				m.Project,
				err,
			)
		}

		for _, member := range members {
			if member == nil {
				continue
			}

			item := Member{
				ID:          member.ID,
				Username:    member.Username,
				Name:        member.Name,
				AccessLevel: int(member.AccessLevel),
			}

			m.Items = append(
				m.Items,
				item,
			)

			debug.Log(
				"Member: ID=%d username=%q name=%q role=%q",
				item.ID,
				item.Username,
				item.Name,
				item.Role(),
			)
		}

		if response == nil || response.NextPage == 0 {
			break
		}

		page = int(response.NextPage)
	}

	debug.Log(
		"Project members refresh complete: project=%q members=%d",
		m.Project,
		len(m.Items),
	)

	return nil
}

func (m Member) Role() string {
	switch m.AccessLevel {
	case 10:
		return "Guest"
	case 20:
		return "Reporter"
	case 30:
		return "Developer"
	case 40:
		return "Maintainer"
	case 50:
		return "Owner"
	default:
		return "Unknown"
	}
}

