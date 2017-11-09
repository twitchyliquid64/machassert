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
  or "copy defaults" {
    action = "COPY"
    source_path = "defaults"
    destination_path = "~/.fisher/defaults.hcl"
  }

  or "assert other condition on failure thingy" {
    action = "ASSERT"

    assert "check thingy" {
      kind = "exists"
      file_path = "/bin/sillyness"
    }
  }
}
```

The above example does the following:

1. Checks `/bin/fisher` exists on the target system. If it does not, the assertion fails and the script terminates.
2. Checks the `~/.fisher/defaults.hcl` file exists on the target system. If it does not, the `COPY` action runs, copying the file, as well as running the assertion in the other `OR` block (which will fail if `/bin/sillyness` does not exist).

#### Available assertions

| Kind          | Description           | Parameters  |
| ------------- |:----------------------| ------------|
| exists | Fails if the path at `file_path` does not exist. | `file_path` |
| !exists | Fails if the path at `file_path` does exists. | `file_path` |
| md5_match | Fails if the file at `file_path` does not have an MD5 hash that matches `hash`. | `file_path`, `hash` |
| file_match | Fails if the file at `file_path` does not match the file at `base_path`. Base path should be present on the machine from which machassert is being executed. | `file_path`, `base_path` |
| regex_contents_match | Fails if `regex` does not match any line in `file_path`. | `regex`, `file_path` |

#### Available actions

Actions are run if the assertion which contains it does not hold to be true.

| Action          | Description           | Additional fields required  |
| ------------- |:----------------------| ------------|
| FAIL | Default. Immediately fail and stop iterating through assertions. |  None. |
| COPY | Copy a file from the local machine to the machine being asserted on. | The `OR` block must contain parameters `source_path` & `destination_path` |
| ASSERT | Specify another set of assertions to run. | Additional named `assert` blocks must be present. |

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
      // prompt the user for credentials
      kind = "prompt"
  }
}
machine "frontend-3" {
  kind = "ssh"
  destination = "10.5.32.3"
  auth {
      // use the current users ssh private key
      kind = "user-key"
  }
}
machine "frontend-4" {
  kind = "ssh"
  destination = "10.5.32.4"
  auth {
      // use the PEM encoded private key at /etc/secret.pem
      kind = "key-file"
      key = "/etc/secret.pem"
  }
}
```
