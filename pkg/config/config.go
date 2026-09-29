package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const placeholderToken = "PLACE_HOLDER"

type Config struct {
	GitProtocol       string                `yaml:"git_protocol"`
	Editor            string                `yaml:"editor"`
	Browser           string                `yaml:"browser"`
	GlamourStyle      string                `yaml:"glamour_style"`
	CheckUpdate       bool                  `yaml:"check_update"`
	DisplayHyperlinks bool                  `yaml:"display_hyperlinks"`
	Host              string                `yaml:"host"`
	NoPrompt          bool                  `yaml:"no_prompt"`
	Hosts             map[string]HostConfig `yaml:"hosts"`
}

type HostConfig struct {
	Token       string `yaml:"token"`
	APIHost     string `yaml:"api_host"`
	APIProtocol string `yaml:"api_protocol"`
	GitProtocol string `yaml:"git_protocol"`
	User        string `yaml:"user"`
}

type State struct {
	RecentRepos []string `yaml:"recentrepos"`
}

var ErrConfigurationRequired = errors.New(
	"Lazy MR configuration required",
)

func Load() (*Config, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", fmt.Errorf(
			"get home directory: %w",
			err,
		)
	}

	lazyMRConfig := filepath.Join(
		home,
		".config",
		"lazymr",
		"config.yml",
	)

	glabConfig := filepath.Join(
		home,
		".config",
		"glab-cli",
		"config.yml",
	)

	if _, err := os.Stat(lazyMRConfig); err == nil {
		cfg, err := readConfig(lazyMRConfig)
		if err == nil {
			if err := validateConfig(cfg); err == nil {
				return cfg, lazyMRConfig, nil
			}
		}
	}

	if _, err := os.Stat(glabConfig); err == nil {
		cfg, err := readConfig(glabConfig)
		if err == nil {
			if err := validateConfig(cfg); err == nil {
				if err := writeConfig(
					lazyMRConfig,
					cfg,
				); err != nil {
					return nil, lazyMRConfig, err
				}

				return cfg, lazyMRConfig, nil
			}
		}
	}

	cfg := newDefaultConfig()

	if err := writeConfig(lazyMRConfig, cfg); err != nil {
		return nil, lazyMRConfig, err
	}

	return nil, lazyMRConfig, ErrConfigurationRequired
}

func newDefaultConfig() *Config {
	return &Config{
		GitProtocol:       "ssh",
		Editor:            "nvim",
		Browser:           "",
		GlamourStyle:      "",
		CheckUpdate:       true,
		DisplayHyperlinks: true,
		Host:              "gitlab.example.com",
		NoPrompt:          false,
		Hosts: map[string]HostConfig{
			"gitlab.example.com": {
				Token:       placeholderToken,
				APIHost:     "gitlab.example.com",
				APIProtocol: "https",
				GitProtocol: "ssh",
				User:        "",
			},
		},
	}
}

func readConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"read config %s: %w",
			path,
			err,
		)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf(
			"parse config %s: %w",
			path,
			err,
		)
	}

	return &cfg, nil
}

func writeConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf(
			"create config directory %s: %w",
			dir,
			err,
		)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf(
			"serialize config: %w",
			err,
		)
	}

	// The configuration contains GitLab tokens.
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf(
			"write config %s: %w",
			path,
			err,
		)
	}

	return nil
}

func LoadState() (*State, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", fmt.Errorf(
			"get home directory: %w",
			err,
		)
	}

	lazyMRState := filepath.Join(
		home,
		".local",
		"state",
		"lazymr",
		"state.yml",
	)

	lazyGitState := filepath.Join(
		home,
		".local",
		"state",
		"lazygit",
		"state.yml",
	)

	// Use our own state if it already exists.
	if _, err := os.Stat(lazyMRState); err == nil {
		state, err := readState(lazyMRState)
		if err != nil {
			return nil, lazyMRState, err
		}

		return state, lazyMRState, nil
	}

	// Bootstrap from LazyGit if available.
	if _, err := os.Stat(lazyGitState); err == nil {
		state, err := readState(lazyGitState)
		if err != nil {
			return nil, lazyGitState, err
		}

		if err := writeState(lazyMRState, state); err != nil {
			return nil, lazyMRState, err
		}

		return state, lazyMRState, nil
	}

	// No existing state. Create ours from scratch.
	state := &State{}

	if err := writeState(lazyMRState, state); err != nil {
		return nil, lazyMRState, err
	}

	return state, lazyMRState, nil
}

func SaveState(state *State) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf(
			"get home directory: %w",
			err,
		)
	}

	path := filepath.Join(
		home,
		".local",
		"state",
		"lazymr",
		"state.yml",
	)

	return writeState(path, state)
}

func readState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"read state %s: %w",
			path,
			err,
		)
	}

	var state State

	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf(
			"parse state %s: %w",
			path,
			err,
		)
	}

	return &state, nil
}

func writeState(path string, state *State) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf(
			"create state directory %s: %w",
			dir,
			err,
		)
	}

	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf(
			"serialize state: %w",
			err,
		)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf(
			"write state %s: %w",
			path,
			err,
		)
	}

	return nil
}

func (s *State) AddRecentRepo(path string) {
	path = filepath.Clean(path)

	updated := make([]string, 0, len(s.RecentRepos)+1)

	// Current repository goes first.
	updated = append(updated, path)

	// Keep existing repositories, excluding the current one.
	for _, repo := range s.RecentRepos {
		repo = filepath.Clean(repo)

		if repo == path {
			continue
		}

		updated = append(updated, repo)
	}

	s.RecentRepos = updated
}

func validateConfig(cfg *Config) error {
	if len(cfg.Hosts) == 0 {
		return errors.New(
			"no GitLab hosts configured",
		)
	}

	if cfg.Host != "" {
		if host, ok := cfg.Hosts[cfg.Host]; ok {
			if host.Token != "" &&
				host.Token != placeholderToken {
				return nil
			}
		}
	}

	for hostname, host := range cfg.Hosts {
		if host.Token == "" ||
			host.Token == placeholderToken {
			continue
		}

		cfg.Host = hostname

		return nil
	}

	return errors.New(
		"no authenticated GitLab host configured",
	)
}


func (s *State) GetRecentRepo() (string, bool) {
	if s == nil || len(s.RecentRepos) == 0 {
		return "", false
	}

	return filepath.Clean(s.RecentRepos[0]), true
}
