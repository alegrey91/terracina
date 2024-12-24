// At the time of writing Terracina doesn't formally support a boolean
// type, but historically this has magically worked. Lots of TF code
// relies on this so we test it now.
variable "a" {
    default = true
}

variable "b" {
    default = false
}
