name = "frontend"

assert "binary" {
  kind = "exists"
  file_path = "/bin/ls"
  or {
    action = "FAIL"
  }
}
