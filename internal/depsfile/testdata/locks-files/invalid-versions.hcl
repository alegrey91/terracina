provider "terracina.io/test/foo" {
  version = "" # ERROR: Invalid provider version number
}

provider "terracina.io/test/bar" {
  # The "v" prefix is not expected here
  version = "v1.0.0" # ERROR: Invalid provider version number
}

provider "terracina.io/test/baz" {
  # Must be written in the canonical form, with three parts
  version = "1.0" # ERROR: Invalid provider version number
}

provider "terracina.io/test/boop" {
  # Must be written in the canonical form, with three parts
  version = "1" # ERROR: Invalid provider version number
}

provider "terracina.io/test/blep" {
  # Mustn't use redundant extra zero padding
  version = "1.02" # ERROR: Invalid provider version number
}

provider "terracina.io/test/huzzah" { # ERROR: Missing required argument
}

provider "terracina.io/test/huzznot" {
  version = null # ERROR: Missing required argument
}
