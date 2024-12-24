terracina {
  backend "foo" {
    path = "relative/path/to/terracina.tfstate"
  }
}

resource "aws_instance" "web" {
  ami = "ami-1234"
  security_groups = [
    "foo",
    "bar",
  ]
}
