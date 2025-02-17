package lsp

type HoverRequest struct {
	Request
	Params TextDocumentPositionParams `json:"params"`
}

type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}

func NewHoverResponse(id int, contents string) HoverResponse {
	return HoverResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: HoverResult{Contents: contents},
	}
}
