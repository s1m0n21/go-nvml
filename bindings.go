package nvml

//#cgo CFLAGS: -I /usr/src/gdk/nvml/include -I /usr/include/nvidia/gdk
//#cgo LDFLAGS: -lnvidia-ml -L /usr/src/gdk/nvml/lib
//#include <nvml.h>
import "C"

import "golang.org/x/xerrors"

func err(result C.nvmlReturn_t) error {
	if result == C.NVML_SUCCESS {
		return nil
	}

	return xerrors.New(C.GoString(C.nvmlErrorString(result)))
}

func Init() error {
	return err(C.nvmlInit_v2())
}

func Shutdown() error {
	return err(C.nvmlShutdown())
}

func GetDeviceCount() (uint, error) {
	var count C.uint
	result := C.nvmlDeviceGetCount_v2(&count)
	return uint(count), err(result)
}

func GetDeviceByIndex(idx uint) (Device, error) {
	var dev C.nvmlDevice_t
	result := C.nvmlDeviceGetHandleByIndex_v2(C.uint(idx), &dev)
	return Device{dev}, err(result)
}

func GetDriverVersion() (string, error) {
	var version [C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE]C.char
	result := C.nvmlSystemGetDriverVersion(&version[0], C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE)
	return C.GoString(&version[0]), err(result)
}
