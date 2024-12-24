# This is a simple configuration with HCP Terracina mode minimally
# activated, but it's suitable only for testing things that we can exercise
# without actually accessing HCP Terracina, such as checking of invalid
# command-line options to "terracina init".

terracina {
  cloud {
    organization = "PLACEHOLDER"
    workspaces {
        name = "PLACEHOLDER"
    }
  }
}
