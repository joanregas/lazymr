package gitlab

import (
	"fmt"
	"io"
	"time"

	gitlabapi "gitlab.com/gitlab-org/api/client-go"
)

type PipelineJobs struct {
	PipelineID    int64
	ID            int64
	Name          string
	Status        string
	Stage         string
	Ref           string
	SHA           string
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

	Coverage float64
	User     string
	Runner   string
}

func (c *Client) GetPipelineJobs(
	project string,
	pipelineID int64,
) ([]PipelineJobs, error) {
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

	var result []PipelineJobs

	for {
		jobs, response, err := c.api.Jobs.ListPipelineJobs(
			project,
			pipelineID,
			options,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"get jobs for GitLab pipeline %d in project %q: %w",
				pipelineID,
				project,
				err,
			)
		}

		for _, job := range jobs {
			result = append(result, PipelineJobs{
				ID:            job.ID,
				Name:          job.Name,
				Status:        job.Status,
				Stage:         job.Stage,
				Ref:           job.Ref,
				SHA:           job.Pipeline.Sha,
				WebURL:        job.WebURL,
				Tag:           job.Tag,
				AllowFailure:  job.AllowFailure,
				FailureReason: job.FailureReason,
				CreatedAt:     job.CreatedAt,
				StartedAt:     job.StartedAt,
				FinishedAt:    job.FinishedAt,
				ErasedAt:      job.ErasedAt,
				Duration:      job.Duration,
				QueuedDuration: job.QueuedDuration,
				Coverage:      job.Coverage,
			})
		}

		if response == nil || response.CurrentPage >= response.TotalPages {
			break
		}

		options.Page = response.NextPage
	}

	return result, nil
}

func (c *Client) GetJob(
	project string,
	jobID int64,
) (*PipelineJobs, error) {
	if project == "" {
		return nil, fmt.Errorf("GitLab project path cannot be empty")
	}

	if jobID <= 0 {
		return nil, fmt.Errorf(
			"GitLab job ID must be greater than zero",
		)
	}

	job, _, err := c.api.Jobs.GetJob(project, int64(jobID))
	if err != nil {
		return nil, fmt.Errorf(
			"get GitLab job %d in project %q: %w",
			jobID,
			project,
			err,
		)
	}

	result := &PipelineJobs{
		ID:            job.ID,
		Name:          job.Name,
		Status:        job.Status,
		Stage:         job.Stage,
		Ref:           job.Ref,
		SHA:           job.Pipeline.Sha,
		WebURL:        job.WebURL,
		Tag:           job.Tag,
		AllowFailure:  job.AllowFailure,
		FailureReason: job.FailureReason,
		CreatedAt:     job.CreatedAt,
		StartedAt:     job.StartedAt,
		FinishedAt:    job.FinishedAt,
		ErasedAt:      job.ErasedAt,
		Duration:      job.Duration,
		QueuedDuration: job.QueuedDuration,
		Coverage:      job.Coverage,
	}

	return result, nil
}


func (c *Client) GetJobTrace(
	project string,
	jobID int64,
) (string, error) {
	if project == "" {
		return "", fmt.Errorf("GitLab project is empty")
	}

	if jobID <= 0 {
		return "", fmt.Errorf("invalid GitLab job ID: %d", jobID)
	}

	trace, _, err := c.api.Jobs.GetTraceFile(
		project,
		jobID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"get job %d trace: %w",
			jobID,
			err,
		)
	}

	data, err := io.ReadAll(trace)
	if err != nil {
		return "", fmt.Errorf(
			"read job %d trace: %w",
			jobID,
			err,
		)
	}

	return string(data), nil
}

// i dont want a general "retry pipeline" but x job
func (c *Client) RetryJob(
	project string,
	jobID int64,
) (*PipelineJobs, error) {
	if project == "" {
		return nil, fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	if jobID <= 0 {
		return nil, fmt.Errorf(
			"GitLab job ID must be greater than zero",
		)
	}

	job, _, err := c.api.Jobs.RetryJob(
		project,
		int64(jobID),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"retry GitLab job %d in project %q: %w",
			jobID,
			project,
			err,
		)
	}

	if job == nil {
		return nil, fmt.Errorf(
			"GitLab returned an empty job after retrying job %d",
			jobID,
		)
	}

	return &PipelineJobs{
		PipelineID:    job.Pipeline.ID,
		ID:            job.ID,
		Name:          job.Name,
		Status:        job.Status,
		Stage:         job.Stage,
		Ref:           job.Ref,
		SHA:           job.Pipeline.Sha,
		WebURL:        job.WebURL,
		Tag:           job.Tag,
		AllowFailure:  job.AllowFailure,
		FailureReason: job.FailureReason,
		CreatedAt:     job.CreatedAt,
		StartedAt:     job.StartedAt,
		FinishedAt:    job.FinishedAt,
		ErasedAt:      job.ErasedAt,
		Duration:      job.Duration,
		QueuedDuration: job.QueuedDuration,
		Coverage:      job.Coverage,
	}, nil
}



func (job *PipelineJobs) IsRunning() bool {
	switch job.Status {
	case "created",
		"pending",
		"preparing",
		"scheduled",
		"running",
		"waiting_for_callback",
		"waiting_for_resource",
		"canceling":
		return true
	default:
		return false
	}
}


// StartJob starts a manual GitLab CI/CD job.
func (c *Client) PlayJob(
	project string,
	jobID int64,
) (*PipelineJobs, error) {
	if project == "" {
		return nil, fmt.Errorf(
			"GitLab project path cannot be empty",
		)
	}

	if jobID <= 0 {
		return nil, fmt.Errorf(
		"GitLab job ID must be greater than zero",
		)
	}

	job, _, err := c.api.Jobs.PlayJob(
		project,
		jobID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"start GitLab job %d in project %q: %w",
			jobID,
			project,
			err,
		)
	}

	if job == nil {
		return nil, fmt.Errorf(
			"GitLab returned an empty job after starting job %d",
			jobID,
		)
	}

	return &PipelineJobs{
		PipelineID:    job.Pipeline.ID,
		ID:            job.ID,
		Name:          job.Name,
		Status:        job.Status,
		Stage:         job.Stage,
		Ref:           job.Ref,
		SHA:           job.Pipeline.Sha,
		WebURL:        job.WebURL,
		Tag:           job.Tag,
		AllowFailure:  job.AllowFailure,
		FailureReason: job.FailureReason,
		CreatedAt:     job.CreatedAt,
		StartedAt:     job.StartedAt,
		FinishedAt:     job.FinishedAt,
		ErasedAt:       job.ErasedAt,
		Duration:       job.Duration,
		QueuedDuration: job.QueuedDuration,
		Coverage:       job.Coverage,
	}, nil
}


