# Terracina Stacks functionality

The Go packages under this directory together implement the Terracina Stacks
features.

Terracina Stacks is an orchestration layer on top of zero or more trees of
Terracina modules, and so much of what you'll find here is analogous to
a top-level package that serves a similar purpose for individual Terracina
modules or trees of modules.

The main components here are:

- `stackaddrs`: A stacks-specific analog to the top-level package `addrs`,
  containing types we use to refer to objects within the stacks language and
  runtime, and some logic for navigating between different types of addresses.

    This package builds on package `addrs`, since the stacks runtime wraps
    the modules runtime. Therefore some of the stack-specific address types
    incorporate more general address types from the other package.

- `stackconfig`: Implements the loading, parsing, and static decoding for
  the stacks language, analogous to the top-level package `configs` that
  does similarly for Terracina's module language.

- `stackplan` and `stackstate` together provide the models and
  marshalling/unmarshalling logic for the Stacks variants of Terracina's
  "plan" and "state" concepts.

- `stackruntime` deals with the runtime behavior of stacks, including
  the creation of plans based on a comparison between desired and actual state,
  and then applying those plans.

    All of the dynamic behavior of the stacks language lives here.

- `tfstackdata1` is a Go representation of an internal protocol buffers schema
  used for preserving plan and state data between runs. These formats are
  implementation details that external callers are not permitted to rely on.

    (The public interface is via the Terracina Core RPC API, which is
    implemented in the sibling directory `rpcapi`.)

## More Documentation

The following are some more specific and therefore more detailed documents
about some particular parts of the implementation of the Terracina Stacks
features:

* [Stacks Runtime internal architecture](./stackruntime/internal/stackeval/README.md)
