name = "bad assert action"

assert "check missing" {
  kind = "exists"
  file_path = "sdfsdfsd"

  or "sasadas" {
    action = "ASSERT"
  }
}
