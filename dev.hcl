name = "test"

assert "check ls" {
  kind = "exists"
  file_path = "/bin/ls"
}

assert "check echo" {
  kind = "exists"
  file_path = "/bin/echo"
}

assert "check hash" {
  kind = "md5_match"
  file_path = "/bin/ls"
  hash = "970be6a05c1ccbadbcece0c6db9b3882"
}
