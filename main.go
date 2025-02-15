package main

import (
	"bufio"
	"cdecl-lsp/lsp"
	"cdecl-lsp/rpc"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func main() {
	logger := getLogger("")
	logger.Println("started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error %s", err)
			continue
		}

		handleMessage(logger, method, content)
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
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
	}
	// TODO:
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
