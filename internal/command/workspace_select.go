// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/opentofu/opentofu/internal/tfdiags"
	"github.com/posener/complete"
	"strings"
)

type WorkspaceSelectCommand struct {
	Meta
	LegacyName bool
}

func (c *WorkspaceSelectCommand) Run(args []string) int {
	var diags tfdiags.Diagnostics

	args = c.Meta.process(args)
	diags = envCommandHasWarning(c.LegacyName)
	if c.HasLegacyViewErrors(diags) {
		c.showDiagnostics(diags)
		return 1
	}

	var orCreate bool
	cmdFlags := c.Meta.defaultFlagSet("workspace select")
	c.Meta.varFlagSet(cmdFlags)
	cmdFlags.BoolVar(&orCreate, "or-create", false, "create workspace if it does not exist")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s\n", err.Error()))
		return 1
	}

	args = cmdFlags.Args()
	if len(args) != 1 {
		c.Ui.Error("Expected a single argument: NAME.\n")
		return cli.RunResultHelp
	}

	configPath, err := modulePath(args[1:])
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	backendConfig, backendDiags := c.loadBackendConfig(configPath)
	diags = diags.Append(backendDiags)
	if c.HasLegacyViewErrors(diags) {
		c.showDiagnostics(diags)
		return 1
	}

	current, isOverridden := c.WorkspaceOverridden()
	if isOverridden {
		c.Ui.Error(envIsOverriddenSelectError)
		return 1
	}

	// Load the encryption configuration
	enc, encDiags := c.EncryptionFromPath(configPath)
	diags = diags.Append(encDiags)
	if c.HasLegacyViewErrors(encDiags) {
		c.showDiagnostics(diags)
		return 1
	}

	// Load the backend
	b, backendDiags := c.Backend(&BackendOpts{
		Config: backendConfig,
	}, enc.State())
	diags = diags.Append(backendDiags)
	if c.HasLegacyViewErrors(backendDiags) {
		c.showDiagnostics(diags)
		return 1
	}

	// This command will not write state
	c.ignoreRemoteVersionConflict(b)

	name := args[0]
	if !validWorkspaceName(name) {
		c.Ui.Error(fmt.Sprintf(envInvalidName, name))
		return 1
	}

	states, err := b.Workspaces()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if name == current {
		// already using this workspace
		return 0
	}

	found := false
	for _, s := range states {
		if name == s {
			found = true
			break
		}
	}

	var newState bool

	if !found {
		if orCreate {
			_, err = b.StateMgr(name)
			if err != nil {
				c.Ui.Error(err.Error())
				return 1
			}
			newState = true
		} else {
			c.Ui.Error(fmt.Sprintf(envDoesNotExist, name))
			return 1
		}
	}

	err = c.SetWorkspace(name)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if newState {
		c.Ui.Output(c.Colorize().Color(fmt.Sprintf(
			strings.TrimSpace(envCreated), name)))
	} else {
		c.Ui.Output(
			c.Colorize().Color(
				fmt.Sprintf(envChanged, name),
			),
		)
	}

	c.showDiagnostics(diags)
	if c.HasLegacyViewErrors(diags) {
		return 1
	}

	return 0
}

func (c *WorkspaceSelectCommand) AutocompleteArgs() complete.Predictor {
	return completePredictSequence{
		c.completePredictWorkspaceName(),
		complete.PredictDirs(""),
	}
}

func (c *WorkspaceSelectCommand) AutocompleteFlags() complete.Flags {
	return nil
}

func (c *WorkspaceSelectCommand) Help() string {
	helpText := `
Usage: tofu [global options] workspace select [options] NAME

  Select a different OpenTofu workspace.

Options:

    -or-create=false    Create the OpenTofu workspace if it doesn't exist.

    -var 'foo=bar'      Set a value for one of the input variables in the root
                        module of the configuration. Use this option more than
                        once to set more than one variable.

    -var-file=filename  Load variable values from the given file, in addition
                        to the default files terraform.tfvars and *.auto.tfvars.
                        Use this option more than once to include more than one
                        variables file.
`
	return strings.TrimSpace(helpText)
}

func (c *WorkspaceSelectCommand) Synopsis() string {
	return "Select a workspace"
}
