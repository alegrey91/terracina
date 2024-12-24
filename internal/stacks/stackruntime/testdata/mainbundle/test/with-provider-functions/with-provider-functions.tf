terracina {
  required_providers {
    testing = {
      source  = "hashicorp/testing"
      version = "0.1.0"
    }
  }
}

variable "id" {
  type     = string
  default  = null
  nullable = true # We'll generate an ID if none provided.
}

variable "input" {
  type = string
}

resource "testing_resource" "data" {
  id    = var.id
  value = provider::testing::echo(var.input)
}

output "value" {
  value = testing_resource.data.value
}
