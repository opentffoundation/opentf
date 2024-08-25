// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
)

// ColorizeUi is a Ui implementation that colors its output according
// to the given color schemes for the given type of output.
type ColorizeUi struct {
	Colorize    *colorstring.Colorize
	OutputColor string
	InfoColor   string
	ErrorColor  string
	WarnColor   string
	Ui          cli.Ui
}

func (u *ColorizeUi) Ask(query string) (string, error) {
	return u.Ui.Ask(u.colorize(query, u.OutputColor))
}

func (u *ColorizeUi) AskSecret(query string) (string, error) {
	return u.Ui.AskSecret(u.colorize(query, u.OutputColor))
}

func (u *ColorizeUi) Output(message string) {
	u.Ui.Output(u.colorize(message, u.OutputColor))
}

func (u *ColorizeUi) Info(message string) {
	u.Ui.Info(u.colorize(message, u.InfoColor))
}

func (u *ColorizeUi) Error(message string) {
	u.Ui.Error(u.colorize(message, u.ErrorColor))
}

func (u *ColorizeUi) Warn(message string) {
	u.Ui.Warn(u.colorize(message, u.WarnColor))
}

func (u *ColorizeUi) colorize(message string, color string) string {
	if color == "" {
		return message
	}

	return u.Colorize.Color(fmt.Sprintf("%s%s[reset]", color, message))
}

// pedanticUI is a UI implementation which directs warning messages to the error output stream
type pedanticUI struct {
	cli.Ui
	notifyWarning func()
}

// Warn notifies of the warning
func (pui *pedanticUI) Warn(msg string) {
	pui.Ui.Error(msg)
	pui.notifyWarning()
}
