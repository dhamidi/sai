package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/dhamidi/sai/ebnf/grammar"
	"github.com/dhamidi/sai/ebnf/lex"
	"github.com/dhamidi/sai/ebnf/parse"
	"github.com/spf13/cobra"
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
	cmd.AddCommand(newEbnfParseCmd())

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

			g, err := grammar.Parse(filename, f)
			if err != nil {
				printErrors(err)
				return err
			}

			if err := grammar.Verify(g, startProduction); err != nil {
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

			g, err := lex.LoadGrammar(grammarFile)
			if err != nil {
				return err
			}

			input, err := io.ReadAll(bufio.NewReader(os.Stdin))
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}

			lexer := lex.NewLexer(g, input, "<stdin>")
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

func newEbnfParseCmd() *cobra.Command {
	var startProduction string
	var lexerGrammarFile string

	cmd := &cobra.Command{
		Use:           "parse <parser-grammar>",
		Short:         "Parse input based on EBNF grammars",
		Long:          "Reads input from stdin, tokenizes with lexer grammar, then parses with parser grammar.\nOutputs a concrete syntax tree.",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			parserGrammarFile := args[0]

			// Load parser grammar
			parserGrammar, err := parse.LoadGrammar(parserGrammarFile)
			if err != nil {
				return fmt.Errorf("load parser grammar: %w", err)
			}

			// Load lexer grammar (defaults to parser grammar if not specified)
			var lexerGrammar grammar.Grammar
			if lexerGrammarFile != "" {
				lexerGrammar, err = lex.LoadGrammar(lexerGrammarFile)
				if err != nil {
					return fmt.Errorf("load lexer grammar: %w", err)
				}
			} else {
				lexerGrammar = parserGrammar
			}

			// Read input
			input, err := io.ReadAll(bufio.NewReader(os.Stdin))
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}

			// Tokenize
			lexer := lex.NewLexer(lexerGrammar, input, "<stdin>")
			tokens, err := lexer.Tokenize()
			if err != nil && err != io.EOF {
				return fmt.Errorf("tokenize: %w", err)
			}

			// Parse
			parser := parse.NewParser(parserGrammar, tokens)
			node, err := parser.Parse(startProduction)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			if node == nil {
				fmt.Println("Parse returned nil (no match)")
				fmt.Println("Tokens:")
				for _, t := range tokens {
					fmt.Printf("  %s\n", t)
				}
				return nil
			}

			// Print CST
			printCST(node, 0)
			return nil
		},
	}

	cmd.Flags().StringVar(&startProduction, "start", "CompilationUnit", "start production for parsing")
	cmd.Flags().StringVar(&lexerGrammarFile, "lexer", "", "lexer grammar file (defaults to parser grammar)")

	return cmd
}

func printCST(node *parse.Node, indent int) {
	if node == nil {
		return
	}

	// Skip empty non-terminal nodes (from zero-match repetitions/options)
	if !node.IsTerminal() && !node.IsError() && len(node.Children) == 0 {
		return
	}

	prefix := strings.Repeat("  ", indent)

	if node.IsError() {
		fmt.Printf("%sERROR: %s [%s]\n", prefix, node.Error, node.Span.Start)
		return
	}

	if node.IsTerminal() {
		fmt.Printf("%s%s %q [%s]\n", prefix, node.Kind, node.Token.Literal, node.Span.Start)
		return
	}

	fmt.Printf("%s%s [%s-%s]\n", prefix, node.Kind, node.Span.Start, node.Span.End)
	for _, child := range node.Children {
		printCST(child, indent+1)
	}
}
