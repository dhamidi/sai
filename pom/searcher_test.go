package pom

import "testing"

func TestSearcher_buildQuery(t *testing.T) {
	s := NewSearcher()

	tests := []struct {
		name  string
		query SearchQuery
		want  string
	}{
		{
			name:  "empty query",
			query: SearchQuery{},
			want:  "*:*",
		},
		{
			name:  "text only",
			query: SearchQuery{Text: "guice"},
			want:  "guice",
		},
		{
			name:  "groupId only",
			query: SearchQuery{GroupID: "com.google.inject"},
			want:  "g:com.google.inject",
		},
		{
			name:  "artifactId only",
			query: SearchQuery{ArtifactID: "guice"},
			want:  "a:guice",
		},
		{
			name: "groupId and artifactId",
			query: SearchQuery{
				GroupID:    "com.google.inject",
				ArtifactID: "guice",
			},
			want: "g:com.google.inject AND a:guice",
		},
		{
			name: "full coordinate search",
			query: SearchQuery{
				GroupID:    "com.google.inject",
				ArtifactID: "guice",
				Version:    "3.0",
				Packaging:  "jar",
				Classifier: "javadoc",
			},
			want: "g:com.google.inject AND a:guice AND v:3.0 AND p:jar AND l:javadoc",
		},
		{
			name:  "class name search",
			query: SearchQuery{ClassName: "junit"},
			want:  "c:junit",
		},
		{
			name:  "fully qualified class name",
			query: SearchQuery{FullyQualifiedClassName: "org.specs.runner.JUnit"},
			want:  "fc:org.specs.runner.JUnit",
		},
		{
			name:  "tags search",
			query: SearchQuery{Tags: "sbtplugin"},
			want:  "tags:sbtplugin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.buildQuery(tt.query)
			if got != tt.want {
				t.Errorf("buildQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}
