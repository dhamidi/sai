package scanner

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/parser"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Request struct {
	ID         string
	Path       string
	ClassFiles []string
	ZipFile    string
	CreatedAt  time.Time
}

type Result struct {
	ID        string
	Status    Status
	Request   Request
	Classes   []*java.ClassModel
	Error     string
	Errors    []string
	StartedAt time.Time
	EndedAt   time.Time
	Progress  int
	Total     int
}

func (s *Result) ProgressPercent() int {
	if s.Total == 0 {
		return 0
	}
	return (s.Progress * 100) / s.Total
}

type Scanner struct {
	mu       sync.RWMutex
	scans    map[string]*Result
	requests chan Request
	nextID   int
}

func New() *Scanner {
	s := &Scanner{
		scans:    make(map[string]*Result),
		requests: make(chan Request, 100),
	}
	go s.run()
	return s
}

func (s *Scanner) run() {
	for req := range s.requests {
		s.processScan(req)
	}
}

func (s *Scanner) processScan(req Request) {
	s.mu.Lock()
	result := s.scans[req.ID]
	result.Status = StatusInProgress
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
		result.Status = StatusFailed
		result.Error = errors[0]
	} else {
		result.Status = StatusCompleted
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
				models, err := java.ClassModelsFromSource(data, parser.WithFile(filepath.Base(file)), parser.WithSourcePath(file))
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
				models, err := java.ClassModelsFromSource(data, parser.WithFile(filepath.Base(f.Name)), parser.WithSourcePath(f.Name))
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
				models, err := java.ClassModelsFromSource(data, parser.WithFile(filepath.Base(f.Name)), parser.WithSourcePath(f.Name))
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

func (s *Scanner) Submit(req Request) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	req.ID = fmt.Sprintf("%d", s.nextID)
	req.CreatedAt = time.Now()

	s.scans[req.ID] = &Result{
		ID:      req.ID,
		Status:  StatusPending,
		Request: req,
	}

	s.requests <- req
	return req.ID
}

func (s *Scanner) Get(id string) (*Result, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.scans[id]
	return result, ok
}

func (s *Scanner) List() []*Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := make([]*Result, 0, len(s.scans))
	for _, r := range s.scans {
		results = append(results, r)
	}
	return results
}

func (s *Scanner) AllClasses() []*java.ClassModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []*java.ClassModel
	for _, scan := range s.scans {
		if scan.Status == StatusCompleted {
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
		if scan.Status == StatusCompleted {
			for _, c := range scan.Classes {
				if c.Name == name {
					return c
				}
			}
		}
	}
	return nil
}
