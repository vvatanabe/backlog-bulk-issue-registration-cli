package main

import (
	"os"

	"context"

	"github.com/urfave/cli"
	. "github.com/vvatanabe/backlog-bulk-issue-registration-cli/bbir"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "BACKLOG_HOST",
			Name:   "host",
			Usage:  "(Required) backlog host name. Ex: xxx.backlog.jp, xxx.backlog.com",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_PROJECT_KEY",
			Name:   "project",
			Usage:  "(Required) backlog project key.",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_API_KEY",
			Name:   "key",
			Usage:  "(Required) backlog api key.",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_IMPORT_FILE",
			Name:   "file",
			Usage:  "import file path. If not specified, it will listen on standard input",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_LANG",
			Name:   "lang",
			Value:  "ja",
			Usage:  "language setting. (ja or en)",
		},
		cli.BoolFlag{
			Name:  "progress",
			Usage: "show progress bar",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = "A command line tool for bulk registering of Backlog issue."
	app.UsageText = "backlog-bulk-issue-registration [options]"
	app.Version = FmtVersion()
	app.Flags = flags
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) error {

	cfg := NewConfig(c)

	module, err := NewModule(cfg)
	if err != nil {
		return cli.NewExitError(err.Error(), UnexpectedError)
	}

	msgs := module.GetMessages()

	if err := module.GetProjectRepository().Prefetch(context.Background()); err != nil {
		return cli.NewExitError(msgs.FailedToCallBacklogAPI(err), BacklogAPIRequestError)
	}

	var file *os.File
	if cfg.HasFile() {
		file, err = os.Open(cfg.File)
		if err != nil {
			return cli.NewExitError(msgs.CanNotOpenFile(cfg.File), FileOpenError)
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}

	r := NewCSVReader(file)
	size, lines, err := r.ReadAll()
	if err != nil {
		name := file.Name()
		if name == "" {
			name = "stdin"
		}
		return cli.NewExitError(msgs.CanNotReadFile(name, err), FileReadError)
	}

	commands, err := module.GetCommandConverter().Convert(lines, func() []Callback {
		if cfg.Progress {
			return CreateProgressCallbacks(size, msgs.TagOfValidationProgressBar())
		}
		return []Callback{}
	}()...)
	if err != nil {
		return cli.NewExitError(err.Error(), ValidationIssueError)
	}

	if err := module.GetBulkCommandExecutor().Do(commands, func() []Callback {
		if cfg.Progress {
			return CreateProgressCallbacks(size, msgs.TagOfRegistrationProgressBar())
		}
		return []Callback{}
	}()...); err != nil {
		return cli.NewExitError(err.Error(), RegistrationIssueError)
	}

	return nil
}

func CreateProgressCallbacks(total int, tag string) []Callback {
	bar := pb.New(total).Prefix(tag)
	return []Callback{
		Before(func() {
			bar.Start()
		}),
		Each(func() {
			bar.Increment()
		}),
		After(func() {
			bar.Finish()
		}),
	}
}
