package codebase

import (
	"archive/zip"
	"bytes"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhamidi/javalyzer/java"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const lsName = "javalyzer"

type LSPServer struct {
	codebase *Codebase
	handler  protocol.Handler
	server   *server.Server
	version  string
}

func NewLSPServer(version string) *LSPServer {
	ls := &LSPServer{
		version: version,
	}

	ls.handler = protocol.Handler{
		Initialize:             ls.initialize,
		Initialized:            ls.initialized,
		Shutdown:               ls.shutdown,
		SetTrace:               ls.setTrace,
		TextDocumentDidOpen:    ls.textDocumentDidOpen,
		TextDocumentDidChange:  ls.textDocumentDidChange,
		TextDocumentDidClose:   ls.textDocumentDidClose,
		TextDocumentDidSave:    ls.textDocumentDidSave,
		TextDocumentCompletion: ls.textDocumentCompletion,
	}

	ls.server = server.NewServer(&ls.handler, lsName, false)

	return ls
}

func (ls *LSPServer) RunStdio() error {
	return ls.server.RunStdio()
}

func (ls *LSPServer) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	rootDir := "."
	if params.RootPath != nil && *params.RootPath != "" {
		rootDir = *params.RootPath
	} else if params.RootURI != nil && *params.RootURI != "" {
		if path, err := uriToPath(*params.RootURI); err == nil {
			rootDir = path
		}
	}

	ls.codebase = New(rootDir)

	capabilities := ls.handler.CreateServerCapabilities()

	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: boolPtr(true),
		Change:    intPtr(int(protocol.TextDocumentSyncKindFull)),
		Save: &protocol.SaveOptions{
			IncludeText: boolPtr(true),
		},
	}

	triggerChars := []string{"."}
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: triggerChars,
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &ls.version,
		},
	}, nil
}

func (ls *LSPServer) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	ls.codebase.ScanAll()
	if javaSrc := os.Getenv("JAVA_SRC"); javaSrc != "" {
		ls.scanJavaSrcZip(javaSrc)
	}
	return nil
}

func (ls *LSPServer) scanJavaSrcZip(zipPath string) {
	ext := filepath.Ext(zipPath)
	switch ext {
	case ".zip", ".jar":
		ls.scanZipOrJar(zipPath)
	case ".class":
		ls.scanClassFile(zipPath)
	case ".java":
		ls.codebase.ScanFile(zipPath)
	}
}

func (ls *LSPServer) scanZipOrJar(zipPath string) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return
	}
	defer r.Close()

	var jarFiles []*zip.File
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		switch ext {
		case ".java":
			ls.scanZipEntryJava(f, zipPath)
		case ".class":
			ls.scanZipEntryClass(f, zipPath)
		case ".jar":
			jarFiles = append(jarFiles, f)
		}
	}

	for _, jarFile := range jarFiles {
		ls.scanJarInZip(jarFile, zipPath)
	}
}

func (ls *LSPServer) scanZipEntryJava(f *zip.File, zipPath string) {
	rc, err := f.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return
	}

	virtualPath := "jdk:" + f.Name
	ls.codebase.UpdateFile(virtualPath, content)
}

func (ls *LSPServer) scanZipEntryClass(f *zip.File, zipPath string) {
	rc, err := f.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return
	}

	model, err := java.ClassModelFromReader(bytes.NewReader(data))
	if err != nil || model == nil {
		return
	}

	ls.codebase.AddClassModel(model)
}

func (ls *LSPServer) scanJarInZip(jarFile *zip.File, zipPath string) {
	rc, err := jarFile.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	jarData, err := io.ReadAll(rc)
	if err != nil {
		return
	}

	jarReader, err := zip.NewReader(bytes.NewReader(jarData), int64(len(jarData)))
	if err != nil {
		return
	}

	for _, f := range jarReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		switch ext {
		case ".java":
			ls.scanNestedJarEntryJava(f, jarFile.Name)
		case ".class":
			ls.scanNestedJarEntryClass(f, jarFile.Name)
		}
	}
}

func (ls *LSPServer) scanNestedJarEntryJava(f *zip.File, jarName string) {
	rc, err := f.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return
	}

	virtualPath := "jdk:" + jarName + "!" + f.Name
	ls.codebase.UpdateFile(virtualPath, content)
}

func (ls *LSPServer) scanNestedJarEntryClass(f *zip.File, jarName string) {
	rc, err := f.Open()
	if err != nil {
		return
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return
	}

	model, err := java.ClassModelFromReader(bytes.NewReader(data))
	if err != nil || model == nil {
		return
	}

	ls.codebase.AddClassModel(model)
}

func (ls *LSPServer) scanClassFile(path string) {
	model, err := java.ClassModelFromFile(path)
	if err != nil || model == nil {
		return
	}
	ls.codebase.AddClassModel(model)
}

func (ls *LSPServer) shutdown(ctx *glsp.Context) error {
	return nil
}

func (ls *LSPServer) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (ls *LSPServer) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	path, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return nil
	}
	ls.codebase.UpdateFile(path, []byte(params.TextDocument.Text))
	return nil
}

func (ls *LSPServer) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	path, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return nil
	}
	if len(params.ContentChanges) > 0 {
		change := params.ContentChanges[len(params.ContentChanges)-1]
		if textChange, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			ls.codebase.UpdateFile(path, []byte(textChange.Text))
		}
	}
	return nil
}

func (ls *LSPServer) textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	return nil
}

func (ls *LSPServer) textDocumentDidSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	path, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return nil
	}
	if params.Text != nil {
		ls.codebase.UpdateFile(path, []byte(*params.Text))
	} else {
		ls.codebase.ScanFile(path)
	}
	return nil
}

func (ls *LSPServer) textDocumentCompletion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	path, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return nil, nil
	}

	line := int(params.Position.Line) + 1
	col := int(params.Position.Character)

	file := ls.codebase.GetFile(path)
	if file == nil {
		return nil, nil
	}

	triggerCol := findTriggerPosition(file.Content, line, col)
	if triggerCol < 0 {
		return nil, nil
	}

	completions := ls.codebase.CompletionsAtPoint(path, line, triggerCol)
	if len(completions) == 0 {
		return nil, nil
	}

	var items []protocol.CompletionItem
	for _, c := range completions {
		kind := toProtocolKind(c.Kind)
		detail := c.Detail
		insertText := c.InsertText
		format := protocol.InsertTextFormatSnippet

		items = append(items, protocol.CompletionItem{
			Label:            c.Label,
			Kind:             &kind,
			Detail:           &detail,
			InsertText:       &insertText,
			InsertTextFormat: &format,
		})
	}

	return items, nil
}

func findTriggerPosition(content []byte, line, col int) int {
	lines := strings.Split(string(content), "\n")
	if line <= 0 || line > len(lines) {
		return -1
	}
	lineContent := lines[line-1]

	for i := col - 1; i >= 0; i-- {
		if i < len(lineContent) && lineContent[i] == '.' {
			return i
		}
	}
	return -1
}

func toProtocolKind(kind CompletionKind) protocol.CompletionItemKind {
	switch kind {
	case CompletionKindMethod:
		return protocol.CompletionItemKindMethod
	case CompletionKindField:
		return protocol.CompletionItemKindField
	case CompletionKindClass:
		return protocol.CompletionItemKindClass
	default:
		return protocol.CompletionItemKindText
	}
}

func uriToPath(uri string) (string, error) {
	if strings.HasPrefix(uri, "file://") {
		parsed, err := url.Parse(uri)
		if err != nil {
			return "", err
		}
		return filepath.Clean(parsed.Path), nil
	}
	return uri, nil
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *protocol.TextDocumentSyncKind {
	v := protocol.TextDocumentSyncKind(i)
	return &v
}

func getRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
