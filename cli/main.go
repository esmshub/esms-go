/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/esmshub/esms-go/cli/cmd"
	"go.uber.org/zap"
)

func main() {
	cmd.Execute()
	defer zap.L().Sync()
}
