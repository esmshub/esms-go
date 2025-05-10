/*
Copyright © 2025 James Howe <james@esmshub.com>
*/
package main

import (
	"github.com/esmshub/esms-go/cmd/cli/cmd"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	cmd.Execute(version)
	defer zap.L().Sync()
}
