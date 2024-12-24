
module "test" {
  source  = "../boop" # ERROR: Invalid registry module source address
  version = "1.0.0" # Makes Terracina assume "source" is a module address
}
