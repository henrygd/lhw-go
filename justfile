#!/usr/bin/env just --justfile

# build the final executable
build: build-dotnet build-go

# build the dotnet project to interact directly with lhw
build-dotnet:
	dotnet build -c Release

# build the final windows go executable 
build-go:
	GOOS=windows go build -ldflags "-w -s" -o ./get_temps.exe

# specific to my setup, copies to my windows machine
_build-and-copy: 
	@just build
	cp ./get_temps.exe /mnt/windows_share/

	