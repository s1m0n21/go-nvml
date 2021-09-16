package nvml

//#cgo CFLAGS: -I /usr/src/gdk/nvml/include -I /usr/include/nvidia/gdk
//#cgo LDFLAGS: -lnvidia-ml -L /usr/src/gdk/nvml/lib
//#include <nvml.h>
import "C"
import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Device struct {
	device C.nvmlDevice_t
}

type PciInfo struct {
	BusId          string
	Domain         uint
	Bus            uint
	Device         uint
	PciDeviceId    uint
	PciSubSystemId uint

	reserved0 uint
	reserved1 uint
	reserved2 uint
	reserved3 uint
}

type TemperatureThreshold struct {
	Slowdown uint
	Shutdown uint
}

type UtilizationRate struct {
	GPU    uint
	Memory uint
}

type Memory struct {
	Total uint64
	Free uint64
	Used uint64
}

type Clock struct {
	Graphics    uint
	GraphicsMax uint
	SM          uint
	SMMax       uint
	Mem         uint
	MemMax      uint
}

type Process struct {
	PID        uint
	Name       string
	MemoryUsed uint64
}

func (d *Device) Name() (string, error) {
	var name [C.NVML_DEVICE_NAME_BUFFER_SIZE]C.char
	result := C.nvmlDeviceGetName(d.device, &name[0], C.NVML_DEVICE_NAME_BUFFER_SIZE)
	return C.GoString(&name[0]), err(result)
}

func (d *Device) MinorName() (uint, error) {
	var name C.uint
	result := C.nvmlDeviceGetMinorNumber(d.device, &name)
	return uint(name), err(result)
}

func (d *Device) UUID() (string, error) {
	var uuid [C.NVML_DEVICE_UUID_BUFFER_SIZE]C.char
	result := C.nvmlDeviceGetUUID(d.device, &uuid[0], C.NVML_DEVICE_UUID_BUFFER_SIZE)
	return C.GoString(&uuid[0]), err(result)
}

func (d *Device) PciInfo() (PciInfo, error) {
	var info C.nvmlPciInfo_t
	result := C.nvmlDeviceGetPciInfo_v2(d.device, &info)
	return PciInfo{
		BusId:          C.GoString(&info.busId[0]),
		Domain:         uint(info.domain),
		Bus:            uint(info.bus),
		Device:         uint(info.device),
		PciDeviceId:    uint(info.pciDeviceId),
		PciSubSystemId: uint(info.pciSubSystemId),
		reserved0:      uint(info.reserved0),
		reserved1:      uint(info.reserved1),
		reserved2:      uint(info.reserved2),
		reserved3:      uint(info.reserved3),
	}, err(result)
}

func (d *Device) TemperatureThreshold() (TemperatureThreshold, error) {
	var slowdown C.uint
	var shutdown C.uint

	result := C.nvmlDeviceGetTemperatureThreshold(d.device, C.NVML_TEMPERATURE_THRESHOLD_SLOWDOWN, &slowdown)
	if e := err(result); e != nil {
		return TemperatureThreshold{}, e
	}
	result = C.nvmlDeviceGetTemperatureThreshold(d.device, C.NVML_TEMPERATURE_THRESHOLD_SHUTDOWN, &shutdown)

	return TemperatureThreshold{
		Slowdown: uint(slowdown),
		Shutdown: uint(shutdown),
	}, err(result)
}

func (d *Device) Temperature() (uint, error) {
	var temperature C.uint
	result := C.nvmlDeviceGetTemperature(d.device, C.NVML_TEMPERATURE_GPU, &temperature)
	return uint(temperature), err(result)
}

func (d *Device) UtilizationRates() (UtilizationRate, error) {
	var utilization C.nvmlUtilization_t
	result := C.nvmlDeviceGetUtilizationRates(d.device, &utilization)
	return UtilizationRate{
		GPU:    uint(utilization.gpu),
		Memory: uint(utilization.memory),
	}, err(result)
}

func (d *Device) FanSpeed() (uint, error) {
	var speed C.uint
	result := C.nvmlDeviceGetFanSpeed(d.device, &speed)
	return uint(speed), err(result)
}

func (d *Device) Memory() (Memory, error) {
	var memory C.nvmlMemory_t
	result := C.nvmlDeviceGetMemoryInfo(d.device, &memory)
	return Memory{
		Total: uint64(memory.total),
		Free:  uint64(memory.free),
		Used:  uint64(memory.used),
	}, err(result)
}

func (d *Device) Bar1Memory() (Memory, error) {
	var memory C.nvmlBAR1Memory_t
	result := C.nvmlDeviceGetBAR1MemoryInfo(d.device, &memory)
	return Memory{
		Total: uint64(memory.bar1Total),
		Free:  uint64(memory.bar1Free),
		Used:  uint64(memory.bar1Used),
	}, err(result)
}

func (d *Device) Clock() (Clock, error) {
	var clock, maxClock C.uint
	var res Clock

	result := C.nvmlDeviceGetMaxClockInfo(d.device, C.NVML_CLOCK_GRAPHICS, &maxClock)
	if e := err(result); e != nil {
		return res, e
	}
	res.GraphicsMax = uint(maxClock)

	result = C.nvmlDeviceGetClockInfo(d.device, C.NVML_CLOCK_GRAPHICS, &clock)
	if e := err(result); e != nil {
		return res, e
	}
	res.Graphics = uint(clock)

	result = C.nvmlDeviceGetMaxClockInfo(d.device, C.NVML_CLOCK_SM, &maxClock)
	if e := err(result); e != nil {
		return res, e
	}
	res.SMMax = uint(maxClock)

	result = C.nvmlDeviceGetClockInfo(d.device, C.NVML_CLOCK_SM, &clock)
	if e := err(result); e != nil {
		return res, e
	}
	res.SM = uint(clock)

	result = C.nvmlDeviceGetMaxClockInfo(d.device, C.NVML_CLOCK_MEM, &maxClock)
	if e := err(result); e != nil {
		return res, e
	}
	res.MemMax = uint(maxClock)

	result = C.nvmlDeviceGetClockInfo(d.device, C.NVML_CLOCK_MEM, &clock)
	if e := err(result); e != nil {
		return res, e
	}
	res.Mem = uint(clock)

	return res, nil
}

func (d *Device) RunningProcess() ([]Process, error) {
	var processCount C.uint
	var processes = make([]C.nvmlProcessInfo_t, 1)

	for {
		result := C.nvmlDeviceGetComputeRunningProcesses(d.device, &processCount, &processes[0])
		if result == C.NVML_ERROR_INSUFFICIENT_SIZE {
			processes = make([]C.nvmlProcessInfo_t, uint(processCount))
		} else if e := err(result); e != nil {
			return nil, e
		} else {
			break
		}
	}

	var res = make([]Process, uint(processCount))
	var i uint = 0
	for ; i < uint(processCount); i++ {
		b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/comm", uint(processes[i].pid)))
		if err != nil {
			return nil, err
		}

		res[i].PID = uint(processes[i].pid)
		res[i].MemoryUsed = uint64(processes[i].usedGpuMemory)
		res[i].Name = strings.Trim(string(b), "\n")
	}

	return res, nil
}

func (d *Device) PowerUsage() (uint, error) {
	var usage C.uint
	result := C.nvmlDeviceGetPowerUsage(d.device, &usage)
	return uint(usage), err(result)
}

func (d *Device) ComputeMode() (string, error) {
	var mode C.nvmlComputeMode_t
	result := C.nvmlDeviceGetComputeMode(d.device, &mode)
	return computeModeString(mode), err(result)
}

func (d *Device) PerformanceState() (string, error) {
	var state C.nvmlPstates_t
	result := C.nvmlDeviceGetPerformanceState(d.device, &state)
	return stateString(state), err(result)
}

func (d *Device) PowerState() (string, error) {
	var state C.nvmlPstates_t
	result := C.nvmlDeviceGetPowerState(d.device, &state)
	return stateString(state), err(result)
}

func computeModeString(mode C.nvmlComputeMode_t) string {
	switch mode {
	case C.NVML_COMPUTEMODE_DEFAULT:
		return "Default"
	case C.NVML_COMPUTEMODE_EXCLUSIVE_THREAD:
		return "ExclusiveThread"
	case C.NVML_COMPUTEMODE_PROHIBITED:
		return "Prohibited"
	case C.NVML_COMPUTEMODE_EXCLUSIVE_PROCESS:
		return "ExclusiveProcess"
	default:
		return "Unknown"
	}
}

func stateString(state C.nvmlPstates_t) string {
	if state == C.NVML_PSTATE_UNKNOWN {
		return "UNKNOWN"
	}

	return fmt.Sprintf("P%v", state)
}
