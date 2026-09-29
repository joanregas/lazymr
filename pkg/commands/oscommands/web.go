package oscommands

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (c *OSCommand) OpenURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	var command string
	var args []string

	switch runtime.GOOS {
	case "linux":
		command = "xdg-open"
		args = []string{url}

	case "darwin":
		command = "open"
		args = []string{url}

	case "windows":
		command = "rundll32"
		args = []string{
			"url.dll,FileProtocolHandler",
			url,
		}

	default:
		return fmt.Errorf(
			"unsupported operating system: %s",
			runtime.GOOS,
		)
	}

	if err := exec.Command(command, args...).Start(); err != nil {
		return fmt.Errorf(
			"open URL %q: %w",
			url,
			err,
		)
	}

	return nil
}

