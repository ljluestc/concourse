package commands

import (
	"fmt"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/fly/commands/internal/displayhelpers"
	"github.com/concourse/concourse/fly/rc"
	"github.com/fatih/color"
)

type DestroyPipelineCommand struct {
	Pipeline        string `short:"p" long:"pipeline" required:"true" description:"Pipeline to destroy"`
	Team            string `short:"n" long:"team" description:"Name of the team that owns the pipeline"`
	SkipInteractive bool   `short:"n" long:"non-interactive" description:"Destroy the pipeline without confirmation"`
}

func (command *DestroyPipelineCommand) Execute(args []string) error {
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	var team = target.Team()
	if command.Team != "" {
		team, err = target.FindTeam(command.Team)
		if err != nil {
			return err
		}
	}

	pipelineName := command.Pipeline
	teamName := command.Team

	// Check if the pipeline exists before prompting
	pipelines, err := target.Team(teamName).ListPipelines()
	if err != nil {
		return err
	}

	pipelineExists := false
	for _, pipeline := range pipelines {
		if pipeline.Name == pipelineName {
			pipelineExists = true
			break
		}
	}

	if !pipelineExists {
		fmt.Printf("`%s` does not exist\n", pipelineName)
		return nil
	}

	if !command.SkipInteractive {
		displayhelpers.WarningSection("!!! this will remove all data for pipeline `" + pipelineName + "`")
		confirm := false
		err = displayhelpers.ConfirmInteractive(&confirm, "are you sure?")
		if err != nil {
			return err
		}

		if !confirm {
			displayhelpers.Failf("bailing out")
			return nil
		}
	}

	found, err := target.Team(teamName).DeletePipeline(pipelineName)
	if err != nil {
		return err
	}

	if found {
		fmt.Println(color.GreenString("`%s` deleted", pipelineName))
	} else {
		fmt.Println(color.RedString("`%s` does not exist", pipelineName))
	}

	return nil
}
