terracina {
  required_providers {
    template  = { version = "2.1.1" }
    null      = { source = "hashicorp/null", version = "2.1.0" }
    terracina = { source = "terracina.io/builtin/terracina" }
  }
}
