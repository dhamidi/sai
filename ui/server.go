package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/scanner"
)

//go:embed static templates
var embeddedFS embed.FS

type Server struct {
	scanner    *scanner.Scanner
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
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"limit": func(n int, classes []*java.ClassModel) []*java.ClassModel {
			if n >= len(classes) {
				return classes
			}
			return classes[:n]
		},
		"formatJavadoc": func(javadoc string) template.HTML {
			if javadoc == "" {
				return ""
			}
			lines := strings.Split(javadoc, "\n")
			var result []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "/**" || line == "*/" {
					continue
				}
				if strings.HasPrefix(line, "* ") {
					line = line[2:]
				} else if line == "*" {
					line = ""
				}
				result = append(result, template.HTMLEscapeString(line))
			}
			return template.HTML(strings.Join(result, "<br>"))
		},
		"hasJavadoc": func(javadoc string) bool {
			return javadoc != ""
		},
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
		scanner:    scanner.New(),
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
	s.mux.HandleFunc("GET /sidebar", s.handleSidebar)
	s.mux.HandleFunc("GET /", s.handleIndex)

	if javaSrc := os.Getenv("JAVA_SRC"); javaSrc != "" {
		s.scanner.Submit(scanner.Request{ZipFile: javaSrc})
	}

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
	var req scanner.Request

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
			tmpFile, err := os.CreateTemp("", "sai-*.zip")
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

func (o *overlayFSType) ReadDir(name string) ([]fs.DirEntry, error) {
	entries := make(map[string]fs.DirEntry)

	if rd, ok := o.secondary.(fs.ReadDirFS); ok {
		if list, err := rd.ReadDir(name); err == nil {
			for _, e := range list {
				entries[e.Name()] = e
			}
		}
	}

	if rd, ok := o.primary.(fs.ReadDirFS); ok {
		if list, err := rd.ReadDir(name); err == nil {
			for _, e := range list {
				entries[e.Name()] = e
			}
		}
	}

	result := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		result = append(result, e)
	}
	return result, nil
}

type ClassViewData struct {
	Classes      []*java.ClassModel
	ActiveClass  *java.ClassModel
	KnownClasses map[string]bool
	Implementers []*java.ClassModel
	TotalMatches int
	HasMore      bool
}

func (s *Server) handleClass(w http.ResponseWriter, r *http.Request) {
	className := r.PathValue("className")
	allClasses := s.scanner.AllClasses()

	knownClasses := make(map[string]bool)
	for _, c := range allClasses {
		knownClasses[c.Name] = true
	}

	const maxResults = 20
	classes := allClasses
	if len(allClasses) > maxResults {
		classes = allClasses[:maxResults]
	}

	data := ClassViewData{
		Classes:      classes,
		TotalMatches: len(allClasses),
		HasMore:      len(allClasses) > maxResults,
		KnownClasses: knownClasses,
	}

	if className != "" {
		data.ActiveClass = s.scanner.FindClass(className)
		if data.ActiveClass == nil {
			http.Error(w, "class not found", http.StatusNotFound)
			return
		}
		if data.ActiveClass.Kind == java.ClassKindInterface {
			for _, c := range allClasses {
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
		Scans []*scanner.Result
	}{
		Scans: s.scanner.List(),
	}
	s.render(w, "index.html", data)
}

func (s *Server) handleSidebar(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	activeClassName := r.URL.Query().Get("active")

	const maxResults = 20

	allClasses := s.scanner.AllClasses()

	var classes []*java.ClassModel
	var totalMatches int
	if query == "" {
		totalMatches = len(allClasses)
		if len(allClasses) > maxResults {
			classes = allClasses[:maxResults]
		} else {
			classes = allClasses
		}
	} else {
		for _, c := range allClasses {
			if strings.Contains(strings.ToLower(c.Name), query) ||
				strings.Contains(strings.ToLower(c.SimpleName), query) {
				totalMatches++
				if len(classes) < maxResults {
					classes = append(classes, c)
				}
			}
		}
	}

	data := struct {
		Classes         []*java.ClassModel
		ActiveClassName string
		TotalMatches    int
		HasMore         bool
	}{
		Classes:         classes,
		ActiveClassName: activeClassName,
		TotalMatches:    totalMatches,
		HasMore:         totalMatches > maxResults,
	}
	s.render(w, "_sidebar.html", data)
}
