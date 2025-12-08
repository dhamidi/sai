package pom

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Scope string

const (
	ScopeCompile  Scope = "compile"
	ScopeProvided Scope = "provided"
	ScopeRuntime  Scope = "runtime"
	ScopeTest     Scope = "test"
	ScopeSystem   Scope = "system"
)

type ArtifactKey struct {
	GroupID    string
	ArtifactID string
}

func (k ArtifactKey) String() string {
	return k.GroupID + ":" + k.ArtifactID
}

type ResolvedDependency struct {
	GroupID    string
	ArtifactID string
	Version    string
	Scope      Scope
	Type       string
	Classifier string
	Optional   bool
	Depth      int
}

func (d ResolvedDependency) Key() ArtifactKey {
	return ArtifactKey{GroupID: d.GroupID, ArtifactID: d.ArtifactID}
}

type VersionRequirement struct {
	Raw    string
	IsHard bool
	Ranges []VersionRange
	Soft   *Version
}

type VersionRange struct {
	Min          *Version
	Max          *Version
	MinInclusive bool
	MaxInclusive bool
}

func (r VersionRange) Contains(v *Version) bool {
	if r.Min != nil {
		cmp := CompareVersions(v, r.Min)
		if r.MinInclusive && cmp < 0 {
			return false
		}
		if !r.MinInclusive && cmp <= 0 {
			return false
		}
	}
	if r.Max != nil {
		cmp := CompareVersions(v, r.Max)
		if r.MaxInclusive && cmp > 0 {
			return false
		}
		if !r.MaxInclusive && cmp >= 0 {
			return false
		}
	}
	return true
}

type Version struct {
	Raw    string
	Tokens []VersionToken
}

type VersionToken struct {
	Value     string
	IsNumeric bool
	Separator string
}

type POMFetcher interface {
	FetchPOM(groupID, artifactID, version string) (*Project, error)
}

type Resolver struct {
	fetcher POMFetcher
}

func NewResolver(fetcher POMFetcher) *Resolver {
	return &Resolver{fetcher: fetcher}
}

type ExclusionSet map[ArtifactKey]struct{}

func (e ExclusionSet) Contains(key ArtifactKey) bool {
	if _, ok := e[key]; ok {
		return true
	}
	wildcard := ArtifactKey{GroupID: "*", ArtifactID: "*"}
	if _, ok := e[wildcard]; ok {
		return true
	}
	groupWildcard := ArtifactKey{GroupID: key.GroupID, ArtifactID: "*"}
	if _, ok := e[groupWildcard]; ok {
		return true
	}
	return false
}

func (e ExclusionSet) Add(key ArtifactKey) {
	e[key] = struct{}{}
}

func (e ExclusionSet) Merge(other ExclusionSet) ExclusionSet {
	result := make(ExclusionSet)
	for k := range e {
		result[k] = struct{}{}
	}
	for k := range other {
		result[k] = struct{}{}
	}
	return result
}

type resolutionState struct {
	resolved     map[ArtifactKey]*ResolvedDependency
	requirements map[ArtifactKey][]*VersionRequirement
	firstSeen    map[ArtifactKey]int
}

func (r *Resolver) Resolve(project *Project) ([]ResolvedDependency, error) {
	state := &resolutionState{
		resolved:     make(map[ArtifactKey]*ResolvedDependency),
		requirements: make(map[ArtifactKey][]*VersionRequirement),
		firstSeen:    make(map[ArtifactKey]int),
	}

	exclusions := make(ExclusionSet)
	err := r.resolveTransitive(project, ScopeCompile, 0, exclusions, state)
	if err != nil {
		return nil, err
	}

	if err := r.mediateVersions(state); err != nil {
		return nil, err
	}

	result := make([]ResolvedDependency, 0, len(state.resolved))
	for _, dep := range state.resolved {
		result = append(result, *dep)
	}
	return result, nil
}

func (r *Resolver) resolveTransitive(
	project *Project,
	parentScope Scope,
	depth int,
	exclusions ExclusionSet,
	state *resolutionState,
) error {
	depMgmt := make(map[ArtifactKey]Dependency)
	if project.DependencyManagement != nil {
		for _, d := range project.DependencyManagement.Dependencies {
			key := ArtifactKey{GroupID: d.GroupID, ArtifactID: d.ArtifactID}
			depMgmt[key] = d
		}
	}

	for _, dep := range project.Dependencies {
		key := ArtifactKey{GroupID: dep.GroupID, ArtifactID: dep.ArtifactID}

		if exclusions.Contains(key) {
			continue
		}

		if dep.Optional == "true" && depth > 0 {
			continue
		}

		scope := resolveScope(dep.Scope, parentScope)
		if scope == "" {
			continue
		}

		version := dep.Version
		if version == "" {
			if managed, ok := depMgmt[key]; ok {
				version = managed.Version
			}
		}

		req, err := ParseVersionRequirement(version)
		if err != nil {
			return fmt.Errorf("invalid version requirement %q for %s: %w", version, key, err)
		}

		if _, exists := state.firstSeen[key]; !exists {
			state.firstSeen[key] = depth
		}

		state.requirements[key] = append(state.requirements[key], req)

		if existing, ok := state.resolved[key]; ok {
			if depth < existing.Depth {
				existing.Depth = depth
			}
			continue
		}

		resolved := &ResolvedDependency{
			GroupID:    dep.GroupID,
			ArtifactID: dep.ArtifactID,
			Version:    version,
			Scope:      scope,
			Type:       dep.Type,
			Classifier: dep.Classifier,
			Optional:   dep.Optional == "true",
			Depth:      depth,
		}
		state.resolved[key] = resolved

		childExclusions := exclusions.Merge(buildExclusionSet(dep.Exclusions))

		if r.fetcher != nil && shouldResolveTransitive(scope) {
			childPOM, err := r.fetcher.FetchPOM(dep.GroupID, dep.ArtifactID, version)
			if err != nil || childPOM == nil {
				continue
			}
			if err := r.resolveTransitive(childPOM, scope, depth+1, childExclusions, state); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildExclusionSet(exclusions []Exclusion) ExclusionSet {
	set := make(ExclusionSet)
	for _, ex := range exclusions {
		key := ArtifactKey{GroupID: ex.GroupID, ArtifactID: ex.ArtifactID}
		set.Add(key)
	}
	return set
}

func shouldResolveTransitive(scope Scope) bool {
	switch scope {
	case ScopeCompile, ScopeRuntime:
		return true
	default:
		return false
	}
}

func resolveScope(depScope string, parentScope Scope) Scope {
	if depScope == "" {
		depScope = string(ScopeCompile)
	}
	scope := Scope(depScope)

	switch scope {
	case ScopeTest:
		if parentScope != "" && parentScope != ScopeTest {
			return ""
		}
		return ScopeTest
	case ScopeProvided:
		return ""
	case ScopeRuntime:
		if parentScope == ScopeCompile {
			return ScopeRuntime
		}
		return parentScope
	case ScopeCompile:
		if parentScope == "" {
			return ScopeCompile
		}
		return parentScope
	case ScopeSystem:
		return ""
	default:
		return ScopeCompile
	}
}

func (r *Resolver) mediateVersions(state *resolutionState) error {
	for key, reqs := range state.requirements {
		resolved, ok := state.resolved[key]
		if !ok {
			continue
		}

		version, err := mediateVersion(reqs)
		if err != nil {
			return fmt.Errorf("version conflict for %s: %w", key, err)
		}
		resolved.Version = version
	}
	return nil
}

func mediateVersion(reqs []*VersionRequirement) (string, error) {
	var hardReqs []*VersionRequirement
	var softReqs []*VersionRequirement

	for _, req := range reqs {
		if req.IsHard {
			hardReqs = append(hardReqs, req)
		} else {
			softReqs = append(softReqs, req)
		}
	}

	if len(hardReqs) > 0 {
		return mediateHardRequirements(hardReqs)
	}

	if len(softReqs) > 0 {
		return softReqs[0].Soft.Raw, nil
	}

	return "", errors.New("no version requirements")
}

func mediateHardRequirements(reqs []*VersionRequirement) (string, error) {
	var exactVersions []*Version
	var ranges []VersionRange

	for _, req := range reqs {
		for _, r := range req.Ranges {
			if r.Min != nil && r.Max != nil && r.Min.Raw == r.Max.Raw && r.MinInclusive && r.MaxInclusive {
				exactVersions = append(exactVersions, r.Min)
			} else {
				ranges = append(ranges, r)
			}
		}
	}

	if len(exactVersions) > 0 {
		v := exactVersions[0]
		for _, other := range exactVersions[1:] {
			if v.Raw != other.Raw {
				return "", fmt.Errorf("conflicting exact versions: %s vs %s", v.Raw, other.Raw)
			}
		}
		for _, r := range ranges {
			if !r.Contains(v) {
				return "", fmt.Errorf("exact version %s does not satisfy range", v.Raw)
			}
		}
		return v.Raw, nil
	}

	if len(ranges) == 0 {
		return "", errors.New("no version ranges")
	}

	candidates := findCandidateVersions(ranges)
	if len(candidates) == 0 {
		return "", errors.New("no version satisfies all constraints")
	}

	best := candidates[0]
	for _, c := range candidates[1:] {
		if CompareVersions(c, best) > 0 {
			best = c
		}
	}
	return best.Raw, nil
}

func findCandidateVersions(ranges []VersionRange) []*Version {
	var candidates []*Version

	for _, r := range ranges {
		if r.Min != nil {
			if satisfiesAllRanges(r.Min, ranges) {
				candidates = append(candidates, r.Min)
			}
		}
		if r.Max != nil {
			if satisfiesAllRanges(r.Max, ranges) {
				candidates = append(candidates, r.Max)
			}
		}
	}

	return candidates
}

func satisfiesAllRanges(v *Version, ranges []VersionRange) bool {
	for _, r := range ranges {
		if !r.Contains(v) {
			return false
		}
	}
	return true
}

var rangePattern = regexp.MustCompile(`^[\[\(].*[\]\)]$`)

func ParseVersionRequirement(s string) (*VersionRequirement, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty version requirement")
	}

	if !strings.ContainsAny(s, "[](,)") {
		v := ParseVersion(s)
		return &VersionRequirement{
			Raw:    s,
			IsHard: false,
			Soft:   v,
		}, nil
	}

	ranges, err := parseVersionRanges(s)
	if err != nil {
		return nil, err
	}

	return &VersionRequirement{
		Raw:    s,
		IsHard: true,
		Ranges: ranges,
	}, nil
}

func parseVersionRanges(s string) ([]VersionRange, error) {
	var ranges []VersionRange

	parts := splitRanges(s)
	for _, part := range parts {
		r, err := parseVersionRange(part)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, r)
	}

	return ranges, nil
}

func splitRanges(s string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for _, c := range s {
		switch c {
		case '[', '(':
			if depth == 0 && current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			current.WriteRune(c)
			depth++
		case ']', ')':
			current.WriteRune(c)
			depth--
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		case ',':
			if depth == 0 {
				continue
			}
			current.WriteRune(c)
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func parseVersionRange(s string) (VersionRange, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return VersionRange{}, fmt.Errorf("invalid range: %s", s)
	}

	minInclusive := s[0] == '['
	maxInclusive := s[len(s)-1] == ']'

	inner := s[1 : len(s)-1]
	parts := strings.SplitN(inner, ",", 2)

	var r VersionRange
	r.MinInclusive = minInclusive
	r.MaxInclusive = maxInclusive

	if len(parts) == 1 {
		v := ParseVersion(strings.TrimSpace(parts[0]))
		r.Min = v
		r.Max = v
		r.MinInclusive = true
		r.MaxInclusive = true
		return r, nil
	}

	minStr := strings.TrimSpace(parts[0])
	maxStr := strings.TrimSpace(parts[1])

	if minStr != "" {
		r.Min = ParseVersion(minStr)
	}
	if maxStr != "" {
		r.Max = ParseVersion(maxStr)
	}

	return r, nil
}

func ParseVersion(s string) *Version {
	s = strings.TrimSpace(s)
	tokens := tokenizeVersion(s)
	tokens = trimNullTokens(tokens)

	return &Version{
		Raw:    s,
		Tokens: tokens,
	}
}

func tokenizeVersion(s string) []VersionToken {
	var tokens []VersionToken
	var current strings.Builder
	var currentIsNumeric bool
	separator := ""

	for i, c := range s {
		switch {
		case c == '.' || c == '-' || c == '_':
			if current.Len() > 0 {
				tokens = append(tokens, VersionToken{
					Value:     current.String(),
					IsNumeric: currentIsNumeric,
					Separator: separator,
				})
				current.Reset()
			}
			separator = string(c)
		case unicode.IsDigit(c):
			if current.Len() > 0 && !currentIsNumeric {
				tokens = append(tokens, VersionToken{
					Value:     current.String(),
					IsNumeric: false,
					Separator: separator,
				})
				current.Reset()
				separator = ""
			}
			current.WriteRune(c)
			currentIsNumeric = true
		case unicode.IsLetter(c):
			if current.Len() > 0 && currentIsNumeric {
				tokens = append(tokens, VersionToken{
					Value:     current.String(),
					IsNumeric: true,
					Separator: separator,
				})
				current.Reset()
				separator = ""
			}
			current.WriteRune(c)
			currentIsNumeric = false
		default:
			if i == 0 {
				currentIsNumeric = unicode.IsDigit(c)
			}
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, VersionToken{
			Value:     current.String(),
			IsNumeric: currentIsNumeric,
			Separator: separator,
		})
	}

	return tokens
}

func trimNullTokens(tokens []VersionToken) []VersionToken {
	for len(tokens) > 0 {
		last := tokens[len(tokens)-1]
		if isNullToken(last) {
			tokens = tokens[:len(tokens)-1]
		} else {
			break
		}
	}
	return tokens
}

func isNullToken(t VersionToken) bool {
	if t.IsNumeric {
		return t.Value == "0"
	}
	lower := strings.ToLower(t.Value)
	return lower == "" || lower == "final" || lower == "ga" || lower == "release"
}

func CompareVersions(a, b *Version) int {
	maxLen := len(a.Tokens)
	if len(b.Tokens) > maxLen {
		maxLen = len(b.Tokens)
	}

	for i := 0; i < maxLen; i++ {
		var tokA, tokB VersionToken
		var aExists, bExists bool

		if i < len(a.Tokens) {
			tokA = a.Tokens[i]
			aExists = true
		} else {
			tokA = nullToken()
			aExists = false
		}

		if i < len(b.Tokens) {
			tokB = b.Tokens[i]
			bExists = true
		} else {
			tokB = nullToken()
			bExists = false
		}

		if aExists != bExists {
			existing := tokA
			if bExists {
				existing = tokB
			}
			if !existing.IsNumeric {
				order := qualifierOrder(strings.ToLower(existing.Value))
				releaseOrder := qualifierOrder("")
				if aExists {
					if order < releaseOrder {
						return -1
					}
					return 1
				} else {
					if order < releaseOrder {
						return 1
					}
					return -1
				}
			}
		}

		cmp := compareTokens(tokA, tokB)
		if cmp != 0 {
			return cmp
		}
	}

	return 0
}

func nullToken() VersionToken {
	return VersionToken{Value: "", IsNumeric: false}
}

func compareTokens(a, b VersionToken) int {
	if a.IsNumeric && b.IsNumeric {
		aNum, _ := strconv.ParseInt(a.Value, 10, 64)
		bNum, _ := strconv.ParseInt(b.Value, 10, 64)
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}

		return compareSeparators(a.Separator, b.Separator)
	}

	if a.IsNumeric && !b.IsNumeric {
		return 1
	}
	if !a.IsNumeric && b.IsNumeric {
		return -1
	}

	aQual := qualifierOrder(strings.ToLower(a.Value))
	bQual := qualifierOrder(strings.ToLower(b.Value))

	if aQual != bQual {
		if aQual < bQual {
			return -1
		}
		return 1
	}

	aLower := strings.ToLower(a.Value)
	bLower := strings.ToLower(b.Value)
	if aLower < bLower {
		return -1
	}
	if aLower > bLower {
		return 1
	}

	return compareSeparators(a.Separator, b.Separator)
}

func qualifierOrder(q string) int {
	switch q {
	case "alpha", "a":
		return 1
	case "beta", "b":
		return 2
	case "milestone", "m":
		return 3
	case "rc", "cr":
		return 4
	case "snapshot":
		return 5
	case "", "final", "ga", "release":
		return 6
	case "sp":
		return 7
	default:
		return 6
	}
}

func compareSeparators(a, b string) int {
	aOrder := separatorOrder(a)
	bOrder := separatorOrder(b)
	if aOrder < bOrder {
		return -1
	}
	if aOrder > bOrder {
		return 1
	}
	return 0
}

func separatorOrder(s string) int {
	switch s {
	case "-":
		return 1
	case ".":
		return 2
	default:
		return 0
	}
}
