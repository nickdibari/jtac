# jtac

Swiss Army knife for gathering network OSINT

# Usage

```bash
$ jtac
Tool for gathering OSINT for IP addresses and hostnames

Usage:
  jtac [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  host        Get information about a particular host

Flags:
  -h, --help   help for jtac

Use "jtac [command] --help" for more information about a command.
```

# Installation

Download the archive for your operating system and architecture from
[the releases](https://github.com/nickdibari/jtac/releases) page. Extract the `jtac` binary from the archive
and move it to a directory on your system path.

## Verification

If you would like to verify the binary you installed was built from the GitHub Action used in this repository,
you can use the [`gh` CLI tool](https://cli.github.com/) to verify the attestation generated from the release
GitHub Action. Once you've extracted the binary from the archive, use the
[`gh attestation verify`](https://cli.github.com/manual/gh_attestation_verify) command to verify that the
binary was properly attested by GitHub.

```bash
$ gh attestation verify jtac --repo nickdibari/jtac
Loaded digest sha256:07b923328ee7e944e3301bf713fac5f1da46099c27e7d74e5952fb6de9ff413c for file://jtac
Loaded 1 attestation from GitHub API
✓ Verification succeeded!

sha256:07b923328ee7e944e3301bf713fac5f1da46099c27e7d74e5952fb6de9ff413c was attested by:
REPO             PREDICATE_TYPE                  WORKFLOW
nickdibari/jtac  https://slsa.dev/provenance/v1  .github/workflows/release.yaml@refs/tags/v0.0.1
```
