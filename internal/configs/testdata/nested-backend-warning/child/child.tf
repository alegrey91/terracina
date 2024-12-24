terracina {
  # Only the root module can declare a backend. Terracina should emit a warning
  # about this child module backend declaration.
  backend "ignored" {
  }
}
