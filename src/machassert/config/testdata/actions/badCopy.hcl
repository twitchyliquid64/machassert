name = "test"

assert "missing copy fields" {
  kind = "exists"
  file_path = "~/dev.hcl"
  or "apply files" {
    action = "COPY"
  }
}
