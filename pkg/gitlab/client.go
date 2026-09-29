package gitlab

import (
	"fmt"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"

	"lazymr/pkg/config"
)

type Client struct {
	api *gitlabapi.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
	host := cfg.Hosts[cfg.Host]

	apiHost := host.APIHost
	if apiHost == "" {
		apiHost = cfg.Host
	}

	apiProtocol := host.APIProtocol
	if apiProtocol == "" {
		apiProtocol = "https"
	}

	baseURL := fmt.Sprintf(
		"%s://%s/api/v4/",
		apiProtocol,
		apiHost,
	)

	api, err := gitlabapi.NewClient(
		host.Token,
		gitlabapi.WithBaseURL(baseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("create GitLab client: %w", err)
	}

	return &Client{
		api: api,
	}, nil
}

func (c *Client) GetProject(project string) (*gitlabapi.Project, error) {
	p, _, err := c.api.Projects.GetProject(project, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"get GitLab project %q: %w",
			project,
			err,
		)
	}

	return p, nil
}

func (c *Client) GetProjectByID(
	projectID int64,
) (*gitlabapi.Project, error) {
	if projectID <= 0 {
		return nil, fmt.Errorf(
			"GitLab project ID must be greater than zero",
		)
	}

	project, _, err := c.api.Projects.GetProject(
		projectID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get GitLab project %d: %w",
			projectID,
			err,
		)
	}

	return project, nil
}
