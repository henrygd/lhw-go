$dllPath = Join-Path $PSScriptRoot "LibreHardwareMonitorLib.dll"

Add-Type -Path $dllPath

$computer = New-Object LibreHardwareMonitor.Hardware.Computer

$computer.IsCpuEnabled = $true
$computer.IsGpuEnabled = $true
$computer.IsStorageEnabled = $true
$computer.IsMemoryEnabled = $true
$computer.IsMotherboardEnabled = $true
$computer.IsPsuEnabled = $true

$computer.Open()

$readings = @()

foreach ($hw in $computer.Hardware) {
    $hwUpdated = $false
    foreach ($sensor in $hw.Sensors) {
        if ($sensor.SensorType -eq [LibreHardwareMonitor.Hardware.SensorType]::Temperature) {
            if (-not $hwUpdated) {
                $hw.Update()
                $hwUpdated = $true
            }
            $readings += [PSCustomObject]@{
                name  = $sensor.Name
                value = $sensor.Value
            }
        }
    }
}

$computer.Close()

# Output as JSON
if ($readings.Count -eq 1) {
    # For single element, manually wrap in array brackets
    $json = $readings | ConvertTo-Json -Depth 2 -Compress
    Write-Output "[$json]"
} else {
    $readings | ConvertTo-Json -Depth 2 -Compress
}

