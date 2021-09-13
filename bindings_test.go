package nvml

import (
    "github.com/stretchr/testify/require"
    "testing"
)

func TestGetDeviceCount(t *testing.T) {
    require.NoError(t, Init())
    defer Shutdown()

    deviceCount, err := GetDeviceCount()
    require.NoError(t, err)

    t.Logf("GetDeviceCount: %d", deviceCount)
}

func TestGetDriverVersion(t *testing.T) {
    require.NoError(t, Init())
    defer Shutdown()

    version, err := GetDriverVersion()
    require.NoError(t, err)

    t.Logf("GetDriverVersion: %s", version)
}

func TestDevice(t *testing.T) {
    require.NoError(t, Init())
    defer Shutdown()

    device, err := GetDeviceByIndex(0)
    require.NoError(t, err)

    name, err := device.Name()
    require.NoError(t, err)
    t.Logf("name: %s", name)

    mname, err := device.MinorName()
    require.NoError(t, err)
    t.Logf("minor_name: %d", mname)

    uuid, err := device.UUID()
    require.NoError(t, err)
    t.Logf("uuid: %s", uuid)

    mem, err := device.Memory()
    require.NoError(t, err)
    t.Logf("memory: %+v", mem)

    bar1, err := device.Bar1Memory()
    require.NoError(t, err)
    t.Logf("bar1_memory: %+v", bar1)

    clock, err := device.Clock()
    require.NoError(t, err)
    t.Logf("clock: %+v", clock)

    temperature, err := device.Temperature()
    require.NoError(t, err)
    t.Logf("temperature: %+v", temperature)

    threshold, err := device.TemperatureThreshold()
    require.NoError(t, err)
    t.Logf("temperature_threshold: %+v", threshold)

    fanSpeed, err := device.FanSpeed()
    require.NoError(t, err)
    t.Logf("fan_speed: %d %%", fanSpeed)

    utilization, err := device.UtilizationRates()
    require.NoError(t, err)
    t.Logf("utilization: %+v", utilization)

    pci, err := device.PciInfo()
    require.NoError(t, err)
    t.Logf("pci: %+v", pci)

    running, err := device.RunningProcess()
    require.NoError(t, err)
    t.Logf("running_processes: %+v", running)

    power, err := device.PowerUsage()
    require.NoError(t, err)
    t.Logf("power_usage: %d", power)

    computeMode, err := device.ComputeMode()
    require.NoError(t, err)
    t.Logf("compute_mode: %s", computeMode)
}
