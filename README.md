# machassert

Assertions for system administrators. Puppet but easy to understand, and without the nonsense.

You write assertions to check files-exist/packages-installed/process-running/etc. You write actions which occur when an assertion fails: Copy this file, install this package etc. By combining these two features, you get a scriptable, easy to understand means to deploy software across a fleet.

## How to use

You have two files: `target` files (tell machassert which machines you want to run assertions on and how to connect to them), and `assertion` files (assertions and actions).

You run machassert like this: `./massert --target <target-file> assert <assertion-file>`

### Assertion files

Assertion files have a name, then a list of `assert` sections. Each `assert` block is an assertion containing information about the kind of assertion, any information the assertion needs,
and (optionally) actions to take if the assertion fails.

```hcl
name = "frontend"

assert "fisher installed" {
  kind = "exists"
  file_path = "/bin/fisher"
}

assert "fisher default config" {
  kind = "exists"
  file_path = "~/.fisher/defaults.hcl"
  or {
    action = "COPY"
    source_path = "defaults"
    destination_path = "~/.fisher/defaults.hcl"
  }
}
```

The above example does the following:

1. Checks `/bin/fisher` exists on the target system. If it does not, the assertion fails and the script terminates.
2. Checks the `~/.fisher/defaults.hcl` file exists on the target system. If it does not, the `COPY` action runs, copying the file.

### Target files

If no target file is specified, the assertions are run on the local system.

```hcl
name = "Frontend servers"

machine "frontend-1" {
  kind = "ssh"
  destination = "10.5.32.1"
  auth {
      password = "amazingly_secure"
  }
}
machine "frontend-2" {
  kind = "ssh"
  destination = "10.5.32.2"
  auth {
      password = "amazingly_secure"
  }
}
```
