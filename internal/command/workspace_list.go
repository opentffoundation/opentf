// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/opentofu/opentofu/internal/tfdiags"
	"github.com/posener/complete"
)

type WorkspaceListCommand struct {
	Meta
	LegacyName bool
}

func (c *WorkspaceListCommand) Run(args []string) int {
	args = c.Meta.process(args)

	var diags tfdiags.Diagnostics

	if c.LegacyName {
		envDiags := envCommandInvoked()
		diags = diags.Append(envDiags)

		c.showDiagnostics(diags)
		if c.View.HasErrors(envDiags) {
			return 1
		}
	}

	cmdFlags := c.Meta.defaultFlagSet("workspace list")
	c.Meta.varFlagSet(cmdFlags)
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s\n", err.Error()))
		return 1
	}

	args = cmdFlags.Args()
	configPath, err := modulePath(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// Load the encryption configuration
	enc, encDiags := c.EncryptionFromPath(configPath)
	diags = diags.Append(encDiags)
	if c.View.HasErrors(encDiags) {
		c.showDiagnostics(encDiags)
		return 1
	}

	backendConfig, backendDiags := c.loadBackendConfig(configPath)
	diags = diags.Append(backendDiags)
	if c.View.HasErrors(diags) {
		c.showDiagnostics(diags)
		return 1
	}

	// Load the backend
	b, backendDiags := c.Backend(&BackendOpts{
		Config: backendConfig,
	}, enc.State())
	diags = diags.Append(backendDiags)
	if c.View.HasErrors(diags) {
		c.showDiagnostics(diags)
		return 1
	}

	// This command will not write state
	c.ignoreRemoteVersionConflict(b)

	states, err := b.Workspaces()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	env, isOverridden := c.WorkspaceOverridden()

	var out bytes.Buffer
	for _, s := range states {
		if s == env {
			out.WriteString("* ")
		} else {
			out.WriteString("  ")
		}
		out.WriteString(s + "\n")
	}

	c.Ui.Output(out.String())

	if isOverridden {
		c.Ui.Output(envIsOverriddenNote)
	}

	return 0
}

func (c *WorkspaceListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictDirs("")
}

func (c *WorkspaceListCommand) AutocompleteFlags() complete.Flags {
	return nil
}

func (c *WorkspaceListCommand) Help() string {
	helpText := `
Usage: tofu [global options] workspace list [options]

  List OpenTofu workspaces.

Options:

  -var 'foo=bar'     Set a value for one of the input variables in the root
                     module of the configuration. Use this option more than
                     once to set more than one variable.

  -var-file=filename Load variable values from the given file, in addition
                     to the default files terraform.tfvars and *.auto.tfvars.
                     Use this option more than once to include more than one
                     variables file.
`
	return strings.TrimSpace(helpText)
}

func (c *WorkspaceListCommand) Synopsis() string {
	return "List Workspaces"
}
