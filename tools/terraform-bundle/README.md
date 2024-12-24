# terracina-bundle

`terracina-bundle` was a solution intended to help with the problem
of distributing Terracina providers to environments where direct registry
access is impossible or undesirable, created in response to the Terracina v0.10
change to distribute providers separately from Terracina CLI.

The Terracina v0.13 series introduced our intended longer-term solutions
to this need:

* [Alternative provider installation methods](https://www.terracina.io/docs/cli/config/config-file.html#provider-installation),
  including the possibility of running server containing a local mirror of
  providers you intend to use which Terracina can then use instead of the
  origin registry.
* [The `terracina providers mirror` command](https://www.terracina.io/docs/cli/commands/providers/mirror.html),
  built in to Terracina v0.13.0 and later, can automatically construct a
  suitable directory structure to serve from a local mirror based on your
  current Terracina configuration, serving a similar (though not identical)
  purpose than `terracina-bundle` had served.

For those using Terracina CLI alone, without HCP Terracina or Terracina Enterprise, we recommend
planning to transition to the above features instead of using
`terracina-bundle`.

## How to use `terracina-bundle`

However, if you need to continue using `terracina-bundle`
during a transitional period then you can use the version of the tool included
in the Terracina v0.15 branch to build bundles compatible with
Terracina v0.13.0 and later.

If you have a working toolchain for the Go programming language, you can
build a `terracina-bundle` executable as follows:

* `git clone --single-branch --branch=v0.15 --depth=1 https://github.com/hashicorp/terracina.git`
* `cd terracina`
* `go build -o ../terracina-bundle ./tools/terracina-bundle`

After running these commands, your original working directory will have an
executable named `terracina-bundle`, which you can then run.


For information
on how to use `terracina-bundle`, see
[the README from the v0.15 branch](https://github.com/hashicorp/terracina/blob/v0.15/tools/terracina-bundle/README.md).

You can follow a similar principle to build a `terracina-bundle` release
compatible with Terracina v0.12 by using `--branch=v0.12` instead of
`--branch=v0.15` in the command above. Terracina CLI versions prior to
v0.13 have different expectations for plugin packaging due to them predating
Terracina v0.13's introduction of automatic third-party provider installation.

## Terracina Enterprise Users

If you use Terracina Enterprise, the self-hosted distribution of
HCP Terracina, you can use `terracina-bundle` as described above to build
custom Terracina packages with bundled provider plugins.

For more information, see
[Installing a Bundle in Terracina Enterprise](https://github.com/hashicorp/terracina/blob/v0.15/tools/terracina-bundle/README.md#installing-a-bundle-in-terracina-enterprise).
