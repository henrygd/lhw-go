package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SensorReading struct {
	Name  string
	Value float64
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

		var readings []SensorReading

		// Read all sensor lines until we hit an empty line or EOF
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				break // Empty line indicates end of sensor data
			}

			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				log.Printf("Invalid sensor line format: %s", line)
				continue
			}

			name := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])

			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				log.Printf("Failed to parse temperature value '%s': %v", valueStr, err)
				continue
			}

			readings = append(readings, SensorReading{
				Name:  name,
				Value: value,
			})
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Failed to read output: %v", err)
			return
		}

		if len(readings) > 0 {
			log.Println(readings)
		} else {
			log.Printf("No sensor readings received")
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
