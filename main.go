package main

import (
	"bufio"
	"cdecl-lsp/rpc"
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
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
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
