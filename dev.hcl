name = "test"

assert "check ls" {
  kind = "exists"
  file_path = "/bin/ls"
}

assert "check echo" {
  kind = "exists"
  file_path = "/bin/echo"
}

assert "check subshard" {
  kind = "exists"
  file_path = "/Applications/Subshard.app/Contents/MacOS/bin/subshard"
}
