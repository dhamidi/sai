package codebase

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dhamidi/javalyzer/java"
	"github.com/dhamidi/javalyzer/java/parser"
)

type Codebase struct {
	mu      sync.RWMutex
	rootDir string
	files   map[string]*FileInfo
	classes []*java.ClassModel
}

type FileInfo struct {
	Path     string
	Content  []byte
	AST      *parser.Node
	Classes  []*java.ClassModel
	ParseErr error
}

func New(rootDir string) *Codebase {
	return &Codebase{
		rootDir: rootDir,
		files:   make(map[string]*FileInfo),
	}
}

func (c *Codebase) RootDir() string {
	return c.rootDir
}

func (c *Codebase) ScanAll() error {
	return filepath.Walk(c.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".java" {
			c.ScanFile(path)
		}
		return nil
	})
}

func (c *Codebase) ScanFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return c.UpdateFile(path, content)
}

func (c *Codebase) UpdateFile(path string, content []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.updateFileLocked(path, content)
}

func (c *Codebase) updateFileLocked(path string, content []byte) error {
	p := parser.ParseCompilationUnit(bytes.NewReader(content), parser.WithFile(filepath.Base(path)), parser.WithPositions())
	ast := p.Finish()

	var classes []*java.ClassModel
	var parseErr error
	if ast != nil {
		classes, parseErr = java.ClassModelsFromSource(content, parser.WithFile(filepath.Base(path)))
	}

	c.files[path] = &FileInfo{
		Path:     path,
		Content:  content,
		AST:      ast,
		Classes:  classes,
		ParseErr: parseErr,
	}

	c.rebuildClassesLocked()
	return nil
}

func (c *Codebase) rebuildClassesLocked() {
	var all []*java.ClassModel
	for _, f := range c.files {
		if f.Classes != nil {
			all = append(all, f.Classes...)
		}
	}
	java.ResolveInnerClassReferences(all)
	c.classes = all
}

func (c *Codebase) RemoveFile(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.files, path)
	c.rebuildClassesLocked()
}

func (c *Codebase) GetFile(path string) *FileInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.files[path]
}

func (c *Codebase) AllClasses() []*java.ClassModel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.classes
}

func (c *Codebase) FindClass(name string) *java.ClassModel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, cls := range c.classes {
		if cls.Name == name {
			return cls
		}
	}
	return nil
}

func (c *Codebase) TypeAtPoint(path string, line, column int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	f := c.files[path]
	if f == nil || f.AST == nil {
		return ""
	}

	pos := parser.Position{Line: line, Column: column}
	return java.TypeAtPoint(f.AST, pos, c.classes)
}

func (c *Codebase) CompletionsAtPoint(path string, line, column int) []CompletionItem {
	typeName := c.TypeAtPoint(path, line, column)
	if typeName == "" {
		return nil
	}

	typeName = strings.TrimSuffix(typeName, "[]")

	cls := c.FindClass(typeName)
	if cls == nil {
		return nil
	}

	var items []CompletionItem

	for _, m := range cls.Methods {
		if m.Visibility != java.VisibilityPublic {
			continue
		}
		items = append(items, CompletionItem{
			Label:      m.Name,
			Kind:       CompletionKindMethod,
			Detail:     formatMethodSignature(m),
			InsertText: formatMethodInsert(m),
		})
	}

	for _, f := range cls.Fields {
		if f.Visibility != java.VisibilityPublic {
			continue
		}
		items = append(items, CompletionItem{
			Label:      f.Name,
			Kind:       CompletionKindField,
			Detail:     f.Type.Name,
			InsertText: f.Name,
		})
	}

	// Record components have implicit accessor methods
	for _, rc := range cls.RecordComponents {
		items = append(items, CompletionItem{
			Label:      rc.Name,
			Kind:       CompletionKindMethod,
			Detail:     rc.Type.Name,
			InsertText: rc.Name + "()",
		})
	}

	return items
}

type CompletionKind int

const (
	CompletionKindMethod CompletionKind = iota
	CompletionKindField
	CompletionKindClass
)

type CompletionItem struct {
	Label      string
	Kind       CompletionKind
	Detail     string
	InsertText string
}

func formatMethodSignature(m java.MethodModel) string {
	var params []string
	for _, p := range m.Parameters {
		params = append(params, p.Type.Name+" "+p.Name)
	}
	return m.ReturnType.Name + " " + m.Name + "(" + strings.Join(params, ", ") + ")"
}

func formatMethodInsert(m java.MethodModel) string {
	if len(m.Parameters) == 0 {
		return m.Name + "()"
	}
	var placeholders []string
	for i, p := range m.Parameters {
		name := p.Name
		if name == "" {
			name = p.Type.Name
		}
		placeholders = append(placeholders, "${"+itoa(i+1)+":"+name+"}")
	}
	return m.Name + "(" + strings.Join(placeholders, ", ") + ")"
}

func itoa(i int) string {
	return string(rune('0' + i))
}
