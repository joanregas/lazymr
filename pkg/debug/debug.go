package debug

import (
	"fmt"
	"os"
)

const logFile = "/tmp/lazymr-debug.log"

func Log(format string, args ...any) {
	f, err := os.OpenFile(
		logFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintf(f, format+"\n", args...)
}

func Reset() {
	_ = os.WriteFile(logFile, nil, 0644)
}
