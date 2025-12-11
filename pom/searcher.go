package pom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const DefaultSearchURL = "https://search.maven.org/solrsearch/select"

type Searcher struct {
	BaseURL    string
	httpClient *http.Client
}

func NewSearcher() *Searcher {
	return &Searcher{
		BaseURL:    DefaultSearchURL,
		httpClient: &http.Client{},
	}
}

type SearchQuery struct {
	Text                    string
	GroupID                 string
	ArtifactID              string
	Version                 string
	Packaging               string
	Classifier              string
	ClassName               string
	FullyQualifiedClassName string
	Tags                    string
	Rows                    int
	Start                   int
	Core                    string // "gav" for all versions
}

type SearchResponse struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Response       Response       `json:"response"`
}

type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
}

type Response struct {
	NumFound int         `json:"numFound"`
	Start    int         `json:"start"`
	Docs     []SearchDoc `json:"docs"`
}

type SearchDoc struct {
	ID            string   `json:"id"`
	GroupID       string   `json:"g"`
	ArtifactID    string   `json:"a"`
	Version       string   `json:"v"`
	LatestVersion string   `json:"latestVersion"`
	Packaging     string   `json:"p"`
	Timestamp     int64    `json:"timestamp"`
	VersionCount  int      `json:"versionCount"`
	Text          []string `json:"text"`
	EC            []string `json:"ec"` // extension classifier
	Tags          []string `json:"tags"`
}

func (s *Searcher) Search(q SearchQuery) (*SearchResponse, error) {
	queryStr := s.buildQuery(q)
	reqURL := fmt.Sprintf("%s?q=%s&wt=json", s.BaseURL, url.QueryEscape(queryStr))

	if q.Rows > 0 {
		reqURL += fmt.Sprintf("&rows=%d", q.Rows)
	} else {
		reqURL += "&rows=20"
	}
	if q.Start > 0 {
		reqURL += fmt.Sprintf("&start=%d", q.Start)
	}
	if q.Core != "" {
		reqURL += fmt.Sprintf("&core=%s", url.QueryEscape(q.Core))
	}

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request returned HTTP %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &result, nil
}

func (s *Searcher) buildQuery(q SearchQuery) string {
	var parts []string

	if q.Text != "" {
		parts = append(parts, q.Text)
	}
	if q.GroupID != "" {
		parts = append(parts, fmt.Sprintf("g:%s", q.GroupID))
	}
	if q.ArtifactID != "" {
		parts = append(parts, fmt.Sprintf("a:%s", q.ArtifactID))
	}
	if q.Version != "" {
		parts = append(parts, fmt.Sprintf("v:%s", q.Version))
	}
	if q.Packaging != "" {
		parts = append(parts, fmt.Sprintf("p:%s", q.Packaging))
	}
	if q.Classifier != "" {
		parts = append(parts, fmt.Sprintf("l:%s", q.Classifier))
	}
	if q.ClassName != "" {
		parts = append(parts, fmt.Sprintf("c:%s", q.ClassName))
	}
	if q.FullyQualifiedClassName != "" {
		parts = append(parts, fmt.Sprintf("fc:%s", q.FullyQualifiedClassName))
	}
	if q.Tags != "" {
		parts = append(parts, fmt.Sprintf("tags:%s", q.Tags))
	}

	if len(parts) == 0 {
		return "*:*"
	}
	return strings.Join(parts, " AND ")
}

func (s *Searcher) SearchText(text string) (*SearchResponse, error) {
	return s.Search(SearchQuery{Text: text, Rows: 20})
}

func (s *Searcher) SearchByGroupID(groupID string) (*SearchResponse, error) {
	return s.Search(SearchQuery{GroupID: groupID, Rows: 20})
}

func (s *Searcher) SearchByArtifactID(artifactID string) (*SearchResponse, error) {
	return s.Search(SearchQuery{ArtifactID: artifactID, Rows: 20})
}

func (s *Searcher) SearchAllVersions(groupID, artifactID string) (*SearchResponse, error) {
	return s.Search(SearchQuery{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Core:       "gav",
		Rows:       100,
	})
}

func (s *Searcher) SearchByClassName(className string) (*SearchResponse, error) {
	return s.Search(SearchQuery{ClassName: className, Rows: 20})
}
