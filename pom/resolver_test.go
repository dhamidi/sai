package pom

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0", "1.0"},
		{"1.0.0", "1.0.0"},
		{"1.0-SNAPSHOT", "1.0-SNAPSHOT"},
		{"1.0-alpha-1", "1.0-alpha-1"},
		{"1.0.0-beta.2", "1.0.0-beta.2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v := ParseVersion(tt.input)
			if v.Raw != tt.expected {
				t.Errorf("ParseVersion(%q).Raw = %q, want %q", tt.input, v.Raw, tt.expected)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"1.0", "1.0", 0},
		{"1.0", "1.1", -1},
		{"1.1", "1.0", 1},
		{"1.0", "2.0", -1},
		{"1.0.0", "1.0", 0},
		{"1.0", "1.0.0", 0},
		{"1.0-alpha", "1.0-beta", -1},
		{"1.0-beta", "1.0-alpha", 1},
		{"1.0-alpha", "1.0", -1},
		{"1.0", "1.0-alpha", 1},
		{"1.0-SNAPSHOT", "1.0", -1},
		{"1.0", "1.0-SNAPSHOT", 1},
		{"1.0-rc1", "1.0-beta1", 1},
		{"1.0-beta1", "1.0-rc1", -1},
		{"1.0-sp1", "1.0", 1},
		{"1.0", "1.0-sp1", -1},
		{"1.0.0.Final", "1.0.0", 0},
		{"1.0.0.GA", "1.0.0", 0},
		{"1.0.0.RELEASE", "1.0.0", 0},
		{"2.0", "10.0", -1},
		{"10.0", "2.0", 1},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			a := ParseVersion(tt.a)
			b := ParseVersion(tt.b)
			result := CompareVersions(a, b)
			if result != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestParseVersionRequirement_Soft(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"1.0"},
		{"1.0.0"},
		{"2.3.4"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			req, err := ParseVersionRequirement(tt.input)
			if err != nil {
				t.Fatalf("ParseVersionRequirement(%q) error = %v", tt.input, err)
			}
			if req.IsHard {
				t.Errorf("ParseVersionRequirement(%q).IsHard = true, want false", tt.input)
			}
			if req.Soft == nil {
				t.Errorf("ParseVersionRequirement(%q).Soft = nil, want non-nil", tt.input)
			}
		})
	}
}

func TestParseVersionRequirement_Hard(t *testing.T) {
	tests := []struct {
		input         string
		expectedCount int
	}{
		{"[1.0]", 1},
		{"[1.0,2.0]", 1},
		{"[1.0,2.0)", 1},
		{"(1.0,2.0]", 1},
		{"(1.0,2.0)", 1},
		{"(,1.0]", 1},
		{"[1.5,)", 1},
		{"(,1.0],[1.2,)", 2},
		{"(,1.1),(1.1,)", 2},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			req, err := ParseVersionRequirement(tt.input)
			if err != nil {
				t.Fatalf("ParseVersionRequirement(%q) error = %v", tt.input, err)
			}
			if !req.IsHard {
				t.Errorf("ParseVersionRequirement(%q).IsHard = false, want true", tt.input)
			}
			if len(req.Ranges) != tt.expectedCount {
				t.Errorf("ParseVersionRequirement(%q) has %d ranges, want %d", tt.input, len(req.Ranges), tt.expectedCount)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		rangeStr string
		version  string
		expected bool
	}{
		{"[1.0]", "1.0", true},
		{"[1.0]", "1.1", false},
		{"[1.0,2.0]", "1.0", true},
		{"[1.0,2.0]", "1.5", true},
		{"[1.0,2.0]", "2.0", true},
		{"[1.0,2.0]", "0.9", false},
		{"[1.0,2.0]", "2.1", false},
		{"(1.0,2.0)", "1.0", false},
		{"(1.0,2.0)", "2.0", false},
		{"(1.0,2.0)", "1.5", true},
		{"[1.0,2.0)", "1.0", true},
		{"[1.0,2.0)", "2.0", false},
		{"(1.0,2.0]", "1.0", false},
		{"(1.0,2.0]", "2.0", true},
		{"(,1.0]", "0.5", true},
		{"(,1.0]", "1.0", true},
		{"(,1.0]", "1.1", false},
		{"[1.5,)", "1.5", true},
		{"[1.5,)", "2.0", true},
		{"[1.5,)", "1.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.rangeStr+" contains "+tt.version, func(t *testing.T) {
			req, err := ParseVersionRequirement(tt.rangeStr)
			if err != nil {
				t.Fatalf("ParseVersionRequirement(%q) error = %v", tt.rangeStr, err)
			}
			if len(req.Ranges) != 1 {
				t.Fatalf("Expected 1 range, got %d", len(req.Ranges))
			}
			r := req.Ranges[0]
			v := ParseVersion(tt.version)
			result := r.Contains(v)
			if result != tt.expected {
				t.Errorf("Range %q contains %q = %v, want %v", tt.rangeStr, tt.version, result, tt.expected)
			}
		})
	}
}

func TestResolveScope(t *testing.T) {
	tests := []struct {
		depScope    string
		parentScope Scope
		expected    Scope
	}{
		{"compile", ScopeCompile, ScopeCompile},
		{"compile", "", ScopeCompile},
		{"", "", ScopeCompile},
		{"runtime", ScopeCompile, ScopeRuntime},
		{"runtime", ScopeRuntime, ScopeRuntime},
		{"test", "", ScopeTest},
		{"test", ScopeCompile, ""},
		{"provided", ScopeCompile, ""},
		{"system", ScopeCompile, ""},
	}

	for _, tt := range tests {
		name := tt.depScope + " in " + string(tt.parentScope)
		t.Run(name, func(t *testing.T) {
			result := resolveScope(tt.depScope, tt.parentScope)
			if result != tt.expected {
				t.Errorf("resolveScope(%q, %q) = %q, want %q", tt.depScope, tt.parentScope, result, tt.expected)
			}
		})
	}
}

func TestExclusionSet(t *testing.T) {
	set := make(ExclusionSet)
	key := ArtifactKey{GroupID: "org.example", ArtifactID: "lib"}

	if set.Contains(key) {
		t.Error("Empty set should not contain key")
	}

	set.Add(key)
	if !set.Contains(key) {
		t.Error("Set should contain added key")
	}

	other := ArtifactKey{GroupID: "org.other", ArtifactID: "lib"}
	if set.Contains(other) {
		t.Error("Set should not contain key not added")
	}
}

func TestExclusionSet_Wildcard(t *testing.T) {
	set := make(ExclusionSet)
	wildcard := ArtifactKey{GroupID: "*", ArtifactID: "*"}
	set.Add(wildcard)

	key := ArtifactKey{GroupID: "org.example", ArtifactID: "lib"}
	if !set.Contains(key) {
		t.Error("Wildcard should match any key")
	}
}

func TestExclusionSet_Merge(t *testing.T) {
	set1 := make(ExclusionSet)
	set1.Add(ArtifactKey{GroupID: "a", ArtifactID: "1"})

	set2 := make(ExclusionSet)
	set2.Add(ArtifactKey{GroupID: "b", ArtifactID: "2"})

	merged := set1.Merge(set2)

	if !merged.Contains(ArtifactKey{GroupID: "a", ArtifactID: "1"}) {
		t.Error("Merged set should contain key from set1")
	}
	if !merged.Contains(ArtifactKey{GroupID: "b", ArtifactID: "2"}) {
		t.Error("Merged set should contain key from set2")
	}
}

type mockFetcher struct {
	poms map[string]*Project
}

func (m *mockFetcher) FetchPOM(groupID, artifactID, version string) (*Project, error) {
	key := groupID + ":" + artifactID + ":" + version
	if pom, ok := m.poms[key]; ok {
		return pom, nil
	}
	return nil, nil
}

func TestResolver_BasicResolve(t *testing.T) {
	project := &Project{
		GroupID:    "com.example",
		ArtifactID: "app",
		Version:    "1.0.0",
		Dependencies: []Dependency{
			{GroupID: "org.lib", ArtifactID: "core", Version: "2.0.0"},
			{GroupID: "org.lib", ArtifactID: "util", Version: "1.5.0", Scope: "runtime"},
		},
	}

	resolver := NewResolver(nil)
	deps, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	found := make(map[string]bool)
	for _, d := range deps {
		found[d.ArtifactID] = true
	}

	if !found["core"] {
		t.Error("Expected core dependency")
	}
	if !found["util"] {
		t.Error("Expected util dependency")
	}
}

func TestResolver_Exclusions(t *testing.T) {
	project := &Project{
		GroupID:    "com.example",
		ArtifactID: "app",
		Version:    "1.0.0",
		Dependencies: []Dependency{
			{
				GroupID:    "org.lib",
				ArtifactID: "core",
				Version:    "2.0.0",
				Exclusions: []Exclusion{
					{GroupID: "org.excluded", ArtifactID: "lib"},
				},
			},
		},
	}

	fetcher := &mockFetcher{
		poms: map[string]*Project{
			"org.lib:core:2.0.0": {
				GroupID:    "org.lib",
				ArtifactID: "core",
				Version:    "2.0.0",
				Dependencies: []Dependency{
					{GroupID: "org.excluded", ArtifactID: "lib", Version: "1.0.0"},
					{GroupID: "org.included", ArtifactID: "lib", Version: "1.0.0"},
				},
			},
		},
	}

	resolver := NewResolver(fetcher)
	deps, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	for _, d := range deps {
		if d.GroupID == "org.excluded" {
			t.Error("Excluded dependency should not be resolved")
		}
	}

	found := false
	for _, d := range deps {
		if d.GroupID == "org.included" {
			found = true
		}
	}
	if !found {
		t.Error("Non-excluded dependency should be resolved")
	}
}

func TestResolver_OptionalNotTransitive(t *testing.T) {
	project := &Project{
		GroupID:    "com.example",
		ArtifactID: "app",
		Version:    "1.0.0",
		Dependencies: []Dependency{
			{GroupID: "org.lib", ArtifactID: "core", Version: "2.0.0"},
		},
	}

	fetcher := &mockFetcher{
		poms: map[string]*Project{
			"org.lib:core:2.0.0": {
				GroupID:    "org.lib",
				ArtifactID: "core",
				Version:    "2.0.0",
				Dependencies: []Dependency{
					{GroupID: "org.optional", ArtifactID: "lib", Version: "1.0.0", Optional: "true"},
				},
			},
		},
	}

	resolver := NewResolver(fetcher)
	deps, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	for _, d := range deps {
		if d.GroupID == "org.optional" {
			t.Error("Optional transitive dependency should not be resolved")
		}
	}
}

func TestResolver_VersionMediation_SoftFirst(t *testing.T) {
	fetcher := &mockFetcher{
		poms: map[string]*Project{
			"org.lib:a:1.0.0": {
				GroupID:    "org.lib",
				ArtifactID: "a",
				Version:    "1.0.0",
				Dependencies: []Dependency{
					{GroupID: "org.shared", ArtifactID: "lib", Version: "2.0.0"},
				},
			},
		},
	}

	projectWithDirect := &Project{
		GroupID:    "com.example",
		ArtifactID: "app",
		Version:    "1.0.0",
		Dependencies: []Dependency{
			{GroupID: "org.shared", ArtifactID: "lib", Version: "1.0.0"},
			{GroupID: "org.lib", ArtifactID: "a", Version: "1.0.0"},
		},
	}

	resolver := NewResolver(fetcher)
	deps, err := resolver.Resolve(projectWithDirect)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	for _, d := range deps {
		if d.GroupID == "org.shared" && d.ArtifactID == "lib" {
			if d.Version != "1.0.0" {
				t.Errorf("Expected version 1.0.0 (first seen), got %s", d.Version)
			}
		}
	}
}

func TestResolver_TransitiveScope(t *testing.T) {
	project := &Project{
		GroupID:    "com.example",
		ArtifactID: "app",
		Version:    "1.0.0",
		Dependencies: []Dependency{
			{GroupID: "org.lib", ArtifactID: "core", Version: "1.0.0"},
		},
	}

	fetcher := &mockFetcher{
		poms: map[string]*Project{
			"org.lib:core:1.0.0": {
				GroupID:    "org.lib",
				ArtifactID: "core",
				Version:    "1.0.0",
				Dependencies: []Dependency{
					{GroupID: "org.test", ArtifactID: "lib", Version: "1.0.0", Scope: "test"},
					{GroupID: "org.provided", ArtifactID: "lib", Version: "1.0.0", Scope: "provided"},
					{GroupID: "org.compile", ArtifactID: "lib", Version: "1.0.0", Scope: "compile"},
				},
			},
		},
	}

	resolver := NewResolver(fetcher)
	deps, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	for _, d := range deps {
		if d.GroupID == "org.test" {
			t.Error("Test-scoped transitive dependency should not be included")
		}
		if d.GroupID == "org.provided" {
			t.Error("Provided-scoped transitive dependency should not be included")
		}
	}

	found := false
	for _, d := range deps {
		if d.GroupID == "org.compile" {
			found = true
		}
	}
	if !found {
		t.Error("Compile-scoped transitive dependency should be included")
	}
}

func TestMediateVersion_HardRequirements(t *testing.T) {
	tests := []struct {
		name        string
		reqs        []string
		expected    string
		shouldError bool
	}{
		{
			name:     "single exact version",
			reqs:     []string{"[1.0.0]"},
			expected: "1.0.0",
		},
		{
			name:     "compatible range and exact",
			reqs:     []string{"[1.0.0,2.0.0]", "[1.5.0]"},
			expected: "1.5.0",
		},
		{
			name:        "conflicting exact versions",
			reqs:        []string{"[1.0.0]", "[2.0.0]"},
			shouldError: true,
		},
		{
			name:        "exact version outside range",
			reqs:        []string{"[1.0.0]", "(1.0.0,2.0.0)"},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqs []*VersionRequirement
			for _, r := range tt.reqs {
				req, err := ParseVersionRequirement(r)
				if err != nil {
					t.Fatalf("ParseVersionRequirement(%q) error = %v", r, err)
				}
				reqs = append(reqs, req)
			}

			result, err := mediateVersion(reqs)
			if tt.shouldError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("mediateVersion() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("mediateVersion() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

func TestDependencyManagement(t *testing.T) {
	project := &Project{
		GroupID:    "com.example",
		ArtifactID: "app",
		Version:    "1.0.0",
		DependencyManagement: &DependencyManagement{
			Dependencies: []Dependency{
				{GroupID: "org.managed", ArtifactID: "lib", Version: "3.0.0"},
			},
		},
		Dependencies: []Dependency{
			{GroupID: "org.managed", ArtifactID: "lib"},
		},
	}

	resolver := NewResolver(nil)
	deps, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if len(deps) != 1 {
		t.Fatalf("Expected 1 dependency, got %d", len(deps))
	}

	if deps[0].Version != "3.0.0" {
		t.Errorf("Expected version from dependencyManagement (3.0.0), got %s", deps[0].Version)
	}
}
