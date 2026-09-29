package oscommands

import "os"

// OSCommand provides the operating-system operations needed by the application.
type OSCommand struct{}

func NewOSCommand() *OSCommand {
	return &OSCommand{}
}

// GetCurrentDir returns the application's current working directory.
func (c *OSCommand) GetCurrentDir() (string, error) {
	return os.Getwd()
}

// ChangeDir changes the application's current working directory.
func (c *OSCommand) ChangeDir(path string) error {
	return os.Chdir(path)
}

// Exists reports whether the given path exists.
func (c *OSCommand) Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

