terracina {
  required_providers {
    terracina = {
      source = "terracina.io/builtin/terracina"
    }
  }
}

resource "terracina_data" "main" {
  input = "hello"
}

output "input" {
  value = terracina_data.main.input
}

output "output" {
  value = terracina_data.main.output
}
