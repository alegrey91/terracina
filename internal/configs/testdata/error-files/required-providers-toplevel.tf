# A top-level required_providers block is not valid, but we have a specialized
# error for it to hint the user to move it into a terracina block.
required_providers { # ERROR: Invalid required_providers block
}

# This one is valid, and what the user should rewrite the above to be like.
terracina {
  required_providers {}
}
