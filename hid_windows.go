package ledger

/*
#cgo LDFLAGS:
#include "./hidapi/hidapi/hidapi.h"
#include "./hidapi/windows/hid.c"
*/
import "C"
import (
	"crypto/rand"
	"fmt"
	"unsafe"
)

// HidDevice Wrapper for HIDAPI library
type HidDevice struct {
	Info      HidDeviceInfo
	hidHandle *C.hid_device

	channel int
}

// Open Open Ledger device for communication.
//
// return {error} Error value.
//
// example
// devices := ledger.GetDevices(0)
// if devices != nil && len(devices) > 0 {
// 	device := devices[0]
// 	if err := device.Open(); err == nil {
// 		...
// 		device.Close()
// 	} else {
// 		fmt.Printf("Open device ERROR: %v\n", err)
// 	}
// }
//
func (device *HidDevice) Open() error {
	device.closeHandle()
	path := C.CString(device.Info.Path)
	defer C.free(unsafe.Pointer(path))
	device.hidHandle = C.hid_open_path(path)
	if device.hidHandle == nil {
		return fmt.Errorf("cannot open device with path %v", device.Info.Path)
	}
	return nil
}

func (device *HidDevice) closeHandle() {
	if device.hidHandle != nil {
		C.hid_close(device.hidHandle)
		device.hidHandle = nil
	}
}

// Read Read data from Ledger
func (device *HidDevice) read() []byte {
	buff := make([]byte, ReadBuffMaxSize)
	returnedLength := C.hid_read(device.hidHandle, (*C.uchar)(&buff[0]), ReadBuffMaxSize)
	if returnedLength == -1 {
		return nil
	}
	return buff[:returnedLength]
}

// Close Close communication with Ledger device.
//
// example
// devices := ledger.GetDevices(0)
// if devices != nil && len(devices) > 0 {
// 	device := devices[0]
// 	if err := device.Open(); err == nil {
// 		...
// 		device.Close()
// 	} else {
// 		fmt.Printf("Open device ERROR: %v\n", err)
// 	}
// }
//
func (device *HidDevice) Close() {
	device.closeHandle()
}

// GetInfo Get HID device info
func (device *HidDevice) GetInfo() *HidDeviceInfo {
	return &device.Info
}

// Write Write data to Ledger
func (device *HidDevice) write(buffer []byte, writeLength int) int {
	if device.hidHandle == nil {
		return -1
	}

	if writeLength <= 0 || writeLength > len(buffer) {
		return -1
	}

	returnedLength := C.hid_write(device.hidHandle, (*C.uchar)(&buffer[0]), C.ulonglong(writeLength))
	if returnedLength < 0 {
		return -1
	}

	return int(returnedLength)
}

// GetDevices Enumerate Ledger devices.
//
// param {int} productId USB Product ID filter, 0 - all.
// return {[]*HidDevice} Discovered Ledger devices.
//
// example
// devices := ledger.GetDevices(0)
// if devices != nil && len(devices) > 0 {
// 	device := devices[0]
// 	if err := device.Open(); err == nil {
// 		...
// 		device.Close()
// 	} else {
// 		fmt.Printf("Open device ERROR: %v\n", err)
// 	}
// }
//
func GetDevices(productID int) []*Ledger {
	devs := C.hid_enumerate(C.ushort(LedgerUSBVendorID), C.ushort(productID))
	if devs == nil {
		return nil
	}
	defer C.hid_free_enumeration(devs)
	devices := make([]*Ledger, 0)

	for dev := devs; dev != nil; dev = dev.next {
		if dev.usage_page != 65440 {
			continue
		}
		device := &HidDevice{}
		b := make([]byte, 2)
		_, err := rand.Read(b)
		if err != nil {
			return nil
		}
		device.channel = int(b[1])<<8 | int(b[0])
		device.Info.VendorID = uint16(dev.vendor_id)
		device.Info.ProductID = uint16(dev.product_id)
		if dev.path != nil {
			device.Info.Path = C.GoString((*C.char)(dev.path))
		}
		device.Info.UsagePage = uint16(dev.usage_page)
		device.Info.Usage = uint16(dev.usage)
		devices = append(devices, &Ledger{hid: device})
	}

	return devices
}
