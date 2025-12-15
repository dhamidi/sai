package format

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/dhamidi/sai/java/parser"
)

var testcasesDir string
var testFilter string

func init() {
	flag.StringVar(&testcasesDir, "testcases", "", "directory containing .java test files")
	flag.StringVar(&testFilter, "filter", "", "filter test files by substring match on filename")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// TestRoundTrip_Testcases runs round-trip tests on all .java files in the testcases directory.
// Each file becomes a subtest that can be targeted with: go test -run TestRoundTrip_Testcases/filename
// Use -filter to filter files by substring: go test ./format -filter=String
// Skipped during pre-commit hooks when IN_GIT_PRECOMMIT=1
func TestRoundTrip_Testcases(t *testing.T) {
	if os.Getenv("IN_GIT_PRECOMMIT") == "1" {
		t.Skip("skipping roundtrip tests during pre-commit")
	}

	dir := testcasesDir
	if dir == "" {
		// Default to testcases/ relative to project root
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		// Walk up to find testcases directory
		for d := wd; d != "/"; d = filepath.Dir(d) {
			candidate := filepath.Join(d, "testcases")
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				dir = candidate
				break
			}
		}
		if dir == "" {
			t.Skip("testcases directory not found; use -testcases flag to specify")
		}
	}

	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".java") {
			// Apply filter if specified
			if testFilter != "" && !strings.Contains(path, testFilter) {
				return nil
			}
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk testcases directory: %v", err)
	}

	if len(files) == 0 {
		if testFilter != "" {
			t.Skipf("no .java files matching filter %q found in %s", testFilter, dir)
		}
		t.Skipf("no .java files found in %s", dir)
	}

	for _, file := range files {
		// Create a test name from the relative path
		relPath, err := filepath.Rel(dir, file)
		if err != nil {
			relPath = filepath.Base(file)
		}
		// Replace path separators with underscores for test naming
		testName := strings.ReplaceAll(relPath, string(filepath.Separator), "_")
		testName = strings.TrimSuffix(testName, ".java")

		t.Run(testName, func(t *testing.T) {
			runRoundTripTest(t, file)
		})
	}
}

func runRoundTripTest(t *testing.T, filename string) {
	source, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// Parse original
	origParser := parser.ParseCompilationUnit(bytes.NewReader(source), parser.WithComments())
	origAST := origParser.Finish()
	if origAST == nil {
		t.Fatalf("failed to parse original file")
	}
	if hasParseErrors(origAST) {
		// Skip files that don't parse cleanly - we can't test the formatter on broken input
		t.Skipf("original file has parse errors")
	}

	// Format
	formatted, err := PrettyPrintJava(source)
	if err != nil {
		t.Fatalf("formatter error: %v", err)
	}

	// Parse formatted output
	fmtParser := parser.ParseCompilationUnit(bytes.NewReader(formatted), parser.WithComments())
	fmtAST := fmtParser.Finish()
	if fmtAST == nil {
		t.Fatalf("failed to parse formatted output")
	}
	if hasParseErrors(fmtAST) {
		t.Errorf("formatted output has parse errors!\n\nOriginal parsed cleanly, but after formatting:\n%s",
			formatParseErrors(fmtAST))
		t.Logf("\n=== Formatted output ===\n%s", string(formatted))
		return
	}

	// Compare node counts
	origCounts := countNodeKinds(origAST)
	fmtCounts := countNodeKinds(fmtAST)

	diffs := compareNodeCounts(origCounts, fmtCounts)
	if len(diffs) > 0 {
		t.Errorf("node count mismatch after round-trip formatting:\n\n%s", formatDiffs(diffs))

		// Show where nodes were dropped for each kind that has fewer nodes
		t.Logf("\n=== Likely dropped nodes (by token analysis) ===")
		for _, diff := range diffs {
			if diff.Formatted < diff.Original {
				missing := findMissingNodes(origAST, fmtAST, diff.Kind)
				t.Logf("\n%s (-%d):\n%s",
					diff.Kind.String(),
					diff.Original-diff.Formatted,
					formatMissingNodes(missing, source, 5))
			}
		}
	}
}

// NodeCountDiff represents a difference in node counts between original and formatted AST
type NodeCountDiff struct {
	Kind      parser.NodeKind
	Original  int
	Formatted int
}

func countNodeKinds(node *parser.Node) map[parser.NodeKind]int {
	counts := make(map[parser.NodeKind]int)
	walkAST(node, func(n *parser.Node) {
		counts[n.Kind]++
	})
	return counts
}

func walkAST(node *parser.Node, visit func(*parser.Node)) {
	if node == nil {
		return
	}
	visit(node)
	for _, child := range node.Children {
		walkAST(child, visit)
	}
}

func hasParseErrors(node *parser.Node) bool {
	hasError := false
	walkAST(node, func(n *parser.Node) {
		if n.Kind == parser.KindError {
			hasError = true
		}
	})
	return hasError
}

func formatParseErrors(node *parser.Node) string {
	var errors []string
	walkAST(node, func(n *parser.Node) {
		if n.Kind == parser.KindError && n.Error != nil {
			errors = append(errors, fmt.Sprintf("  - %s at line %d, col %d",
				n.Error.Message, n.Span.Start.Line, n.Span.Start.Column))
		}
	})
	if len(errors) == 0 {
		return "  (error nodes found but no error details)"
	}
	return strings.Join(errors, "\n")
}

func compareNodeCounts(original, formatted map[parser.NodeKind]int) []NodeCountDiff {
	var diffs []NodeCountDiff

	// Collect all kinds from both maps
	allKinds := make(map[parser.NodeKind]bool)
	for k := range original {
		allKinds[k] = true
	}
	for k := range formatted {
		allKinds[k] = true
	}

	for kind := range allKinds {
		origCount := original[kind]
		fmtCount := formatted[kind]
		if origCount != fmtCount {
			// Skip comment nodes - formatting may legitimately change comment structure
			if kind == parser.KindComment || kind == parser.KindLineComment {
				continue
			}
			diffs = append(diffs, NodeCountDiff{
				Kind:      kind,
				Original:  origCount,
				Formatted: fmtCount,
			})
		}
	}

	// Sort by severity (most dropped nodes first)
	sort.Slice(diffs, func(i, j int) bool {
		diffI := diffs[i].Original - diffs[i].Formatted
		diffJ := diffs[j].Original - diffs[j].Formatted
		return diffI > diffJ
	})

	return diffs
}

func formatDiffs(diffs []NodeCountDiff) string {
	var sb strings.Builder
	sb.WriteString("Kind                          Original  Formatted  Delta\n")
	sb.WriteString("------------------------------------------------------------\n")
	for _, d := range diffs {
		delta := d.Formatted - d.Original
		sign := "+"
		if delta < 0 {
			sign = ""
		}
		sb.WriteString(fmt.Sprintf("%-30s %8d  %9d  %s%d\n",
			d.Kind.String(), d.Original, d.Formatted, sign, delta))
	}
	return sb.String()
}

func formatCountSummary(counts map[parser.NodeKind]int) string {
	// Sort by kind name for consistent output
	type kv struct {
		kind  parser.NodeKind
		count int
	}
	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].kind.String() < sorted[j].kind.String()
	})

	var sb strings.Builder
	for _, item := range sorted {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", item.kind.String(), item.count))
	}
	return sb.String()
}

// NodeLocation captures a node's location for debugging
type NodeLocation struct {
	Kind   parser.NodeKind
	Line   int
	Column int
	Token  string // Token literal if available
}

func (nl NodeLocation) String() string {
	if nl.Token != "" {
		return fmt.Sprintf("%s %q at line %d, col %d", nl.Kind.String(), nl.Token, nl.Line, nl.Column)
	}
	return fmt.Sprintf("%s at line %d, col %d", nl.Kind.String(), nl.Line, nl.Column)
}

// collectNodeLocations gathers all nodes of a specific kind with their locations
func collectNodeLocations(node *parser.Node, kind parser.NodeKind) []NodeLocation {
	var locations []NodeLocation
	walkAST(node, func(n *parser.Node) {
		if n.Kind == kind {
			loc := NodeLocation{
				Kind:   n.Kind,
				Line:   n.Span.Start.Line,
				Column: n.Span.Start.Column,
			}
			if n.Token != nil {
				loc.Token = n.Token.Literal
			}
			locations = append(locations, loc)
		}
	})
	return locations
}

// findMissingNodes compares nodes of a specific kind between original and formatted AST
// and returns nodes that appear in original but not in formatted (by token+line matching)
func findMissingNodes(origAST, fmtAST *parser.Node, kind parser.NodeKind) []NodeLocation {
	origLocs := collectNodeLocations(origAST, kind)
	fmtLocs := collectNodeLocations(fmtAST, kind)

	// Build a set of formatted tokens for quick lookup
	// We use token value as key since positions change after formatting
	fmtTokens := make(map[string]int)
	for _, loc := range fmtLocs {
		if loc.Token != "" {
			fmtTokens[loc.Token]++
		}
	}

	// Find tokens that appear more times in original than formatted
	origTokenCounts := make(map[string][]NodeLocation)
	for _, loc := range origLocs {
		if loc.Token != "" {
			origTokenCounts[loc.Token] = append(origTokenCounts[loc.Token], loc)
		}
	}

	var missing []NodeLocation
	for token, locs := range origTokenCounts {
		origCount := len(locs)
		fmtCount := fmtTokens[token]
		if origCount > fmtCount {
			// Some instances of this token were dropped
			// Report the first (origCount - fmtCount) locations
			for i := 0; i < origCount-fmtCount && i < len(locs); i++ {
				missing = append(missing, locs[i])
			}
		}
	}

	// Sort by line number
	sort.Slice(missing, func(i, j int) bool {
		if missing[i].Line != missing[j].Line {
			return missing[i].Line < missing[j].Line
		}
		return missing[i].Column < missing[j].Column
	})

	return missing
}

// formatMissingNodes formats a list of missing nodes for display
func formatMissingNodes(missing []NodeLocation, source []byte, maxShow int) string {
	if len(missing) == 0 {
		return "  (none found - positions may have shifted)\n"
	}

	lines := strings.Split(string(source), "\n")
	var sb strings.Builder

	shown := 0
	for _, loc := range missing {
		if shown >= maxShow {
			sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(missing)-maxShow))
			break
		}

		sb.WriteString(fmt.Sprintf("  - %s\n", loc.String()))

		// Show the source line for context
		if loc.Line > 0 && loc.Line <= len(lines) {
			line := lines[loc.Line-1]
			if len(line) > 120 {
				line = line[:120] + "..."
			}
			sb.WriteString(fmt.Sprintf("    > %s\n", line))
		}
		shown++
	}

	return sb.String()
}
