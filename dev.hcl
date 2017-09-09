name = "test"

assert "check ls" {
  kind = "exists"
  file_path = "/bin/ls"
}

assert "check echo" {
  kind = "exists"
  file_path = "/bin/echo"
}

assert "dev assertionspec latest" {
  kind = "file_match"
  file_path = "~/dev.hcl"
  base_path = "dev.hcl"
  or "apply files" {
    action = "COPY"
    source_path = "dev.hcl"
    destination_path = "~/dev.hcl"
  }
}
