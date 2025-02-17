package lsp

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: make(map[string]string),
	}
}

type DidOpenTextDocumentNotification struct {
	Request
	Params DidOpenTextDocumentNotificationParams `json:"params"`
}

type DidOpenTextDocumentNotificationParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidChangeTextDocumentNotification struct {
	Request
	Params DidChangeTextDocumentNotificationParams `json:"params"`
}

type DidChangeTextDocumentNotificationParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Range *Range `json:"range"`
	Text  string `json:"text"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}
