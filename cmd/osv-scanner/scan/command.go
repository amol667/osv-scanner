package scan

import (
	"io"

	"github.com/google/osv-scanner/v2/cmd/osv-scanner/scan/image"
	"github.com/google/osv-scanner/v2/cmd/osv-scanner/scan/source"
	"github.com/urfave/cli/v3"
)

const sourceSubCommand = "source"

const DefaultSubcommand = sourceSubCommand

var Subcommands = []string{sourceSubCommand, "image"}

func Command(stdout, stderr io.Writer) *cli.Command {
	return &cli.Command{
		Name:        "scan",
		Usage:       "scans projects and container images for dependencies, and checks them against the OSV database.",
		Description: "Recursively scans projects and container images for dependencies and checks them against the OSV vulnerability database for known security issues.",
		Commands: []*cli.Command{
			source.Command(stdout, stderr),
			image.Command(stdout, stderr),
		},
	}
}
