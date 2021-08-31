package ledger

import (
	"unicode/utf16"
	"unsafe"
)

const (
	// LedgerUSBVendorID Ledger USB vendor ID
	LedgerUSBVendorID = 0x2c97
	// READBUFFMAXSIZE Max read buffer size
	READBUFFMAXSIZE = 2048
)

// IHidDevice interface
type IHidDevice interface {
	Open() error
	Close()
	Exchange(apdu []byte) ([]byte, error)
	GetInfo() *HidDeviceInfo
}

// HidDeviceInfo hidapi info structure
type HidDeviceInfo struct {
	/** Platform-specific device path */
	Path string
	/** Device Vendor ID */
	VendorID uint16
	/** Device Product ID */
	ProductID uint16
	/** Serial Number */
	SerialNumber string
	/** Device Release Number in binary-coded decimal,
	  also known as Device Version Number */
	ReleaseNumber uint16
	/** Manufacturer String */
	Manufacturer string
	/** Product string */
	Product string
	/** Usage Page for this Device/Interface
	  (Windows/Mac/hidraw only) */
	UsagePage uint16
	/** Usage for this Device/Interface
	  (Windows/Mac/hidraw only) */
	Usage uint16
	/** The USB interface which this logical device
	  represents.
	  * Valid on both Linux implementations in all cases.
	  * Valid on the Windows implementation only if the device
	    contains more than one interface.
	  * Valid on the Mac implementation if and only if the device
	    is a USB HID device. */
	InterfaceNumber int
}

// resize byte array
func addbuf(buf []uint16, newcap int) (newbuf []uint16) {
	newbuf = make([]uint16, newcap)
	copy(newbuf, buf)
	return
}

// Utf16prt2str Convert UTF16 to string
func Utf16prt2str(p uintptr) (str string) {
	len := 0
	buf := make([]uint16, 64)
	for a := (*(*uint16)(unsafe.Pointer(p))); a != 0; len++ {
		if len >= cap(buf) {
			buf = addbuf(buf, len*2)
		}
		buf[len] = a
		p += 2 //uint16 occupies 2 bytes
		a = (*(*uint16)(unsafe.Pointer(p)))
	}
	str = string(utf16.Decode(buf[:len]))
	return
}
