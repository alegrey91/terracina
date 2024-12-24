required_providers {
  test = {
    source = "terracina.io/builtin/test"
  }
}

provider "test" "main" {
}

component "main" {
  source = "./"

  providers = {
    test = provider.test.main
  }
}
