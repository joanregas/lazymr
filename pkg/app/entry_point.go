package app

import (
	"fmt"
	"log"
	"os"
	"runtime"
  "errors"


	"lazymr/pkg/debug"
	"lazymr/pkg/config"
)

type BuildInfo struct {
	Commit      string
	Date        string
	Version     string
	BuildSource string
}

type cliArgs struct {
	PrintVersionInfo bool
	Debug            bool
}

func Start(buildInfo *BuildInfo) {
	cliArgs := parseCliArgsAndEnvVars()

	if cliArgs.Debug {
		debug.Reset()
		debug.Log("lazymr starting")
	}

	mergeBuildInfo(buildInfo)

	if cliArgs.PrintVersionInfo {
		printVersionInfo(buildInfo)
		debug.Log("Version: " + buildInfo.Version)

		return
	}

/*	app, err := NewApp(buildInfo, cliArgs)
	if err != nil {
		if cliArgs.Debug {
			debug.Log("NewApp failed: %v", err)
		}

		log.Fatal(err)
	}*/
	app, err := NewApp(buildInfo, cliArgs)
	if err != nil {
		if errors.Is(err, config.ErrConfigurationRequired) {
			return
		}
		if cliArgs.Debug {
			debug.Log("NewApp failed: %v", err)
		}
		log.Fatal(err)
	}

	if cliArgs.Debug {
		debug.Log("NewApp initialized successfully")
	}

	if err := app.Run(); err != nil {
		if cliArgs.Debug {
			debug.Log("App.Run failed: %v", err)
		}

		log.Fatal(err)
	}
}

func parseCliArgsAndEnvVars() *cliArgs {
	args := &cliArgs{}

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			args.PrintVersionInfo = true

		case "--debug", "-d":
			args.Debug = true

		default:
			fmt.Printf("unknown argument: %s\n", arg)
		}
	}

	if os.Getenv("DEBUG") == "TRUE" {
		args.Debug = true
	}

	return args
}

func mergeBuildInfo(buildInfo *BuildInfo) {
	if buildInfo.Version != "" {
		return
	}
	buildInfo.Version = "0.0.1"
}

func printVersionInfo(buildInfo *BuildInfo) {
	fmt.Printf(
		"commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s\n",
		buildInfo.Commit,
		buildInfo.Date,
		buildInfo.BuildSource,
		buildInfo.Version,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
