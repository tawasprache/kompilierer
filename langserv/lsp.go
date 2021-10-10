package langserv

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/fehlerberichtung"
	"Tawa/kompilierer/getypisiertast"
	"Tawa/kompilierer/typisierung"
	"context"
	"log"
	"path"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

type server struct {
	rootURI string
	files   map[string]string
}

func (s *server) Initialize(ctx context.Context, conn jsonrpc2.JSONRPC2, params lsp.InitializeParams) (*lsp.InitializeResult, *lsp.InitializeError) {
	s.rootURI = string(params.RootURI)
	s.files = map[string]string{}

	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Options: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TDSKFull,
				},
			},
		},
	}, nil
}

func zuP(von lexer.Position) lsp.Position {
	return lsp.Position{
		Line:      von.Line - 1,
		Character: von.Column - 1,
	}
}

func zuR(von getypisiertast.Span) lsp.Range {
	return lsp.Range{
		Start: zuP(von.Von),
		End:   zuP(von.Zu),
	}
}

func zuNurR(von lexer.Position) lsp.Range {
	v2 := von
	v2.Column++
	return lsp.Range{
		Start: zuP(von),
		End:   zuP(v2),
	}
}

func zuDiag(e error) []lsp.Diagnostic {
	switch feh := e.(type) {
	case participle.Error:
		return []lsp.Diagnostic{
			{
				Range:    zuNurR(feh.Position()),
				Severity: lsp.Error,
				Message:  feh.Message(),
			},
		}
	case fehlerberichtung.PositionError:
		return []lsp.Diagnostic{
			{
				Range:    zuR(feh.Span),
				Severity: lsp.Error,
				Message:  feh.Text,
			},
		}
	case fehlerberichtung.VerketteterFehler:
		fehler := []lsp.Diagnostic{}
		for _, it := range feh.Fehler {
			fehler = append(fehler, zuDiag(it)...)
		}
		return fehler
	default:
		panic("e " + repr.String(e))
	}
}

func (s *server) doDiag(ctx context.Context, conn jsonrpc2.JSONRPC2, uri lsp.DocumentURI, content string) {
	diags := lsp.PublishDiagnosticsParams{
		URI: uri,
	}

	dat := ast.Modul{}
	feh := ast.Parser.ParseString(path.Base(string(uri)), content, &dat)
	if feh != nil {
		diags.Diagnostics = append(diags.Diagnostics, zuDiag(feh)...)
		return
	}

	ktx := typisierung.NeuKontext()
	genannt, err := typisierung.Auflösenamen(ktx, dat, "User")
	if err != nil {
		log.Println("TODO")
		return
	}

	_, err = typisierung.Typiere(ktx, genannt, "User")
	if err != nil {
		diags.Diagnostics = append(diags.Diagnostics, zuDiag(err)...)
	}

	conn.Notify(ctx, "textDocument/publishDiagnostics", diags)
}

func (s *server) Initialized(ctx context.Context, conn jsonrpc2.JSONRPC2, params struct{}) {

}

func (s *server) DidOpen(ctx context.Context, conn jsonrpc2.JSONRPC2, params lsp.DidOpenTextDocumentParams) {
	s.files[strings.TrimPrefix(string(params.TextDocument.URI), s.rootURI)] = params.TextDocument.Text
	go s.doDiag(ctx, conn, params.TextDocument.URI, params.TextDocument.Text)
}

func (s *server) DidChange(ctx context.Context, conn jsonrpc2.JSONRPC2, params lsp.DidChangeTextDocumentParams) {
	s.files[strings.TrimPrefix(string(params.TextDocument.URI), s.rootURI)] = params.ContentChanges[0].Text
	go s.doDiag(ctx, conn, params.TextDocument.URI, params.ContentChanges[0].Text)
}

func (s *server) DidClose(ctx context.Context, conn jsonrpc2.JSONRPC2, params lsp.DidCloseTextDocumentParams) {
	delete(s.files, strings.TrimPrefix(string(params.TextDocument.URI), s.rootURI))
}
