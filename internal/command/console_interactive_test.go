// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/opentofu/opentofu/internal/configs/configschema"
	"github.com/opentofu/opentofu/internal/providers"
	"github.com/zclconf/go-cty/cty"
)

func TestConsole_multiline_interactive(t *testing.T) {
	td := t.TempDir()
	testCopyDir(t, testFixturePath("console-multiline-vars"), td)
	defer testChdir(t, td)()

	p := testProvider()
	p.GetProviderSchemaResponse = &providers.GetProviderSchemaResponse{
		ResourceTypes: map[string]providers.Schema{
			"test_instance": {
				Block: &configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"value": {Type: cty.String, Optional: true},
					},
				},
			},
		},
	}
	ui := cli.NewMockUi()
	view, _ := testView(t)
	c := &ConsoleCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(p),
			Ui:               ui,
			View:             view,
		},
	}

	type testCase struct {
		input    string
		expected string
	}

	tests := map[string]testCase{
		"single_line": {
			input:    `var.counts.lalala`,
			expected: "1\n",
		},
		"basic_multi_line": {
			input: `
			var.counts.lalala
			var.counts.lololo`,
			expected: "1\n2\n",
		},
		"backets_multi_line": {
			input: `
			var.counts.lalala
			split(
			"_",
			"lalala_lolol_lelelele"
			)`,
			expected: "1\ntolist([\n  \"lalala\",\n  \"lolol\",\n  \"lelelele\",\n])\n",
		},
		"baces_multi_line": {
			input: `
			{ 
			for key, value in var.counts : key => value 
			if value == 1
			}`,
			expected: "{\n  \"lalala\" = 1\n}\n",
		},
	}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			var output bytes.Buffer
			defer testStdinPipe(t, strings.NewReader(tc.input))()
			outCloser := testStdoutCapture(t, &output)

			args := []string{}
			code := c.Run(args)
			outCloser()

			if code != 0 {
				t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
			}

			got := output.String()
			if got != tc.expected {
				t.Fatalf("unexpected output\ngot: %q\nexpected: %q", got, tc.expected)
			}
		})
	}
}
