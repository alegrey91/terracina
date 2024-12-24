module "grandchild" {
  # NOTE: This only works because our caller told Terracina to treat
  # the parent directory as a whole as a module package, and so
  # the "./terracina/modules/child" directory should contain both
  # "child" and "grandchild" sub directories that we can traverse between.
  # This is the same as local paths between different directories inside
  # a single git repository or distribution archive.
  source = "../grandchild"
}
