terracina {
  required_providers {
    usererror = {
      source = "foo/terracina-provider-foo" # ERROR: Invalid provider type
    }
    badname = {
      source = "foo/terracina-foo" # ERROR: Invalid provider type
    }
  }
}
