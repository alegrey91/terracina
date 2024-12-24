terracina {
  required_providers {
    test = {
      source = "terracina.io/builtin/test"
    }
  }
}

provider "test" {
  arg = "foo"
}

module "b" {
  source = "../module-b"
}
