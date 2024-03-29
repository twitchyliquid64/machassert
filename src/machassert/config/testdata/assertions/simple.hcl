name = "frontend"

assert "binary" {
  kind = "exists"
  order = 1
  file_path = "/bin/ls"
}

assert "thing" {
  kind = "exists"
  file_path = "/bin/ls"
  or {
    source_path = "dev.hcl"
    destination_path = "~/dev.hcl"
  }
}
