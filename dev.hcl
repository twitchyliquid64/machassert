name = "test"

assert "check ls" {
  kind = "exists"
  order = 3
  file_path = "/bin/ls"
}

assert "check echo" {
  kind = "exists"
  order = 1
  file_path = "/bin/echo"
}

assert "dev assertionspec latest" {
  kind = "file_match"
  file_path = "~/dev.hcl"
  base_path = "dev.hcl"
  order = 2
  or "apply files" {
    action = "COPY"
    source_path = "dev.hcl"
    destination_path = "~/dev.hcl"
  }
}

assert "check other" {
  kind = "exists"
  order = 1
  file_path = "/bin/yelp"
  or "double check" {
    action = "ASSERT"

    assert "check echo 2" {
      kind = "exists"
      file_path = "/bin/echso"

      or "double check" {
        action = "ASSERT"

        assert "check echo 3" {
          kind = "exists"
          file_path = "/bin/echo"
        }
      }
    }
  }
}
