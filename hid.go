package ledger

// #cgo LDFLAGS:
// #include "./hidapi/hidapi/hidapi.h"
// #include "./hidapi/windows/hid.c"
//
import "C"
import (
	"crypto/rand"
	"fmt"
	"unicode/utf16"
	"unsafe"
)

const READ_BUFF_MAXSIZE = 2048

/** hidapi info structure */
type HidDeviceInfo struct {
	/** Platform-specific device path */
	Path 		string
	/** Device Vendor ID */
	VendorId	uint16
	/** Device Product ID */
	ProductId	uint16
	/** Serial Number */
	SerialNumber	string
	/** Device Release Number in binary-coded decimal,
	    also known as Device Version Number */
	ReleaseNumber	uint16
	/** Manufacturer String */
	Manufacturer	string
	/** Product string */
	Product		string
	/** Usage Page for this Device/Interface
	    (Windows/Mac/hidraw only) */
	UsagePage	uint16
	/** Usage for this Device/Interface
	    (Windows/Mac/hidraw only) */
	Usage		uint16
	/** The USB interface which this logical device
	    represents.
	    * Valid on both Linux implementations in all cases.
	    * Valid on the Windows implementation only if the device
	      contains more than one interface.
	    * Valid on the Mac implementation if and only if the device
	      is a USB HID device. */
	InterfaceNumber	int
}

type HidDevice struct {
	Info		HidDeviceInfo
	hidHandle	*C.hid_device

	channel		int
}

func (device *HidDevice) Open() error {
	device.closeHandle()
	path := C.CString(device.Info.Path)
	defer C.free(unsafe.Pointer(path))
	device.hidHandle = C.hid_open_path(path)
	if (device.hidHandle == nil) {
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

func (device *HidDevice) read() []byte {
	buff := make([]byte, READ_BUFF_MAXSIZE)
	returnedLength := C.hid_read(device.hidHandle, (*C.uchar)(&buff[0]), READ_BUFF_MAXSIZE)
	if returnedLength == -1 {
		return nil
	}
	return buff[:returnedLength]
}

func (device *HidDevice) readTimeout(timeout int) []byte {
	buff := make([]byte, READ_BUFF_MAXSIZE)
	returnedLength := C.hid_read_timeout(device.hidHandle, (*C.uchar)(&buff[0]), READ_BUFF_MAXSIZE, C.int(timeout))
	if returnedLength == -1 {
		return nil
	}
	return buff[:returnedLength]
}

func (device *HidDevice) close() {
	device.closeHandle();
}

func (device *HidDevice) setNonBlocking(blockStatus int) bool {
	res := C.hid_set_nonblocking(device.hidHandle, C.int(blockStatus))
	if res < 0 {
		return false
	}
	return true
}

func (device *HidDevice) write(buffer []byte, writeLength int) int {
	if device.hidHandle == nil {
		return -1
	}

	if writeLength <= 0 || writeLength > len(buffer) {
		return -1
	}

	returnedLength := C.hid_write(device.hidHandle, (*C.uchar)(&buffer[0]), C.ulonglong(writeLength));
	if returnedLength < 0 {
		return -1
	}

	return int(returnedLength)
}

func addbuf(buf []uint16, newcap int)(newbuf []uint16) {
	newbuf = make([]uint16,newcap)
	copy(newbuf,buf)
	return
}

func Utf16prt2str(p uintptr)(str string){
	len := 0
	buf := make([]uint16,64)
	for a := (*(*uint16)(unsafe.Pointer(p))); a != 0; len++ {
		if len >= cap(buf){
			buf = addbuf(buf,len*2)
		}
		buf[len] = a
		p += 2//uint16 occupies 2 bytes
		a = (*(*uint16)(unsafe.Pointer(p)))
	}
	str = string(utf16.Decode(buf[:len]))
	return
}

func GetDevices(vendorId int, productId int) []*HidDevice {
	devs := C.hid_enumerate(C.ushort(vendorId), C.ushort(productId))
	if devs == nil {
		return nil
	}
	devices := make([]*HidDevice, 0)

	for dev := devs; dev != nil; dev = dev.next {
		device := &HidDevice{}
		b := make([]byte, 2)
		_, err := rand.Read(b)
		if err != nil {
			return nil
		}
		device.channel = int(b[1]) << 8 | int(b[0])
		device.Info.VendorId = uint16(dev.vendor_id)
		device.Info.ProductId = uint16(dev.product_id)
		if (dev.path != nil) {
			device.Info.Path = C.GoString((*C.char)(dev.path))
		}
		if (dev.serial_number != nil) {
			device.Info.SerialNumber = Utf16prt2str(uintptr(unsafe.Pointer(dev.serial_number)))
		}
		if (dev.manufacturer_string != nil) {
			device.Info.Manufacturer = Utf16prt2str(uintptr(unsafe.Pointer(dev.manufacturer_string)))
		}
		if (dev.product_string != nil) {
			device.Info.Product = Utf16prt2str(uintptr(unsafe.Pointer(dev.product_string)))
		}
//		deviceInfo.Release = dev.release_number
//		deviceInfo.Interface = dev.interface_number
		device.Info.UsagePage = uint16(dev.usage_page)
		device.Info.Usage = uint16(dev.usage)
		devices = append(devices, device)
	}
	C.hid_free_enumeration(devs)
	return devices
}

func Deinitialize() {
	C.hid_exit()
}

func Initialize() bool {
	if C.hid_init() != 0 {
		return false
	}
	return true
}

