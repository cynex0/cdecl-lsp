package lsp

import (
	"cdecl-lsp/parser"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func HandleInitialize(content []byte) (*InitializeResponse, error) {
	var request InitializeRequest
	if err := json.Unmarshal(content, &request); err != nil {
		return nil, err
	}

	res := NewInitializeResponse(request.ID)

	return &res, nil
}

func HandleDidOpen(content []byte, state *State) error {
	var request DidOpenTextDocumentNotification
	if err := json.Unmarshal(content, &request); err != nil {
		return err
	}

	state.Documents[request.Params.TextDocument.URI] = request.Params.TextDocument.Text
	return nil
}

func HandleDidChange(content []byte, state *State) error {
	var request DidChangeTextDocumentNotification
	if err := json.Unmarshal(content, &request); err != nil {
		return err
	}
	_, contains := state.Documents[request.Params.TextDocument.URI]
	if contains {
		changes := request.Params.ContentChanges
		for _, change := range changes {
			if change.Range != nil { // if the change capability is set to 2
				// TODO: implement range changes
				return fmt.Errorf("Do not know what to do with a range change: %s", change.Text)
			}

			state.Documents[request.Params.TextDocument.URI] = change.Text
		}
	} else {
		return errors.New("Could not apply changes to a non-opened document")
	}

	return nil
}

func HandleHover(content []byte, state *State) (*HoverResponse, error) {
	var request HoverRequest
	if err := json.Unmarshal(content, &request); err != nil {
		return nil, err
	}

	lines := strings.Split(state.Documents[request.Params.TextDocument.URI], "\n")
	line := lines[request.Params.Position.Line]

	res := NewHoverResponse(request.ID, "")
	if parser.IsDeclaration(line) {
		explain, err := parser.Explain(line)
		if err != nil {
			return &res, err
		}

		res.Result.Contents = explain
	}

	return &res, nil
}
