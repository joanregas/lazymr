package gitlab

import (
	"fmt"
	"time"

	"lazymr/pkg/debug"
)

type Pipeline struct {
	ID        int64
	IID       int64
	Status    string
	Source    string
	Ref       string
	SHA       string
	WebURL    string
	UpdatedAt *time.Time
	CreatedAt *time.Time
}

type Pipelines struct {
	Project         string
	MergeRequestIID int
	Items           []Pipeline

	Jobs       []PipelineJobs
	Downstream []Bridge

	loadedDetails map[int64]bool
}

func NewPipelines(
	client *Client,
	project string,
	mergeRequestIID int,
) (*Pipelines, error) {
	if project == "" {
		return nil, fmt.Errorf("GitLab project path cannot be empty")
	}

	if mergeRequestIID <= 0 {
		return nil, fmt.Errorf(
			"GitLab merge request IID must be greater than zero",
		)
	}

	pipelines := &Pipelines{
		Project:         project,
		MergeRequestIID: mergeRequestIID,
		loadedDetails:   make(map[int64]bool),
	}

	if err := pipelines.Refresh(client); err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (p *Pipelines) Refresh(client *Client) error {
	if p.Project == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	if p.MergeRequestIID <= 0 {
		return fmt.Errorf(
			"GitLab merge request IID must be greater than zero",
		)
	}

	pipelines, _, err := client.api.MergeRequests.ListMergeRequestPipelines(
		p.Project,
		int64(p.MergeRequestIID),
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"get pipelines for GitLab merge request !%d in project %q: %w",
			p.MergeRequestIID,
			p.Project,
			err,
		)
	}

	p.Items = make([]Pipeline, 0, len(pipelines))

	// The pipeline list has changed, so previously loaded details
	// must not be considered valid anymore.
	p.Jobs = make([]PipelineJobs, 0)
	p.Downstream = make([]Bridge, 0)
	p.loadedDetails = make(map[int64]bool)

	debug.Log(
		"Found %d pipelines for MR !%d",
		len(pipelines),
		p.MergeRequestIID,
	)

	for _, pipeline := range pipelines {
		p.Items = append(p.Items, Pipeline{
			ID:        pipeline.ID,
			IID:       pipeline.IID,
			Status:    pipeline.Status,
			Source:    pipeline.Source,
			Ref:       pipeline.Ref,
			SHA:       pipeline.SHA,
			WebURL:    pipeline.WebURL,
			UpdatedAt: pipeline.UpdatedAt,
			CreatedAt: pipeline.CreatedAt,
		})
	}

	debug.Log(
		"Pipeline refresh complete: MR=!%d pipelines=%d",
		p.MergeRequestIID,
		len(p.Items),
	)

	return nil
}


func (p *Pipelines) LoadPipelineDetails(
	client *Client,
	projectPath string,
	pipelineID int64,
) error {
	if projectPath == "" {
		return fmt.Errorf("GitLab project path cannot be empty")
	}

	if pipelineID <= 0 {
		return fmt.Errorf(
			"GitLab pipeline ID must be greater than zero",
		)
	}

	if p.loadedDetails[pipelineID] {
		debug.Log(
			"Pipeline details already loaded: project=%q pipeline=%d",
			projectPath,
			pipelineID,
		)
		return nil
	}

	debug.Log(
		"Loading pipeline details: project=%q pipeline=%d",
		projectPath,
		pipelineID,
	)

	jobs, err := client.GetPipelineJobs(
		projectPath,
		pipelineID,
	)
	if err != nil {
		return fmt.Errorf(
			"get jobs for GitLab pipeline %d in project %q: %w",
			pipelineID,
			projectPath,
			err,
		)
	}

	debug.Log(
		"Pipeline %d: found %d jobs",
		pipelineID,
		len(jobs),
	)

	for i := range jobs {
		jobs[i].PipelineID = pipelineID

		debug.Log(
			"  Job: pipeline=%d ID=%d name=%q status=%s stage=%s",
			jobs[i].PipelineID,
			jobs[i].ID,
			jobs[i].Name,
			jobs[i].Status,
			jobs[i].Stage,
		)

		p.Jobs = append(p.Jobs, jobs[i])
	}

	bridges, err := client.GetPipelineBridges(
		projectPath,
		pipelineID,
	)
	if err != nil {
		return fmt.Errorf(
			"get downstream pipelines for GitLab pipeline %d in project %q: %w",
			pipelineID,
			projectPath,
			err,
		)
	}

	debug.Log(
		"Pipeline %d: found %d bridges",
		pipelineID,
		len(bridges),
	)

	for i := range bridges {
		bridges[i].PipelineID = pipelineID

		if bridges[i].DownstreamPipeline != nil {
			debug.Log(
				"  Bridge: pipeline=%d ID=%d name=%q status=%s downstream=%d",
				bridges[i].PipelineID,
				bridges[i].ID,
				bridges[i].Name,
				bridges[i].Status,
				bridges[i].DownstreamPipeline.ID,
			)
		} else {
			debug.Log(
				"  Bridge: pipeline=%d ID=%d name=%q status=%s downstream=<nil>",
				bridges[i].PipelineID,
				bridges[i].ID,
				bridges[i].Name,
				bridges[i].Status,
			)
		}

		p.Downstream = append(p.Downstream, bridges[i])
	}

	p.loadedDetails[pipelineID] = true

	debug.Log(
		"Pipeline details loaded: project=%q pipeline=%d jobs=%d bridges=%d",
		projectPath,
		pipelineID,
		len(jobs),
		len(bridges),
	)

	return nil
}

func (p *Pipelines) RefreshPipelineDetails(
	client *Client,
	projectPath string,
	pipelineID int64,
) error {
	delete(p.loadedDetails, pipelineID)

	filteredJobs := make([]PipelineJobs, 0, len(p.Jobs))

	for _, job := range p.Jobs {
		if job.PipelineID != pipelineID {
			filteredJobs = append(filteredJobs, job)
		}
	}

	p.Jobs = filteredJobs

	filteredDownstream := make([]Bridge, 0, len(p.Downstream))

	for _, bridge := range p.Downstream {
		if bridge.PipelineID != pipelineID {
			filteredDownstream = append(filteredDownstream, bridge)
		}
	}

	p.Downstream = filteredDownstream

	return p.LoadPipelineDetails(
		client,
		projectPath,
		pipelineID,
	)
}

