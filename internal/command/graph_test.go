// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/opentofu/internal/addrs"
	"github.com/opentofu/opentofu/internal/plans"
	"github.com/opentofu/opentofu/internal/states"
)

func TestGraph(t *testing.T) {
	td := t.TempDir()
	testCopyDir(t, testFixturePath("graph"), td)
	defer testChdir(t, td)()

	ui := new(cli.MockUi)
	view, done := testView(t)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(applyFixtureProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
	done(t)

	output := ui.OutputWriter.String()
	if !strings.Contains(output, `provider[\"registry.opentofu.org/hashicorp/test\"]`) {
		t.Fatalf("doesn't look like digraph: %s", output)
	}
}

func TestGraph_multipleArgs(t *testing.T) {
	ui := new(cli.MockUi)
	view, done := testView(t)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(applyFixtureProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{
		"bad",
		"bad",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.OutputWriter.String())
	}
	done(t)
}

func TestGraph_noArgs(t *testing.T) {
	td := t.TempDir()
	testCopyDir(t, testFixturePath("graph"), td)
	defer testChdir(t, td)()

	ui := new(cli.MockUi)
	view, done := testView(t)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(applyFixtureProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
	done(t)

	output := ui.OutputWriter.String()
	if !strings.Contains(output, `provider[\"registry.opentofu.org/hashicorp/test\"]`) {
		t.Fatalf("doesn't look like digraph: %s", output)
	}
}

func TestGraph_noConfig(t *testing.T) {
	td := t.TempDir()
	os.MkdirAll(td, 0755)
	defer testChdir(t, td)()

	ui := new(cli.MockUi)
	view, done := testView(t)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(applyFixtureProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	// Running the graph command without a config should not panic,
	// but this may be an error at some point in the future.
	args := []string{"-type", "apply"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
	done(t)
}

func TestGraph_plan(t *testing.T) {
	testCwd(t)

	plan := &plans.Plan{
		Changes: plans.NewChanges(),
	}
	plan.Changes.Resources = append(plan.Changes.Resources, &plans.ResourceInstanceChangeSrc{
		Addr: addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "test_instance",
			Name: "bar",
		}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
		ChangeSrc: plans.ChangeSrc{
			Action: plans.Delete,
			Before: plans.DynamicValue(`{}`),
			After:  plans.DynamicValue(`null`),
		},
		ProviderAddr: addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("test"),
			Module:   addrs.RootModule,
		},
	})
	beConfig := cty.ObjectVal(map[string]cty.Value{
		"path":          cty.NilVal,
		"workspace_dir": cty.NilVal,
	})
	emptyConfig, err := plans.NewDynamicValue(beConfig, beConfig.Type())
	if err != nil {
		t.Fatal(err)
	}
	plan.Backend = plans.Backend{
		Type:   "local",
		Config: emptyConfig,
	}
	_, configSnap := testModuleWithSnapshot(t, "graph")

	planPath := testPlanFile(t, configSnap, states.NewState(), plan)

	ui := new(cli.MockUi)
	view, done := testView(t)
	c := &GraphCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(applyFixtureProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{
		"-plan", planPath,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
	done(t)

	output := ui.OutputWriter.String()
	if !strings.Contains(output, `provider[\"registry.opentofu.org/hashicorp/test\"]`) {
		t.Fatalf("doesn't look like digraph: %s", output)
	}
}
