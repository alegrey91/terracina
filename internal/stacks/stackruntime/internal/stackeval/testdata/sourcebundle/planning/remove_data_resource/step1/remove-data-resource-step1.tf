
terracina {
  required_providers {
    test = {
        source = "terracina.io/builtin/test"
    }
  }
}

data "test" "test" {
}
