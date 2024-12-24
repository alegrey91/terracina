terracina {
  required_providers {
    test = {
      source = "terracina.io/builtin/test"
    }
  }
}

variable "marker" {
  type = string
}

variable "deps" {
  type    = set(string)
  default = []
}

resource "test_report" "main" {
  marker = var.marker
}

output "marker" {
  value = var.marker
}
