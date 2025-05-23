using System;
using System.Collections.Generic;
using Newtonsoft.Json;
using LibreHardwareMonitor.Hardware;

class SensorReading
{
  public string Name { get; set; }
  public float Value { get; set; }
}

class Program
{
  static void Main()
  {
    var computer = new Computer
    {
      IsCpuEnabled = true,
      IsGpuEnabled = true,
      IsMemoryEnabled = true,
      IsMotherboardEnabled = true,
      IsStorageEnabled = true,
      IsPsuEnabled = true,
      IsNetworkEnabled = true,
    };
    computer.Open();

    var reader = Console.In;
    var writer = Console.Out;

    string line;
    while ((line = reader.ReadLine()) != null)
    {
      if (line.Trim().Equals("getTemps", StringComparison.OrdinalIgnoreCase))
      {
        var result = new List<SensorReading>();
        foreach (var hw in computer.Hardware)
        {
          var updated = false;
          foreach (var sensor in hw.Sensors)
          {
            if (sensor.SensorType == SensorType.Temperature && sensor.Value.HasValue)
            {
              if (!updated)
              {
                hw.Update();
                updated = true;
              }
              result.Add(new SensorReading
              {
                Name = sensor.Name,
                Value = sensor.Value.Value
              });
            }
          }
        }

        var json = JsonConvert.SerializeObject(result);
        writer.WriteLine(json);
        writer.Flush();
      }
    }

    computer.Close();
  }
}
