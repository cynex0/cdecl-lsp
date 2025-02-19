package main

import (
	"bufio"
	"cdecl-lsp/lsp"
	"cdecl-lsp/parser"
	"cdecl-lsp/rpc"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	logger := getLogger("")
	logger.Println("started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := lsp.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error %s", err)
			continue
		}

		handleMessage(&state, logger, method, content)
	}
}

func handleMessage(state *lsp.State, logger *log.Logger, method string, content []byte) {
	logger.Printf("Received message: %s", content)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("could not parse message: %s", err)
		}

		res := rpc.EncodeMessage(lsp.NewInitializeResponse(request.ID))

		writer := os.Stdout
		writer.Write([]byte(res))

		logger.Printf("Initialized: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)
		logger.Printf("Sent response: %s", res)

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("could not parse message: %s", err)
		}

		state.Documents[request.Params.TextDocument.URI] = request.Params.TextDocument.Text

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("could not parse message: %s", err)
		}
		_, contains := state.Documents[request.Params.TextDocument.URI]
		if contains {
			changes := request.Params.ContentChanges
			for _, change := range changes {
				if change.Range != nil { // if the change capability is set to 2
					// TODO: implement range changes
					logger.Printf("Do not know what to do with a range change: %s", change.Text)
					break
				}

				state.Documents[request.Params.TextDocument.URI] = change.Text
			}
		} else {
			logger.Printf("Could not apply changes to a non-opened document")
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("could not parse message: %s", err)
		}

		lines := strings.Split(state.Documents[request.Params.TextDocument.URI], "\n")
		line := lines[request.Params.Position.Line]
		if parser.IsDeclaration(line) {
			res := rpc.EncodeMessage(lsp.NewHoverResponse(request.ID, "You are a type declaration, Harry"))

			writer := os.Stdout

			writer.Write([]byte(res))
			logger.Printf("sent response: %s", res)
		}

	}
}

func getLogger(filename string) *log.Logger {
	if filename == "" {
		filename = defaultLogPath()
	}

	err := os.MkdirAll(filepath.Dir(filename), 0o755)
	if err != nil {
		panic("could not create log directory " + err.Error())
	}

	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		panic("could not open logger file " + filename + err.Error())
	}

	return log.New(logfile, "[cdecl-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}

func defaultLogPath() string {
	var baseDir string

	if xdgStateHome := os.Getenv("XDG_STATE_HOME"); xdgStateHome != "" {
		baseDir = xdgStateHome
	} else if home, err := os.UserHomeDir(); err == nil {
		baseDir = filepath.Join(home, ".local", "state")
	} else {
		baseDir = "/tmp" // Fallback for unusual cases
	}

	if os.Getenv("OS") == "Windows_NT" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			baseDir = appData
		} else {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	}

	return filepath.Join(baseDir, "cdecl-lsp", "cdecl-lsp.log")
}
