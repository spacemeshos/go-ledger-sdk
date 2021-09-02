package ledger

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
	/** Device Release Number in binary-coded decimal,
	  also known as Device Version Number */
	ReleaseNumber uint16
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

