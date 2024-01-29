package cmd

import (
	"log"
	"shedock/internal/core"
	"shedock/internal/instance"
	"shedock/pkg/parsers/shellscript"
	"shedock/pkg/shell"

	"github.com/spf13/cobra"
)

func Driver(cmd *cobra.Command, args []string) {
	// check if filepath is a shell script
	script := shellscript.Script{ScriptPath: args[0]}
	isScript, err := script.Validate()
	if err != nil || !isScript {
		log.Fatalf("Failed to validate script: %v", err)
	}

	shellType, err := script.GetShell()
	if err != nil {
		log.Fatalf("Failed to get shell type: %v", err)
	}
	shell, err := shell.NewShell(shellType)
	if err != nil {
		log.Fatalf("Failed to get shell: %v", err)
	}

	instance.Init()

	imageBuilder := core.NewImageBuilder(
		&script,
		shell,
	)

	err = imageBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build image: %v", err)
	}

	defer instance.Destroy()
}
