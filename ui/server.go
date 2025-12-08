package ui

import (
	"archive/zip"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/dhamidi/javalyzer/java"
)

//go:embed static templates
var embeddedFS embed.FS

type ScanStatus string

const (
	ScanStatusPending    ScanStatus = "pending"
	ScanStatusInProgress ScanStatus = "in_progress"
	ScanStatusCompleted  ScanStatus = "completed"
	ScanStatusFailed     ScanStatus = "failed"
)

type ScanRequest struct {
	ID         string
	Path       string
	ClassFiles []string
	ZipFile    string
	CreatedAt  time.Time
}

type ScanResult struct {
	ID        string
	Status    ScanStatus
	Request   ScanRequest
	Classes   []*java.Class
	Error     string
	StartedAt time.Time
	EndedAt   time.Time
	Progress  int
	Total     int
}

func (s *ScanResult) ProgressPercent() int {
	if s.Total == 0 {
		return 0
	}
	return (s.Progress * 100) / s.Total
}

type Scanner struct {
	mu       sync.RWMutex
	scans    map[string]*ScanResult
	requests chan ScanRequest
	nextID   int
}

func NewScanner() *Scanner {
	s := &Scanner{
		scans:    make(map[string]*ScanResult),
		requests: make(chan ScanRequest, 100),
	}
	go s.run()
	return s
}

func (s *Scanner) run() {
	for req := range s.requests {
		s.processScan(req)
	}
}

func (s *Scanner) processScan(req ScanRequest) {
	s.mu.Lock()
	result := s.scans[req.ID]
	result.Status = ScanStatusInProgress
	result.StartedAt = time.Now()
	s.mu.Unlock()

	var classes []*java.Class
	var scanErr error

	if req.Path != "" {
		classes, scanErr = s.scanDirectory(req.ID, req.Path)
	} else if len(req.ClassFiles) > 0 {
		classes, scanErr = s.scanFiles(req.ID, req.ClassFiles)
	} else if req.ZipFile != "" {
		classes, scanErr = s.scanZipFile(req.ID, req.ZipFile)
	} else {
		scanErr = fmt.Errorf("no path, class files, or zip file provided")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	result.EndedAt = time.Now()
	if scanErr != nil {
		result.Status = ScanStatusFailed
		result.Error = scanErr.Error()
	} else {
		result.Status = ScanStatusCompleted
		result.Classes = classes
	}
}

func (s *Scanner) scanDirectory(id, path string) ([]*java.Class, error) {
	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(p) == ".class" {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.scanFiles(id, files)
}

func (s *Scanner) scanFiles(id string, files []string) ([]*java.Class, error) {
	s.mu.Lock()
	s.scans[id].Total = len(files)
	s.mu.Unlock()

	var classes []*java.Class
	for i, file := range files {
		class, err := java.ParseClassFile(file)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}
		if class != nil {
			classes = append(classes, class)
		}

		s.mu.Lock()
		s.scans[id].Progress = i + 1
		s.mu.Unlock()
	}
	return classes, nil
}

func (s *Scanner) scanZipFile(id, path string) ([]*java.Class, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	var classFiles []*zip.File
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && filepath.Ext(f.Name) == ".class" {
			classFiles = append(classFiles, f)
		}
	}

	s.mu.Lock()
	s.scans[id].Total = len(classFiles)
	s.mu.Unlock()

	var classes []*java.Class
	for i, f := range classFiles {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", f.Name, err)
		}

		class, err := java.ParseClass(rc)
		rc.Close()

		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f.Name, err)
		}
		if class != nil {
			classes = append(classes, class)
		}

		s.mu.Lock()
		s.scans[id].Progress = i + 1
		s.mu.Unlock()
	}
	return classes, nil
}

func (s *Scanner) Submit(req ScanRequest) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	req.ID = fmt.Sprintf("%d", s.nextID)
	req.CreatedAt = time.Now()

	s.scans[req.ID] = &ScanResult{
		ID:      req.ID,
		Status:  ScanStatusPending,
		Request: req,
	}

	s.requests <- req
	return req.ID
}

func (s *Scanner) Get(id string) (*ScanResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.scans[id]
	return result, ok
}

func (s *Scanner) List() []*ScanResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := make([]*ScanResult, 0, len(s.scans))
	for _, r := range s.scans {
		results = append(results, r)
	}
	return results
}

type Server struct {
	scanner    *Scanner
	staticFS   fs.FS
	templates  *template.Template
	mux        *http.ServeMux
	templateFS fs.FS
	funcMap    template.FuncMap
}

func NewServer() (*Server, error) {
	staticFS := overlayFS("ui/static", mustSub(embeddedFS, "static"))
	templateFS := overlayFS("ui/templates", mustSub(embeddedFS, "templates"))

	funcMap := template.FuncMap{
		"linkifyClass": func(knownClasses map[string]bool, className string) template.HTML {
			escaped := template.HTMLEscapeString(className)
			if knownClasses[className] {
				return template.HTML(fmt.Sprintf(`<a href="/c/%s">%s</a>`, escaped, escaped))
			}
			return template.HTML(escaped)
		},
		"linkifyType": func(knownClasses map[string]bool, t java.Type) template.HTML {
			escaped := template.HTMLEscapeString(t.Name)
			var result string
			if knownClasses[t.Name] {
				result = fmt.Sprintf(`<a href="/c/%s">%s</a>`, escaped, escaped)
			} else {
				result = escaped
			}
			for i := 0; i < t.ArrayDepth; i++ {
				result += "[]"
			}
			return template.HTML(result)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	s := &Server{
		scanner:    NewScanner(),
		staticFS:   staticFS,
		templates:  tmpl,
		mux:        http.NewServeMux(),
		templateFS: templateFS,
		funcMap:    funcMap,
	}

	s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	s.mux.HandleFunc("POST /scan", s.handleScan)
	s.mux.HandleFunc("GET /scans/{id}", s.handleGetScan)
	s.mux.HandleFunc("GET /c/{className...}", s.handleClass)
	s.mux.HandleFunc("GET /", s.handleIndex)

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	tmpl, err := template.New("").Funcs(s.funcMap).ParseFS(s.templateFS, "*.html")
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, name, data)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "invalid form data: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		req.Path = r.FormValue("path")

		if classFiles := r.Form["class_files"]; len(classFiles) > 0 {
			req.ClassFiles = classFiles
		}

		if file, header, err := r.FormFile("zipfile"); err == nil {
			defer file.Close()
			tmpFile, err := os.CreateTemp("", "javalyzer-*.zip")
			if err != nil {
				http.Error(w, "failed to create temp file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := io.Copy(tmpFile, file); err != nil {
				tmpFile.Close()
				os.Remove(tmpFile.Name())
				http.Error(w, "failed to save zip file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			tmpFile.Close()
			req.ZipFile = tmpFile.Name()
			_ = header
		}
	}

	if req.Path == "" && len(req.ClassFiles) == 0 && req.ZipFile == "" {
		http.Error(w, "must provide path, class_files, or zipfile", http.StatusBadRequest)
		return
	}

	id := s.scanner.Submit(req)
	http.Redirect(w, r, "/scans/"+id, http.StatusSeeOther)
}

func (s *Server) handleGetScan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, ok := s.scanner.Get(id)
	if !ok {
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	}

	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	s.render(w, "scan.html", result)
}

func mustSub(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

type overlayFSType struct {
	primary   fs.FS
	secondary fs.FS
}

func overlayFS(primaryPath string, secondary fs.FS) fs.FS {
	return &overlayFSType{
		primary:   os.DirFS(primaryPath),
		secondary: secondary,
	}
}

func (o *overlayFSType) Open(name string) (fs.File, error) {
	f, err := o.primary.Open(name)
	if err == nil {
		return f, nil
	}
	return o.secondary.Open(name)
}

func (s *Scanner) AllClasses() []*java.Class {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []*java.Class
	for _, scan := range s.scans {
		if scan.Status == ScanStatusCompleted {
			all = append(all, scan.Classes...)
		}
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name() < all[j].Name()
	})
	return all
}

func (s *Scanner) FindClass(name string) *java.Class {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, scan := range s.scans {
		if scan.Status == ScanStatusCompleted {
			for _, c := range scan.Classes {
				if c.Name() == name {
					return c
				}
			}
		}
	}
	return nil
}

type ClassViewData struct {
	Classes      []*java.Class
	ActiveClass  *java.Class
	KnownClasses map[string]bool
	Implementers []*java.Class
}

func (s *Server) handleClass(w http.ResponseWriter, r *http.Request) {
	className := r.PathValue("className")
	classes := s.scanner.AllClasses()

	knownClasses := make(map[string]bool)
	for _, c := range classes {
		knownClasses[c.Name()] = true
	}

	data := ClassViewData{
		Classes:      classes,
		KnownClasses: knownClasses,
	}

	if className != "" {
		data.ActiveClass = s.scanner.FindClass(className)
		if data.ActiveClass == nil {
			http.Error(w, "class not found", http.StatusNotFound)
			return
		}
		if data.ActiveClass.IsInterface() {
			for _, c := range classes {
				for _, iface := range c.Interfaces() {
					if iface == className {
						data.Implementers = append(data.Implementers, c)
						break
					}
				}
			}
		}
	}

	s.render(w, "class.html", data)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		classes := s.scanner.AllClasses()
		if len(classes) > 0 {
			http.Redirect(w, r, "/c/", http.StatusSeeOther)
			return
		}
	}

	data := struct {
		Scans []*ScanResult
	}{
		Scans: s.scanner.List(),
	}
	s.render(w, "index.html", data)
}
