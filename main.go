package main

import "github.com/ArchiveNetwork/wgcf-cli/cmd"

func main() {
	if err := cmd.Judgement(); err != nil {
		panic(err)
	}
}
