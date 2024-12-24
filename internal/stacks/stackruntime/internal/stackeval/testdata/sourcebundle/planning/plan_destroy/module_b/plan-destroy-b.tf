
terracina {
  required_providers {
    test = {
      source = "terracina.io/builtin/test"

      configuration_aliases = [ test ]
    }
  }
}

variable "from_a" {
  type = string
}

resource "test" "foo" {
  for_module = "b"

  arg = var.from_a
}

output "result" {
  value = test.foo.result
}
