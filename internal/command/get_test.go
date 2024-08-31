// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGet(t *testing.T) {
	wd := tempWorkingDirFixture(t, "get")
	defer testChdir(t, wd.RootModuleDir())()

	ui := cli.NewMockUi()
	view, done := testView(t)
	c := &GetCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			WorkingDir:       wd,
			View:			  view,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
	done(t)

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "- foo in") {
		t.Fatalf("doesn't look like get: %s", output)
	}
}

func TestGet_multipleArgs(t *testing.T) {
	wd := tempWorkingDir(t)
	defer testChdir(t, wd.RootModuleDir())()

	ui := cli.NewMockUi()
	view, done := testView(t)
	c := &GetCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			WorkingDir:       wd,
			View:			  view,
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

func TestGet_update(t *testing.T) {
	wd := tempWorkingDirFixture(t, "get")
	defer testChdir(t, wd.RootModuleDir())()

	ui := cli.NewMockUi()
	view, done := testView(t)
	c := &GetCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			WorkingDir:       wd,
			View:			  view,
		},
	}

	args := []string{
		"-update",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
	done(t)

	output := ui.OutputWriter.String()
	if !strings.Contains(output, `- foo in`) {
		t.Fatalf("doesn't look like get: %s", output)
	}
}

func TestGet_cancel(t *testing.T) {
	// This test runs `tofu get` as if SIGINT (or similar on other
	// platforms) were sent to it, testing that it is interruptible.

	wd := tempWorkingDirFixture(t, "init-registry-module")
	defer testChdir(t, wd.RootModuleDir())()

	// Our shutdown channel is pre-closed so init will exit as soon as it
	// starts a cancelable portion of the process.
	shutdownCh := make(chan struct{})
	close(shutdownCh)

	ui := cli.NewMockUi()
	view, done := testView(t)
	c := &GetCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			WorkingDir:       wd,
			ShutdownCh:       shutdownCh,
			View:			  view,
		},
	}

	args := []string{}
	if code := c.Run(args); code == 0 {
		t.Fatalf("succeeded; wanted error\n%s", ui.OutputWriter.String())
	}
	done(t)

	if got, want := ui.ErrorWriter.String(), `Module installation was canceled by an interrupt signal`; !strings.Contains(got, want) {
		t.Fatalf("wrong error message\nshould contain: %s\ngot:\n%s", want, got)
	}
}
