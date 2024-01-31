#!/bin/sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

apk add curl

if [ "$1" = "--convenience" ]; then
  sh -x alpine-convenience.sh
else
  sh -x alpine-manual.sh
fi

tofu --version