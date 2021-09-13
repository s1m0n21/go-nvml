package main

import (
	"fmt"
	"github.com/s1m0n21/go-nvml"
)

func main() {
	err := nvml.Init()
	if err != nil {
		panic(err)
	}
	defer nvml.Shutdown()

	deviceCount, err := nvml.GetDeviceCount()
	handleErr(err)

	driverVersion, err := nvml.GetDriverVersion()
	handleErr(err)

	fmt.Printf("DriverVersion: %s\n", driverVersion)
	fmt.Printf("DeviceCount:(%d)\n", deviceCount)

	var i uint = 0
	for ; i < deviceCount; i++ {
		device, err := nvml.GetDeviceByIndex(i)
		handleErr(err)

		deviceName, err := device.Name()
		handleErr(err)

		pciInfo, err := device.PciInfo()
		handleErr(err)

		uuid, err := device.UUID()
		handleErr(err)

		computeMode, err := device.ComputeMode()
		handleErr(err)

		utilization, err := device.UtilizationRates()
		handleErr(err)

		memoryInfo, err := device.Memory()
		handleErr(err)

		bar1Memory, err := device.Bar1Memory()
		handleErr(err)

		clockInfo, err := device.Clock()
		handleErr(err)

		temperature, err := device.Temperature()
		handleErr(err)

		temperatureThreshold, err := device.TemperatureThreshold()
		handleErr(err)

		fanSpeed, err := device.FanSpeed()
		handleErr(err)

		powerUsage, err := device.PowerUsage()
		handleErr(err)

		runningProcess, err := device.RunningProcess()
		handleErr(err)

		fmt.Printf("\t%s [%s]\n", deviceName, pciInfo.BusId)
		fmt.Printf("\tUUID: %s\n", uuid)
		fmt.Printf("\tComputeMode: %s\n", computeMode)
		fmt.Printf("\tGPUUtilization: %d %%\n", utilization.GPU)
		fmt.Printf("\tMemoryUtilization: %d %%\n", utilization.Memory)
		fmt.Printf("\tMemory.Total: %d bytes\n", memoryInfo.Total)
		fmt.Printf("\tMemory.Free: %d bytes\n", memoryInfo.Free)
		fmt.Printf("\tMemory.Used: %d bytes\n", memoryInfo.Used)
		fmt.Printf("\tBar1Memory.Total: %d bytes\n", bar1Memory.Total)
		fmt.Printf("\tBar1Memory.Free: %d bytes\n", bar1Memory.Free)
		fmt.Printf("\tBar1Memory.Used: %d bytes\n", bar1Memory.Used)
		fmt.Printf("\tMemClock: %d MHz (max: %d MHz)\n", clockInfo.Mem, clockInfo.MemMax)
		fmt.Printf("\tSMClock: %d MHz (max: %d MHz)\n", clockInfo.SM, clockInfo.SMMax)
		fmt.Printf("\tGraphicsClock: %d MHz (max: %d MHz)\n", clockInfo.Graphics, clockInfo.GraphicsMax)
		fmt.Printf("\tTemperature: %d C (slowdown: %d C, shutdown: %d C)\n",
			temperature, temperatureThreshold.Slowdown, temperatureThreshold.Shutdown)
		fmt.Printf("\tFanSpeed: %d %%\n", fanSpeed)
		fmt.Printf("\tPowerUsage: %d mW\n", powerUsage)
		fmt.Printf("\tRunningProcess:(%d)\n", len(runningProcess))
		for i, p := range runningProcess {
			fmt.Printf("\t\t%d. pid:%d name:%s memoryUsed:%d\n", i, p.PID, p.Name, p.MemoryUsed)
		}
		fmt.Println()
	}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
