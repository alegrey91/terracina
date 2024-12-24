variable "foo" { default = "bar" }

terracina {
    backend "local" {
        path = "${var.foo}"
    }
}
