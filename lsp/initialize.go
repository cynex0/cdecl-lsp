package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeRequestParams struct {
	// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initialize
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initializeResult
	ServerInfo   *ClientInfo        `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	HoverProvider         bool                  `json:"hoverProvider"`
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider"`
}

type SignatureHelpOptions struct {
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			ServerInfo: &ClientInfo{
				Name:    "cdecl-lsp",
				Version: "0.0.1",
			},
			Capabilities: ServerCapabilities{},
		},
	}
}
