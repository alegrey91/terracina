terracina {
  required_providers {
    tls = {
      source  = "hashicorp/tls"
      version = "~> 3.0"
    }
  }
}

# There is no provider in required_providers called "implied", so this
# implicitly declares a dependency on "hashicorp/implied".
resource "implied_foo" "bar" {
}

# There is no provider in required_providers called "terracina", but for
# this name in particular we imply terracina.io/builtin/terracina instead,
# to avoid selecting the now-unmaintained
# registry.terracina.io/hashicorp/terracina.
data "terracina_remote_state" "bar" {
}
