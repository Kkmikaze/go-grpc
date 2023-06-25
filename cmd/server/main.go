package main

import (
	"github.com/Kkmikaze/roketin/cmd"
	"github.com/Kkmikaze/roketin/configs"
)

func init() {
	configs.InitConfigs()
}

func main() {
	cmd.Execute()
}
