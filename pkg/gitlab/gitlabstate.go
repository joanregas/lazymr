package gitlab


type GitlabState struct {
	Repository    *Repository
	Branches      *Branches
	MergeRequests *MergeRequests
	Files         *Files
	Members       *Members
	Pipelines     *Pipelines

	CommandLogs *CommandLogs

	SelectedMergeRequest int
	SelectedFilePath     string
}

func NewGitlabState(
	repository *Repository,
	branches *Branches,
	mergeRequests *MergeRequests,
	members *Members,
	files *Files,
	pipelines *Pipelines,
) *GitlabState {
	return &GitlabState{
		Repository:           repository,
		Branches: 						branches,
		MergeRequests:        mergeRequests,
		Files:                files,
		Members:              members,
		Pipelines:            pipelines,
		CommandLogs:          NewCommandLogs(),
		SelectedMergeRequest: 0,
		SelectedFilePath:     "",
	}
}

func (s *GitlabState) GetRepository() *Repository {
	return s.Repository
}

func (s *GitlabState) GetBranches() *Branches { 
	return s.Branches 
}

func (s *GitlabState) GetMergeRequests() *MergeRequests {
	return s.MergeRequests
}

func (s *GitlabState) GetFiles() *Files {
	return s.Files
}

func (s *GitlabState) GetPipelines() *Pipelines {
	return s.Pipelines
}

func (s *GitlabState) GetSelectedMergeRequest() int {
	return s.SelectedMergeRequest
}

func (s *GitlabState) SetSelectedMergeRequest(index int) {
	s.SelectedMergeRequest = index
	s.SelectedFilePath = ""
}

func (s *GitlabState) GetSelectedFilePath() string {
	return s.SelectedFilePath
}

func (s *GitlabState) SetSelectedFilePath(path string) {
	s.SelectedFilePath = path
}

func (s *GitlabState) GetSelectedFile() *File {
	if s.Files == nil || s.SelectedFilePath == "" {
		return nil
	}

	for i := range s.Files.Items {
		file := &s.Files.Items[i]

		if file.NewPath == s.SelectedFilePath {
			return file
		}
	}

	return nil
}

func (s *GitlabState) GetCommandLogs() *CommandLogs { 
	return s.CommandLogs 
}


