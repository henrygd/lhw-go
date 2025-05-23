Example cgo-free single binary that embeds [LibreHardwareMonitorLib](https://github.com/LibreHardwareMonitor/LibreHardwareMonitor) to read temperature sensors on Windows. It uses a C# wrapper to interact with LHM / .NET.

Build the executable with `just`. This should work on any OS with Go and `dotnet` installed.

```
just
```

Unfortunately it does need to be run as administrator to access the data, like Libre Hardware Monitor itself.

```powershell
.\get_temps.exe
```
