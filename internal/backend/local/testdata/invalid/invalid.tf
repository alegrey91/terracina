# This configuration is intended to be loadable (valid syntax, etc) but to
# fail terracina.Context.Validate.

locals {
  a = local.nonexist
}
