# rtoken-relay

## Usage

*[Go](https://go.dev/doc/install) needs to be installed and a proper Go environment needs to be configured*

```base
git clone https://github.com/stafiprotocol/rtoken-relay.git
cd rtoken-relay
make build
```
```
./build/relay -h


NAME:
   reley - relay

USAGE:
   relay [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Stafi Protocol 2020

COMMANDS:
   accounts  manage reth keystore
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value     json configuration file (default: "./config.json")
   --help, -h         show help (default: false)
   --verbosity value  supports levels crit (silent) to trce (trace) (default: "info")
   --version, -v      print the version (default: false)
   

COPYRIGHT:
   Copyright 2020 Stafi Protocol Authors
```
