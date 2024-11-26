package main

import (
	"github.com/spf13/cobra"

	"github.com/cloudweav/cloudweav/cmd/upgradehelper/cmd"
	_ "github.com/cloudweav/cloudweav/cmd/upgradehelper/cmd/versionguard"
	_ "github.com/cloudweav/cloudweav/cmd/upgradehelper/cmd/vmlivemigratedetector"
)

func main() {
	cobra.CheckErr(cmd.RootCmd.Execute())
}
