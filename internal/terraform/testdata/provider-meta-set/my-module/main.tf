resource "test_resource" "bar" {
  value = "bar"
}

terracina {
  provider_meta "test" {
    baz = "quux-submodule"
  }
}
