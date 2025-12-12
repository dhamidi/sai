package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/dhamidi/sai/ebnflex"
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
	cmd.AddCommand(newEbnfLexCmd())

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

func newEbnfLexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "lex <grammar>",
		Short:         "Tokenize input based on an EBNF grammar",
		Long:          "Reads input from stdin and emits tokens based on the grammar, including source positions.",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			grammarFile := args[0]

			grammar, err := ebnflex.LoadGrammar(grammarFile)
			if err != nil {
				return err
			}

			input, err := io.ReadAll(bufio.NewReader(os.Stdin))
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}

			lexer := ebnflex.NewLexer(grammar, input, "<stdin>")
			for {
				tok, err := lexer.NextToken()
				if err == io.EOF {
					fmt.Println(tok)
					break
				}
				if err != nil {
					return err
				}
				fmt.Println(tok)
			}

			return nil
		},
	}

	return cmd
}
