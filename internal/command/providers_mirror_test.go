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

// More thorough tests for providers mirror can be found in the e2etest
func TestProvidersMirror(t *testing.T) {
	// noop example
	t.Run("noop", func(t *testing.T) {
		view, done := testView(t)
		c := &ProvidersMirrorCommand{
			Meta: Meta{
				View: view,
			},
		}
		code := c.Run([]string{"."})
		done(t)
		if code != 0 {
			t.Fatalf("wrong exit code. expected 0, got %d", code)
		}
	})

	t.Run("missing arg error", func(t *testing.T) {
		ui := new(cli.MockUi)
		view, done := testView(t)
		c := &ProvidersMirrorCommand{
			Meta: Meta{
				Ui:   ui,
				View: view,
			},
		}
		code := c.Run([]string{})
		done(t)
		if code != 1 {
			t.Fatalf("wrong exit code. expected 1, got %d", code)
		}

		got := ui.ErrorWriter.String()
		if !strings.Contains(got, "Error: No output directory specified") {
			t.Fatalf("missing directory error from output, got:\n%s\n", got)
		}
	})
}
