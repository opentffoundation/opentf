# OpenTF Core Codebase Documentation

This directory contains some documentation about the OpenTF Core codebase,
aimed at readers who are interested in making code contributions.

If you're looking for information on _using_ OpenTF, please instead refer
to [the main OpenTF CLI documentation](https://opentf.org/docs/cli/index.html).

## OpenTF Core Architecture Documents

* [OpenTF Core Architecture Summary](./architecture.md): an overview of the
  main components of OpenTF Core and how they interact. This is the best
  starting point if you are diving in to this codebase for the first time.

* [Resource Instance Change Lifecycle](./resource-instance-change-lifecycle.md):
  a description of the steps in validating, planning, and applying a change
  to a resource instance, from the perspective of the provider plugin RPC
  operations. This may be useful for understanding the various expectations
  OpenTF enforces about provider behavior, either if you intend to make
  changes to those behaviors or if you are implementing a new OpenTF plugin
  SDK and so wish to conform to them.

  (If you are planning to write a new provider using the _official_ SDK then
  please refer to [the Extend documentation](https://opentf.org/docs/extend/index.html)
  instead; it presents similar information from the perspective of the SDK
  API, rather than the plugin wire protocol.)

* [Plugin Protocol](./plugin-protocol/): gRPC/protobuf definitions for the
  plugin wire protocol and information about its versioning strategy.

  This documentation is for SDK developers, and is not necessary reading for
  those implementing a provider using the official SDK.

* [How OpenTF Uses Unicode](./unicode.md): an overview of the various
  features of OpenTF that rely on Unicode and how to change those features
  to adopt new versions of Unicode.

## Contribution Guides

* [Contributing to OpenTF](../.github/CONTRIBUTING.md): a complete guideline for those who want to contribute to this project.
