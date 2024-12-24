terracina {
  required_providers {
    test = {
      source = "registry.terracina.io/hashicorp/test"
	}
  }
}

resource "test_object" "test" {
  test_string = provider::test::echo("input")
}
