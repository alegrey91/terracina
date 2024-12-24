
terracina {
  backend "local" {
    path = $invalid
  }
}

variable "input" {
  type = string
}
