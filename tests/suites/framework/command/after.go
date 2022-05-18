package command

import (
	"context"
	"log"

	"github.com/NicholeGit/nade/framework"
	"github.com/spf13/cobra"
)

var (
	gRootCmd   *cobra.Command
	gCtx       context.Context
	gContainer framework.IContainer
)

func After() {
	log.Println("After")
}

func TearDown() {
	log.Println("TearDown")
}
