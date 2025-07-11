package main

import (
	"github.com/anyproto/anytype-cli/cmd"
	"github.com/anyproto/anytype-cli/core"
)

func main() {
	defer func() {
		core.CloseEventReceiver()
		core.CloseGRPCConnection()
	}()
	cmd.Execute()
}
