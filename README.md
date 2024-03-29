# wgcf-cli
Use the standard library to access Cloudflare-WARP api.<br>
Thanks: [@fscarmen](https://github.com/fscarmen/), [fscarmen/warp](https://gitlab.com/fscarmen/warp/), [@badafans](https://github.com/badafans), [badafans/warp-reg](https://github.com/badafans/warp-reg)<br>
## Install
- Install script supports `linux` / `macos` / `android(termux)`
- It will install `wgcf-cli` to `$PREFIX/bin/`
- Termux doesn't need manually set `$PREFIX`, it will use the default `$PREFIX`
```bash
PREFIX=~/.local bash -c "$(curl -L wgcf-cli.vercel.app)"
```
- when without `$PREFIX`, you may need root privileges to run the script
- It will install `wgcf-cli` to `/usr/local/bin/`
```bash
sudo bash -c "$(curl -L wgcf-cli.vercel.app)"
```
- Also, `BETA` Environment is available
```bash
BETA=true bash -c "$(curl -L wgcf-cli.vercel.app)"
```
## Example 
1. Register
```console
❯ wgcf-cli -r
{
    "endpoint": {
        "v4": "162.159.192.7:0",
        "v6": "[2606:4700:d0::a29f:c007]:0"
    },
    "reserved_str": "6nT5",
    "reserved_hex": "0xea74f9",
    "reserved_dec": [
        234,
        116,
        249
    ],
    "private_key": "WIAKvgUlq5fBazhttCvjhEGpu8MmGHcb1H0iHSGlU0Q=",
    "public_key": "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
    "addresses": {
        "v4": "172.16.0.2",
        "v6": "2606:4700:110:8d9c:3c4e:2190:59d1:2d3c"
    }
}
```
2. Bind a License
```console
❯ wgcf-cli -l 9zs5I61a-l9j8m7T5-4pC6k20X
{
    "id": "cd7f4695-e9ef-4bb0-b412-5f4d84919db7",
    "created": "0001-01-01T00:00:00Z",
    "updated": "2023-12-14T12:32:18.689777921Z",
    "premium_data": 0,
    "quota": 0,
    "warp_plus": true,
    "referral_count": 0,
    "referral_renewal_countdown": 0,
    "role": "child"
}
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
