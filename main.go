package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed read-temps.ps1
var psScript []byte

//go:embed LibreHardwareMonitorLib.dll
var dllBytes []byte

type SensorReading struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func main() {
	// if !isElevated() {
	// 	return
	// }

	scriptPath := writeTempFile("read_temp.ps1", psScript)
	dllPath := writeTempFile("LibreHardwareMonitorLib.dll", dllBytes)
	defer os.Remove(scriptPath)
	defer os.Remove(dllPath)

	// Run PowerShell script with its working directory set to DLL location
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	cmd.Dir = filepath.Dir(dllPath)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Script error: %v\n%s", err, stderr.String())
	}

	var readings []SensorReading
	err = json.Unmarshal(out.Bytes(), &readings)
	if err != nil {
		log.Fatalf("JSON parse error: %v\nOutput: %s", err, out.String())
	}

	for _, r := range readings {
		fmt.Printf("%s: %.2f °C\n", r.Name, r.Value)
	}
}

func writeTempFile(name string, data []byte) string {
	path := filepath.Join(os.TempDir(), name)
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		log.Fatalf("Failed to write temp file %s: %v", name, err)
	}
	return path
}
