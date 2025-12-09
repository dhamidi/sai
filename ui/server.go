package ui

import (
	"archive/zip"
	"bytes"
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
	"github.com/dhamidi/javalyzer/java/parser"
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
	Classes   []*java.ClassModel
	Error     string
	Errors    []string
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

	var classes []*java.ClassModel
	var errors []string

	if req.Path != "" {
		classes, errors = s.scanDirectory(req.ID, req.Path)
	} else if len(req.ClassFiles) > 0 {
		classes, errors = s.scanFiles(req.ID, req.ClassFiles)
	} else if req.ZipFile != "" {
		classes, errors = s.scanZipFile(req.ID, req.ZipFile)
	} else {
		errors = append(errors, "no path, class files, or zip file provided")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	result.EndedAt = time.Now()
	result.Classes = classes
	result.Errors = errors
	if len(errors) > 0 && len(classes) == 0 {
		result.Status = ScanStatusFailed
		result.Error = errors[0]
	} else {
		result.Status = ScanStatusCompleted
	}
}

func (s *Scanner) scanDirectory(id, path string) ([]*java.ClassModel, []string) {
	var files []string
	var errors []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, fmt.Sprintf("walk %s: %v", p, err))
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(p)
			if ext == ".class" || ext == ".java" {
				files = append(files, p)
			}
		}
		return nil
	})
	if err != nil {
		errors = append(errors, fmt.Sprintf("walk %s: %v", path, err))
	}
	classes, scanErrors := s.scanFiles(id, files)
	return classes, append(errors, scanErrors...)
}

func (s *Scanner) scanFiles(id string, files []string) ([]*java.ClassModel, []string) {
	s.mu.Lock()
	s.scans[id].Total = len(files)
	s.mu.Unlock()

	var classes []*java.ClassModel
	var errors []string
	for i, file := range files {
		ext := filepath.Ext(file)
		switch ext {
		case ".class":
			class, err := java.ClassModelFromFile(file)
			if err != nil {
				errors = append(errors, fmt.Sprintf("parse %s: %v", file, err))
			} else if class != nil {
				classes = append(classes, class)
			}
		case ".java":
			data, err := os.ReadFile(file)
			if err != nil {
				errors = append(errors, fmt.Sprintf("read %s: %v", file, err))
			} else {
				models, err := java.ClassModelsFromSource(data, parser.WithFile(filepath.Base(file)))
				if err != nil {
					errors = append(errors, fmt.Sprintf("parse %s: %v", file, err))
				} else {
					classes = append(classes, models...)
				}
			}
		}

		s.mu.Lock()
		s.scans[id].Progress = i + 1
		s.mu.Unlock()
	}
	return classes, errors
}

func (s *Scanner) scanZipFile(id, path string) ([]*java.ClassModel, []string) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, []string{fmt.Sprintf("open zip: %v", err)}
	}
	defer r.Close()

	var sourceFiles []*zip.File
	var jarFiles []*zip.File
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		switch ext {
		case ".class", ".java":
			sourceFiles = append(sourceFiles, f)
		case ".jar":
			jarFiles = append(jarFiles, f)
		}
	}

	total := len(sourceFiles)
	for _, jarFile := range jarFiles {
		total += s.countFilesInJar(jarFile)
	}

	s.mu.Lock()
	s.scans[id].Total = total
	s.mu.Unlock()

	var classes []*java.ClassModel
	var errors []string
	progress := 0

	for _, f := range sourceFiles {
		rc, err := f.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("open %s: %v", f.Name, err))
			progress++
			s.mu.Lock()
			s.scans[id].Progress = progress
			s.mu.Unlock()
			continue
		}

		ext := filepath.Ext(f.Name)
		switch ext {
		case ".class":
			class, err := java.ClassModelFromReader(rc)
			rc.Close()
			if err != nil {
				errors = append(errors, fmt.Sprintf("parse %s: %v", f.Name, err))
			} else if class != nil {
				classes = append(classes, class)
			}
		case ".java":
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				errors = append(errors, fmt.Sprintf("read %s: %v", f.Name, err))
			} else {
				models, err := java.ClassModelsFromSource(data, parser.WithFile(filepath.Base(f.Name)))
				if err != nil {
					errors = append(errors, fmt.Sprintf("parse %s: %v", f.Name, err))
				} else {
					classes = append(classes, models...)
				}
			}
		}

		progress++
		s.mu.Lock()
		s.scans[id].Progress = progress
		s.mu.Unlock()
	}

	for _, jarFile := range jarFiles {
		onProgress := func() {
			progress++
			s.mu.Lock()
			s.scans[id].Progress = progress
			s.mu.Unlock()
		}
		jarClasses, jarErrors := s.scanJarInZip(jarFile, onProgress)
		classes = append(classes, jarClasses...)
		errors = append(errors, jarErrors...)
	}

	return classes, errors
}

func (s *Scanner) countFilesInJar(jarFile *zip.File) int {
	rc, err := jarFile.Open()
	if err != nil {
		return 0
	}
	defer rc.Close()

	jarData, err := io.ReadAll(rc)
	if err != nil {
		return 0
	}

	jarReader, err := zip.NewReader(bytes.NewReader(jarData), int64(len(jarData)))
	if err != nil {
		return 0
	}

	count := 0
	for _, f := range jarReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		if ext == ".class" || ext == ".java" {
			count++
		}
	}
	return count
}

func (s *Scanner) scanJarInZip(jarFile *zip.File, onProgress func()) ([]*java.ClassModel, []string) {
	rc, err := jarFile.Open()
	if err != nil {
		return nil, []string{fmt.Sprintf("open jar %s: %v", jarFile.Name, err)}
	}
	defer rc.Close()

	jarData, err := io.ReadAll(rc)
	if err != nil {
		return nil, []string{fmt.Sprintf("read jar %s: %v", jarFile.Name, err)}
	}

	jarReader, err := zip.NewReader(bytes.NewReader(jarData), int64(len(jarData)))
	if err != nil {
		return nil, []string{fmt.Sprintf("open jar %s as zip: %v", jarFile.Name, err)}
	}

	var classes []*java.ClassModel
	var errors []string
	for _, f := range jarReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		if ext != ".class" && ext != ".java" {
			continue
		}

		fileRC, err := f.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("open %s in %s: %v", f.Name, jarFile.Name, err))
			onProgress()
			continue
		}

		switch ext {
		case ".class":
			fmt.Printf("[DEBUG] Parsing class %d: %s in %s\n", len(classes)+1, f.Name, jarFile.Name)
			class, err := java.ClassModelFromReader(fileRC)
			fileRC.Close()
			if err != nil {
				errors = append(errors, fmt.Sprintf("parse %s in %s: %v", f.Name, jarFile.Name, err))
			} else if class != nil {
				classes = append(classes, class)
			}
		case ".java":
			fmt.Printf("[DEBUG] Parsing java %d: %s in %s\n", len(classes)+1, f.Name, jarFile.Name)
			data, err := io.ReadAll(fileRC)
			fileRC.Close()
			if err != nil {
				errors = append(errors, fmt.Sprintf("read %s in %s: %v", f.Name, jarFile.Name, err))
			} else {
				models, err := java.ClassModelsFromSource(data, parser.WithFile(filepath.Base(f.Name)))
				if err != nil {
					errors = append(errors, fmt.Sprintf("parse %s in %s: %v", f.Name, jarFile.Name, err))
				} else {
					classes = append(classes, models...)
				}
			}
		}
		onProgress()
	}

	return classes, errors
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
		"linkifyType": func(knownClasses map[string]bool, t java.TypeModel) template.HTML {
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
		"constructors": func(m *java.ClassModel) []java.MethodModel {
			var ctors []java.MethodModel
			for _, method := range m.Methods {
				if method.Name == "<init>" {
					ctors = append(ctors, method)
				}
			}
			return ctors
		},
		"isConstructor": func(m java.MethodModel) bool {
			return m.Name == "<init>"
		},
		"isStaticInitializer": func(m java.MethodModel) bool {
			return m.Name == "<clinit>"
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

func (s *Scanner) AllClasses() []*java.ClassModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []*java.ClassModel
	for _, scan := range s.scans {
		if scan.Status == ScanStatusCompleted {
			all = append(all, scan.Classes...)
		}
	}
	java.ResolveInnerClassReferences(all)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})
	return all
}

func (s *Scanner) FindClass(name string) *java.ClassModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, scan := range s.scans {
		if scan.Status == ScanStatusCompleted {
			for _, c := range scan.Classes {
				if c.Name == name {
					return c
				}
			}
		}
	}
	return nil
}

type ClassViewData struct {
	Classes      []*java.ClassModel
	ActiveClass  *java.ClassModel
	KnownClasses map[string]bool
	Implementers []*java.ClassModel
}

func (s *Server) handleClass(w http.ResponseWriter, r *http.Request) {
	className := r.PathValue("className")
	classes := s.scanner.AllClasses()

	knownClasses := make(map[string]bool)
	for _, c := range classes {
		knownClasses[c.Name] = true
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
		if data.ActiveClass.Kind == java.ClassKindInterface {
			for _, c := range classes {
				for _, iface := range c.Interfaces {
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
