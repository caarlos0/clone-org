# clone-org [![Build Status](https://travis-ci.org/caarlos0/clone-org.svg?branch=master)](https://travis-ci.org/caarlos0/clone-org) [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

A simple command line tool to clone all repos of a given organization.

I needed to do that so I can `grep` all repos for some stuff. GitHub search
wasn't powerful enough to what I needed, so, here it is.

## Usage

```
NAME:
   clone-org - Clone all repos of a github organization

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --org value, -o value
   --token value, -t value         [$GITHUB_TOKEN]
   --destination value, -d value
   --help, -h                     show help
   --version, -v                  print the version
```

## Notes

* if no destination is provided, the clone will be made in
`/tmp/organization-name`
* a `git clone --depth 1` will be performed, meaning that only the last commit
of the default branch will be available. On future versions this may be
configurable.
