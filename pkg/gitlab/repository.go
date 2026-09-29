package gitlab

import (
	"fmt"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

type Repository struct {
	ID                int64
	Name              string
	Path              string
	PathWithNamespace string
	Namespace         string
	WebURL            string
	SSHURL            string
	HTTPURL           string
	DefaultBranch     string
	Description       string
	project 					string
}

func NewRepository(client *Client, project string) (*Repository, error) {
	repository := &Repository{
		project: project,
}

	//if err := repository.Refresh(client, project); err != nil {
	if err := repository.Refresh(client); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *Repository) Refresh(client *Client) error {
	if r == nil {
		return fmt.Errorf("GitLab repository is nil")
	}

	if r.project == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	p, _, err := client.api.Projects.GetProject(
		r.project,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"get GitLab project %q: %w",
			r.project,
			err,
		)
	}

	r.populate(p)

	return nil
}



func (r *Repository) populate(project *gitlabapi.Project) {
	r.ID = project.ID
	r.Name = project.Name
	r.Path = project.Path
	r.PathWithNamespace = project.PathWithNamespace
	r.WebURL = project.WebURL
	r.SSHURL = project.SSHURLToRepo
	r.HTTPURL = project.HTTPURLToRepo
	r.DefaultBranch = project.DefaultBranch
	r.Description = project.Description

	if project.Namespace != nil {
		r.Namespace = project.Namespace.FullPath
	}
}


