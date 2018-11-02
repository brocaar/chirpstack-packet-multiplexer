package main

import "github.com/brocaar/lora-packet-multiplexer/cmd/lora-packet-multiplexer/cmd"

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
