package gitlab

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"lazymr/pkg/debug"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

type DownstreamPipeline struct {
	ID     int64
	Status string
	Ref    string
	SHA    string
	WebURL string

	ProjectID   int64
	ProjectPath string
}

type Bridge struct {
	PipelineID    int64
	ID            int64
	Name          string
	Status        string
	Stage         string
	Ref           string
	WebURL        string
	Tag           bool
	AllowFailure  bool
	FailureReason string

	CreatedAt      *time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
	ErasedAt       *time.Time
	Duration       float64
	QueuedDuration float64

	DownstreamPipeline *DownstreamPipeline
}

func pipelineProjectPath(webURL string) (string, error) {
	const marker = "/-/pipelines/"

	u, err := url.Parse(webURL)
	if err != nil {
		return "", fmt.Errorf(
			"parse GitLab pipeline URL %q: %w",
			webURL,
			err,
		)
	}

	index := strings.Index(u.Path, marker)

	if index == -1 {
		return "", fmt.Errorf(
			"invalid GitLab pipeline URL: %q",
			webURL,
		)
	}

	projectPath := strings.Trim(
		u.Path[:index],
		"/",
	)

	if projectPath == "" {
		return "", fmt.Errorf(
			"empty GitLab project path in pipeline URL: %q",
			webURL,
		)
	}

	return projectPath, nil
}

func (c *Client) GetPipelineBridges(
	project string,
	pipelineID int64,
) ([]Bridge, error) {
	if project == "" {
		return nil, fmt.Errorf("GitLab project path cannot be empty")
	}

	if pipelineID <= 0 {
		return nil, fmt.Errorf(
			"GitLab pipeline ID must be greater than zero",
		)
	}

	options := &gitlabapi.ListJobsOptions{
		ListOptions: gitlabapi.ListOptions{
			PerPage: 100,
		},
	}

	var result []Bridge
	seen := make(map[int64]bool)

	for {
		bridges, response, err := c.api.Jobs.ListPipelineBridges(
			project,
			pipelineID,
			options,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"get downstream pipelines for GitLab pipeline %d in project %q: %w",
				pipelineID,
				project,
				err,
			)
		}

		for _, bridge := range bridges {
			if bridge == nil {
				continue
			}

			if seen[bridge.ID] {
				debug.Log(
					"Skipping duplicate bridge: pipeline=%d bridgeID=%d",
					pipelineID,
					bridge.ID,
				)
				continue
			}

			seen[bridge.ID] = true

			projectPath := ""

			if bridge.DownstreamPipeline != nil {
				projectPath, err = pipelineProjectPath(
					bridge.DownstreamPipeline.WebURL,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"get downstream project path for pipeline %d: %w",
						bridge.DownstreamPipeline.ID,
						err,
					)
				}

				debug.Log(
					"Downstream pipeline: ID=%d project=%q",
					bridge.DownstreamPipeline.ID,
					projectPath,
				)
			}

			converted := convertBridge(
				bridge,
				projectPath,
			)

			if converted.DownstreamPipeline != nil {
				debug.Log(
					"Bridge: pipeline=%d bridgeID=%d name=%q status=%s downstreamID=%d downstreamProject=%q",
					pipelineID,
					converted.ID,
					converted.Name,
					converted.Status,
					converted.DownstreamPipeline.ID,
					converted.DownstreamPipeline.ProjectPath,
				)
			} else {
				debug.Log(
					"Bridge: pipeline=%d bridgeID=%d name=%q status=%s downstream=<nil>",
					pipelineID,
					converted.ID,
					converted.Name,
					converted.Status,
				)
			}

			result = append(result, converted)
		}

		if response == nil ||
			response.CurrentPage >= response.TotalPages {
			break
		}

		options.Page = response.NextPage
	}

	debug.Log(
		"Pipeline %d: found %d bridges",
		pipelineID,
		len(result),
	)

	return result, nil
}


func convertBridge(
	bridge *gitlabapi.Bridge,
	projectPath string,
) Bridge {
	result := Bridge{
		ID:             bridge.ID,
		Name:           bridge.Name,
		Status:         bridge.Status,
		Stage:          bridge.Stage,
		Ref:            bridge.Ref,
		WebURL:         bridge.WebURL,
		Tag:            bridge.Tag,
		AllowFailure:   bridge.AllowFailure,
		FailureReason:  bridge.FailureReason,
		CreatedAt:      bridge.CreatedAt,
		StartedAt:      bridge.StartedAt,
		FinishedAt:     bridge.FinishedAt,
		ErasedAt:       bridge.ErasedAt,
		Duration:       bridge.Duration,
		QueuedDuration: bridge.QueuedDuration,
	}

	if bridge.DownstreamPipeline != nil {
		result.DownstreamPipeline = &DownstreamPipeline{
			ID:          bridge.DownstreamPipeline.ID,
			Status:      bridge.DownstreamPipeline.Status,
			Ref:         bridge.DownstreamPipeline.Ref,
			SHA:         bridge.DownstreamPipeline.SHA,
			WebURL:      bridge.DownstreamPipeline.WebURL,
			ProjectID:   bridge.DownstreamPipeline.ProjectID,
			ProjectPath: projectPath,
		}
	}

	return result
}

