package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

type SensorReading struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

//go:embed all:bin/Release/net48
var libreReaderFs embed.FS

func main() {
	tempDir := os.TempDir()
	destDir := filepath.Join(tempDir, "get_temps")

	err := copyEmbeddedDir(libreReaderFs, "bin/Release/net48", destDir)
	if err != nil {
		log.Fatalf("Failed to copy embedded directory: %v", err)
	}

	cmd := exec.Command(filepath.Join(destDir, "get_temps.exe"))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start get_temps.exe: %v", err)
	}

	scanner := bufio.NewScanner(stdout)

	for {
		_, err = fmt.Fprintln(stdin, "getTemps")
		if err != nil {
			log.Fatalf("Failed to send command: %v", err)
		}

		if scanner.Scan() {
			line := scanner.Text()
			var readings []SensorReading
			err := json.Unmarshal([]byte(line), &readings)
			if err != nil {
				log.Printf("Failed to parse JSON: %v\nGot: %s", err, line)
				continue
			}

			log.Println(readings)

		} else if err := scanner.Err(); err != nil {
			log.Printf("Failed to read output: %v", err)
			return
		} else {
			log.Printf("No output from get_temps.exe")
			return
		}
		time.Sleep(time.Second)
	}
}

func copyEmbeddedDir(fs embed.FS, srcPath, destPath string) error {
	entries, err := fs.ReadDir(srcPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcEntryPath := path.Join(srcPath, entry.Name())
		destEntryPath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			if err := copyEmbeddedDir(fs, srcEntryPath, destEntryPath); err != nil {
				return err
			}
			continue
		}

		data, err := fs.ReadFile(srcEntryPath)
		if err != nil {
			return err
		}

		if err := os.WriteFile(destEntryPath, data, 0755); err != nil {
			return err
		}
	}

	return nil
}
