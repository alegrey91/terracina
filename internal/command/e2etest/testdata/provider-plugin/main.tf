// the provider-plugin tests uses the -plugin-cache flag so terracina pulls the
// test binaries instead of reaching out to the registry.
terracina {
  required_providers {
    simple5 = {
      source = "registry.terracina.io/hashicorp/simple"
    }
    simple6 = {
      source = "registry.terracina.io/hashicorp/simple6"
    }
  }
}

resource "simple_resource" "test-proto5" {
  provider = simple5
}

resource "simple_resource" "test-proto6" {
  provider = simple6
}
