package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "run [-- args...]",
		Short: "Run the Java project",
		Long: `Run the Java project using java.

This command invokes the main module's Cli class using the module path.
Any additional arguments after -- are passed to the Java program.

The project must be compiled first with 'sai compile'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(args, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print exact command being executed")

	return cmd
}

func runRun(programArgs []string, verbose bool) error {
	projectID, err := detectProjectID()
	if err != nil {
		return err
	}

	javaArgs := []string{
		"-p", "lib:out",
		"-m", fmt.Sprintf("%s.main/%s.main.Cli", projectID, projectID),
	}
	javaArgs = append(javaArgs, programArgs...)

	if verbose {
		fmt.Printf("+ java %s\n", formatArgs(javaArgs))
	}

	javaCmd := exec.Command("java", javaArgs...)
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr
	javaCmd.Stdin = os.Stdin

	return javaCmd.Run()
}
