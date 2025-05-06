/*
Copyright © 2025 James Howe <james@esmshub.com>
*/
package main

import (
	"github.com/esmshub/esms-go/cmd/cli/cmd"
	"go.uber.org/zap"
)

func main() {
	cmd.Execute()
	defer zap.L().Sync()
}
