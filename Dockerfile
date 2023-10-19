# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

FROM alpine:3.18

LABEL maintainer="OpenTofu Team <opentofu@opentofu.org>"

RUN apk add --no-cache git

COPY tofu /usr/local/bin/tofu

ENTRYPOINT ["/usr/local/bin/tofu"]
