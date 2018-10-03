package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/containerd/console"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	pack "github.com/tonistiigi/buildkit-pack"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

func main() {
	app := cli.NewApp()
	app.Name = "pack-cli"
	app.UsageText = `pack-cli PATH | URL | -`
	app.Description = `
debug utility for invoking pack directly from client.
`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "buildkit-addr",
			Usage:  "buildkit daemon address",
			EnvVar: "BUILDKIT_HOST",
			Value:  appdefaults.Address,
		},
		cli.StringFlag{
			Name:  "progress",
			Usage: "Set type of progress (auto, plain, tty). Use plain to show container output",
			Value: "auto",
		},
		cli.StringFlag{
			Name:  "tag, t",
			Usage: "Name and optionally a tag in the 'name:tag' format",
		},
	})
	app.Action = action
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func action(clicontext *cli.Context) error {
	ctx := appcontext.Context()

	c, err := client.New(ctx, clicontext.String("buildkit-addr"), client.WithFailFast())
	if err != nil {
		return err
	}

	pipeR, pipeW := io.Pipe()
	solveOpt, err := newSolveOpt(clicontext, pipeW)
	if err != nil {
		return err
	}
	ch := make(chan *client.SolveStatus)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = c.Build(ctx, *solveOpt, "", pack.Build, ch)
		if err != nil {
			pipeW.CloseWithError(err)
		}
		return err
	})
	eg.Go(func() error {
		var c console.Console
		progressOpt := clicontext.String("progress")

		switch progressOpt {
		case "auto", "tty":
			cf, err := console.ConsoleFromFile(os.Stderr)
			if err != nil && progressOpt == "tty" {
				return err
			}
			c = cf
		case "plain":
		default:
			return errors.Errorf("invalid progress value : %s", progressOpt)
		}

		// not using shared context to not disrupt display but let is finish reporting errors
		return progressui.DisplaySolveStatus(context.TODO(), "", c, os.Stdout, ch)
	})
	eg.Go(func() error {
		if err := loadDockerTar(pipeR); err != nil {
			return err
		}
		return pipeR.Close()
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	logrus.Infof("Loaded the image %q to Docker.", clicontext.String("tag"))
	return nil
}

func newSolveOpt(clicontext *cli.Context, w io.WriteCloser) (*client.SolveOpt, error) {
	buildCtx := clicontext.Args().First()
	if buildCtx == "" {
		return nil, errors.New("please specify build context (e.g. \".\" for the current directory)")
	} else if buildCtx == "-" {
		return nil, errors.New("stdin not supported")
	}

	localDirs := map[string]string{
		"context": buildCtx,
	}

	frontendAttrs := map[string]string{}
	// if target := clicontext.String("target"); target != "" {
	// 	frontendAttrs["target"] = target
	// }

	return &client.SolveOpt{
		Exporter: "docker",
		ExporterAttrs: map[string]string{
			"name": clicontext.String("tag"),
		},
		ExporterOutput: w,
		LocalDirs:      localDirs,
		FrontendAttrs:  frontendAttrs,
	}, nil
}

func loadDockerTar(r io.Reader) error {
	// no need to use moby/moby/client here
	cmd := exec.Command("docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
