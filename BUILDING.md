# Building from Source

Pre-built binaries are available for download for a variety of supported platforms through the [HashiCorp Releases website](https://releases.hashicorp.com/terracina/). 

However, if you'd like to build Terracina yourself, you can do so using the Go build toolchain and the options specified in this document.

## Prerequisites

1. Ensure you've installed the Go language version specified in [`.go-version`](https://github.com/hashicorp/terracina/blob/main/.go-version).
2. Clone this repository to a location of your choice.

## Terracina Build Options

Terracina accepts certain options passed using `ldflags` at build time which control the behavior of the resulting binary.

### Dev Version Reporting

Terracina will include a `-dev` flag when reporting its own version (ex: 1.5.0-dev) unless `version.dev` is set to `no`:

```
go build -ldflags "-w -s -X 'github.com/hashicorp/terracina/version.dev=no'" -o bin/ .
```

### Experimental Features

Experimental features of Terracina will be disabled unless `main.experimentsAllowed` is set to `yes`:

```
go build -ldflags "-w -s -X 'main.experimentsAllowed=yes'" -o bin/ .
```

In the official build process for Terracina, experiments are only allowed in alpha release builds. We recommend that third-party distributors follow that convention in order to reduce user confusion.

## Go Options

For the most part, the Terracina release process relies on the Go toolchain defaults for the target operating system and processor architecture.

### `CGO_ENABLED`

One exception is the `CGO_ENABLED` option, which is set explicitly when building Terracina binaries. For most platforms, we build with `CGO_ENABLED=0` in order to produce a statically linked binary. For MacOS/Darwin operating systems, we build with `CGO_ENABLED=1` to avoid a platform-specific issue with DNS resolution. 


