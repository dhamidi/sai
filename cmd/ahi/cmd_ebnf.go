package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"golang.org/x/exp/ebnf"
)

func newEbnfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ebnf",
		Short:         "EBNF grammar tools",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newEbnfCheckCmd())

	return cmd
}

func newEbnfCheckCmd() *cobra.Command {
	var startProduction string

	cmd := &cobra.Command{
		Use:           "check <file>",
		Short:         "Parse and verify an EBNF grammar file",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]

			f, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer f.Close()

			grammar, err := ebnf.Parse(filename, f)
			if err != nil {
				printErrors(err)
				return err
			}

			if err := ebnf.Verify(grammar, startProduction); err != nil {
				printErrors(err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&startProduction, "start", "", "start production for verification (if empty, only checks syntax)")

	return cmd
}

func printErrors(err error) {
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			fmt.Println(v.Index(i).Interface())
		}
	} else {
		fmt.Println(err)
	}
}
