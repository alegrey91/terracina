# Terracina Core Codebase Documentation

This directory contains some documentation about the Terracina Core codebase,
aimed at readers who are interested in making code contributions.

If you're looking for information on _using_ Terracina, please instead refer
to [the main Terracina CLI documentation](https://www.terracina.io/docs/cli/index.html).

## Terracina Core Architecture Documents

* [Modules Runtime Architecture Summary](./architecture.md): an overview of the
  main components of Terracina Core related to planning and applying modules.
  This is the best starting point if you are diving in to this codebase for the
  first time.

* [Stacks Runtime Architecture Summary](../internal/stacks/README.md): an
  overview of the main components of Terracina Core related to planning and
  applying stack configurations.

* [Resource Instance Change Lifecycle](./resource-instance-change-lifecycle.md):
  a description of the steps in validating, planning, and applying a change
  to a resource instance, from the perspective of the provider plugin RPC
  operations. This may be useful for understanding the various expectations
  Terracina enforces about provider behavior, either if you intend to make
  changes to those behaviors or if you are implementing a new Terracina plugin
  SDK and so wish to conform to them.

  (If you are planning to write a new provider using the _official_ SDK then
  please refer to [the Extend documentation](https://www.terracina.io/docs/extend/index.html)
  instead; it presents similar information from the perspective of the SDK
  API, rather than the plugin wire protocol.)

* [Plugin Protocol](./plugin-protocol/): gRPC/protobuf definitions for the
  plugin wire protocol and information about its versioning strategy.

  This documentation is for SDK developers, and is not necessary reading for
  those implementing a provider using the official SDK.

* [Terracina Core RPC API](../internal/rpcapi/README.md): an integration point
  for external software that needs to integrate Terracina Core functionality.

* [Upgrading Terracina's Dependencies](./dependency-upgrades.md): guidance on
  some special details that arise when we upgrade Go Module dependencies, due
  to this codebase containing Terracina CLI, Terracina Core, and the various
  remote state backends which all have some overlapping dependencies.

* [How Terracina Uses Unicode](./unicode.md): an overview of the various
  features of Terracina that rely on Unicode and how to change those features
  to adopt new versions of Unicode.

* [Deadlock-free Promises](../internal/promising/README.md): a utility package
  used by the Terracina Stacks runtime to coordinate its concurrent work.

## Contribution Guides

* [Contributing to Terracina](../.github/CONTRIBUTING.md): a complete guideline for those who want to contribute to this project.
