package bbir

import (
	"context"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "BACKLOG_HOST",
			Name:   "host, H",
			Usage:  "(Required) backlog host name. Ex: xxx.backlog.jp, xxx.backlog.com",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_PROJECT_KEY",
			Name:   "project, P",
			Usage:  "(Required) backlog project key.",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_API_KEY",
			Name:   "key, K",
			Usage:  "(Required) backlog api key.",
		},
		cli.StringFlag{
			EnvVar: "BACKLOG_LANG",
			Name:   "lang, l",
			Value:  "en",
			Usage:  "language setting. (ja or en)",
		},
		cli.BoolFlag{
			Name:  "progress, p",
			Usage: "show progress bar",
		},
		cli.BoolFlag{
			Name:  "check, c",
			Usage: "check mode (validation only)",
		},
	}
)

func NewCLI() *CLI {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = "A command line tool for bulk registering of Backlog issue."
	app.UsageText = Name + " [options] FILE_PATH"
	app.Version = FmtVersion()
	app.Flags = flags
	app.Action = action
	return &CLI{app}
}

type CLI struct {
	app *cli.App
}

func (cli *CLI) Run(argv []string) error {
	return cli.app.Run(argv)
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

	if cfg.Check {
		return nil
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
