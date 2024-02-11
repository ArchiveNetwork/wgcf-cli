# scripts
# [Review by yourself before use.](../LICENSE)

## install.sh
- wgcf-cli install script

**Usage:**
- Install script supports `linux` / `macos` / `android(termux)`
- It will install `wgcf-cli` to `$PREFIX/bin/`
- Termux doesn't need manually set `$PREFIX`, it will use the default `$PREFIX`
```bash
PREFIX="~/.local" bash -c "$(curl -L wgcf-cli.vercel.app)"
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
## endpoint.sh
- Auto endpoint choose script.
- It **only used** to test the connectivity of the endpoint.
- It **won't** test the speed of the endpoint.
- It **won't** test the stability of the endpoint.
- It **is not** a benchmark script.

In the mainland of China, the result is useable.  
You can use it to choose a endpoint to connect.  
**Only** use it if you have a bad network connection (like in the mainland of China).  
If you have a good network connection, you should use the official endpoint.  
This script's result is based on the ICMP ping, so it's not very accurate.  
***(since the ICMP package is not the same as the UDP package)***

**Usage:**
```bash
sudo bash -c "$(curl -L wgcf-cli.vercel.app/endpoint)" [Arguments]
```
```ini
Arguments:  
[-6/--v6/--ipv6]: Only test IPv6  
[-4/--v4/--ipv4]: Only test IPv4  
[-y/--yes]: By pass the menu  
[-t/--time]: How many IPs to test  
[-s/--store]: Store the test result to ./result.txt
```
**Why use `sudo`?**
- Because the bring up a wireguard interface need root privileges.