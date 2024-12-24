data "test_file" "foo" {
  id = "bar"
}

terracina {
  provider_meta "test" {
    baz = "quux-submodule"
  }
}
