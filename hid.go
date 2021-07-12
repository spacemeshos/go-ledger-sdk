package ledger

import (
	"unicode/utf16"
	"unsafe"
)

const LedgerUSBVendorId = 0x2c97

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
