name = "Frontend server"

machine "frontend-1" {
  kind = "ssh"
  destination = "10.5.32.1"
  auth "use key" {
      // Try a private key at /etc/secret.pem
      kind = "key-file"
      key = "/etc/secret.pem"
  }
}
