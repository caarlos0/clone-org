# clone-org


[![Release](https://img.shields.io/github/release/caarlos0/clone-org.svg?style=flat-square)](https://github.com/caarlos0/clone-org/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/caarlos0/clone-org.svg?style=flat-square)](https://travis-ci.org/caarlos0/clone-org)
<!-- [![Coverage Status](https://img.shields.io/coveralls/caarlos0/clone-org/master.svg?style=flat-square)](https://coveralls.io/github/caarlos0/clone-org?branch=master) -->
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/clone-org?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/clone-org)
[![Godoc](https://godoc.org/github.com/caarlos0/clone-org?status.svg&style=flat-square)](http://godoc.org/github.com/caarlos0/clone-org)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/caarlos0)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)


A simple command line tool to clone all repos of a given organization.

I needed to do that so I can `grep` all repos for some stuff. GitHub search
wasn't powerful enough to do what I needed, so, here it is.

## Install

```sh
brew install caarlos0/tap/clone-org
```

## Usage

```
NAME:
   clone-org - Clone all repos of a github organization

USAGE:
   clone-org [global options] command [command options] [arguments...]

VERSION:
   master

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
