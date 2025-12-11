package main

import (
	"fmt"
	"strings"

	"github.com/dhamidi/sai/pom"
	"github.com/spf13/cobra"
)

func newLibsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "libs",
		Short: "Manage library dependencies",
	}

	cmd.AddCommand(newLibsSearchCmd())
	cmd.AddCommand(newLibsListCmd())

	return cmd
}

func newLibsSearchCmd() *cobra.Command {
	var (
		groupID    string
		artifactID string
		className  string
		rows       int
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for Maven dependencies",
		Long: `Search Maven Central for dependencies.

Examples:
  sai libs search guice
  sai libs search --group com.google.inject
  sai libs search --artifact guice
  sai libs search --class Injector`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			searcher := pom.NewSearcher()

			query := pom.SearchQuery{
				GroupID:    groupID,
				ArtifactID: artifactID,
				ClassName:  className,
				Rows:       rows,
			}

			if len(args) > 0 {
				query.Text = args[0]
			}

			if query.Text == "" && query.GroupID == "" && query.ArtifactID == "" && query.ClassName == "" {
				return fmt.Errorf("provide a search query or use --group, --artifact, or --class flags")
			}

			result, err := searcher.Search(query)
			if err != nil {
				return err
			}

			if result.Response.NumFound == 0 {
				fmt.Println("No results found.")
				return nil
			}

			fmt.Printf("Found %d results:\n\n", result.Response.NumFound)
			for _, doc := range result.Response.Docs {
				version := doc.LatestVersion
				if version == "" {
					version = doc.Version
				}
				coord := fmt.Sprintf("%s:%s:%s", doc.GroupID, doc.ArtifactID, version)
				fmt.Printf("  %s\n", coord)
				if len(doc.Tags) > 0 {
					fmt.Printf("    tags: %s\n", strings.Join(doc.Tags, ", "))
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&groupID, "group", "g", "", "filter by group ID")
	cmd.Flags().StringVarP(&artifactID, "artifact", "a", "", "filter by artifact ID")
	cmd.Flags().StringVarP(&className, "class", "c", "", "search by class name")
	cmd.Flags().IntVarP(&rows, "rows", "n", 20, "number of results to return")

	return cmd
}
