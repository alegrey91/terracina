required_providers {
  terracina = {
    source = "terracina.io/builtin/terracina"
  }
}

provider "terracina" "default" {
}

component "self" {
  source = "./"

  providers = {
    terracina = provider.terracina.default
  }
}

output "obj" {
  type = object({
    input  = string
    output = string
  })
  value = component.self
}
