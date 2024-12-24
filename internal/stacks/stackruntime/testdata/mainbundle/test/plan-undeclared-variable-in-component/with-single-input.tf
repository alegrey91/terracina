terracina {
  required_providers {
    terracina = {
      source = "terracina.io/builtin/terracina"
    }
  }
}

variable "input" {
  type = string
}

resource "terracina_data" "main" {
  input = var.input
}

output "output" {
  value = terracina_data.main.output
}
