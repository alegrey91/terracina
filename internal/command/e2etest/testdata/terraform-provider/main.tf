provider "terracina" {

}

data "terracina_remote_state" "test" {
  backend = "local"
  config = {
    path = "test.tfstate"
  }
}
