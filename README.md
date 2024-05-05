# wgcf-cli
A command-line tool for Cloudflare-WARP API, built using Cobra.<br>
Thanks: [@fscarmen](https://github.com/fscarmen/), [fscarmen/warp](https://gitlab.com/fscarmen/warp/), [@badafans](https://github.com/badafans), [badafans/warp-reg](https://github.com/badafans/warp-reg)<br>
## Example 
```
❯ wgcf-cli 
A command-line tool for Cloudflare-WARP API, built using Cobra.

Usage:
  wgcf-cli [command]

Available Commands:
  bind        check current bind devices
  cancel      cancel a config from original one
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  license     change to a new license
  name        change the device name
  register    Register a new WARP account
  simplify    Simplify a config from original one
  unbind      unbind from original license
  update      update a config
  version     Print the version of wgcf-cli

Flags:
  -c, --config string   set configuration file path (default "wgcf.json")
  -h, --help            help for wgcf-cli

Use "wgcf-cli [command] --help" for more information about a command.
```
## Build 
```bash
make
```
- Available environment variable:
1. `GOFLAGS`, Default: 
```bash
-trimpath -ldflags "-X github.com/ArchiveNetwork/wgcf-cli/constant.Version=$(VERSION) -s -w -buildid=" -v
```
2. `PREFIX`, Default: 
```bash
$(shell go env GOPATH)
```
3. `CGO_ENABLED`, Default: `0`
